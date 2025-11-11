// Package generators provides execution plan generation functionality
// for the DevBox Pack system.
package generators

import (
	"fmt"

	"github.com/labring/devbox-pack/pkg/registry"
	"github.com/labring/devbox-pack/pkg/types"
	"github.com/labring/devbox-pack/pkg/utils"
)

const (
	// DefaultPort is the default port used when no language-specific port is configured
	DefaultPort = 8000
)

// ExecutionPlanGenerator execution plan generator
type ExecutionPlanGenerator struct {
	baseCatalog     map[types.SupportedLanguage]map[string]string
	defaultPorts    map[types.SupportedLanguage]int
	defaultVersions map[types.SupportedLanguage]string
	registry        *registry.ProviderRegistry
}

// NewExecutionPlanGenerator creates a new execution plan generator
func NewExecutionPlanGenerator() *ExecutionPlanGenerator {
	return &ExecutionPlanGenerator{
		baseCatalog:     utils.BaseCatalog,
		defaultPorts:    utils.DefaultPorts.Languages,
		defaultVersions: utils.DefaultVersions,
		registry:        registry.NewProviderRegistry(),
	}
}

// GeneratePlan generates execution plan
func (g *ExecutionPlanGenerator) GeneratePlan(results []types.DetectResult, options types.CLIOptions) (*types.ExecutionPlan, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("no detection results provided")
	}

	// Select the best detection result
	bestResult := g.selectBestResult(results)
	if bestResult == nil {
		return nil, fmt.Errorf("no valid detection result found")
	}

	// Generate execution plan
	plan := &types.ExecutionPlan{
		Provider:    bestResult.Language,
		Runtime:     g.generateRuntime(bestResult, options),
		Environment: g.generateEnvironment(bestResult, options),
		Commands:    g.generateCommands(bestResult, options),
		Port:        g.getPortForResult(bestResult),
	}

	// Add Apt dependencies only when there are values
	aptDeps := g.generateAptDependencies(bestResult, options)
	if len(aptDeps) > 0 {
		plan.Apt = aptDeps
	}

	// Add Evidence only when there are values
	if len(bestResult.Evidence.Files) > 0 || bestResult.Evidence.Reason != "" {
		plan.Evidence = bestResult.Evidence
	}

	return plan, nil
}

// selectBestResult selects the best detection result with backend-first priority
func (g *ExecutionPlanGenerator) selectBestResult(results []types.DetectResult) *types.DetectResult {
	if len(results) == 0 {
		return nil
	}

	// Separate backend and frontend results
	var backendResults []*types.DetectResult
	var frontendResults []*types.DetectResult

	for i := range results {
		result := &results[i]
		if g.isBackendFramework(result) {
			backendResults = append(backendResults, result)
		} else {
			frontendResults = append(frontendResults, result)
		}
	}

	// Prioritize backend frameworks for full-stack applications
	if len(backendResults) > 0 {
		// Return backend result with highest confidence
		bestBackend := backendResults[0]
		for _, result := range backendResults[1:] {
			if result.Confidence > bestBackend.Confidence {
				bestBackend = result
			}
		}
		return bestBackend
	}

	// If no backend results, return frontend result with highest confidence
	bestFrontend := frontendResults[0]
	for _, result := range frontendResults[1:] {
		if result.Confidence > bestFrontend.Confidence {
			bestFrontend = result
		}
	}
	return bestFrontend
}

// isBackendFramework checks if a detection result represents a backend framework
func (g *ExecutionPlanGenerator) isBackendFramework(result *types.DetectResult) bool {
	backendLanguages := []string{"php", "python", "java", "go", "ruby", "rust"}

	// Check if language is backend
	for _, lang := range backendLanguages {
		if result.Language == lang {
			return true
		}
	}

	// Specific backend frameworks that might be detected in other languages
	backendFrameworks := []string{
		"Laravel", "Symfony", "Django", "Flask", "FastAPI", "Spring Boot",
		"Ruby on Rails", "Express", "Koa", "NestJS", "Gin", "Echo", "Fiber",
	}

	for _, framework := range backendFrameworks {
		if result.Framework == framework {
			return true
		}
	}

	return false
}

// generateRuntime generates simplified runtime configuration
func (g *ExecutionPlanGenerator) generateRuntime(result *types.DetectResult, _ types.CLIOptions) types.RuntimeConfig {
	runtime := types.RuntimeConfig{}

	// Get base image from catalog using detected version
	if catalog, exists := g.baseCatalog[types.SupportedLanguage(result.Language)]; exists {
		version := result.Version
		if version == "" {
			// Use default version if none detected
			if defaultVersion, exists := g.defaultVersions[types.SupportedLanguage(result.Language)]; exists {
				version = defaultVersion
			}
		}

		if version != "" {
			if image, exists := catalog[version]; exists {
				runtime.Image = image
			}
		}
	}

	// Add framework if detected
	if result.Framework != "" {
		runtime.Framework = &result.Framework
	}

	return runtime
}

// generateEnvironment generates environment variables (flattened)
func (g *ExecutionPlanGenerator) generateEnvironment(result *types.DetectResult, _ types.CLIOptions) map[string]string {
	provider := g.registry.GetProvider(result.Language)
	if provider == nil {
		return nil
	}

	return provider.GenerateEnvironment(result)
}

// generateAptDependencies generates apt dependencies
func (g *ExecutionPlanGenerator) generateAptDependencies(result *types.DetectResult, _ types.CLIOptions) []string {
	var deps []string

	// Check if native compilation is needed
	if g.needsNativeCompilation(result) {
		deps = append(deps, "build-essential")

		// Add language-specific build dependencies
		switch types.SupportedLanguage(result.Language) {
		case types.LanguagePython:
			deps = append(deps, "python3-dev", "libffi-dev", "libssl-dev")
		case types.LanguageNode:
			deps = append(deps, "node-gyp")
		case types.LanguageRuby:
			deps = append(deps, "ruby-dev", "libsqlite3-dev")
		case types.LanguageJava:
			// Java typically doesn't need additional native compilation deps
		case types.LanguageGo:
			// Go typically doesn't need additional native compilation deps
		case types.LanguagePHP:
			deps = append(deps, "php-dev")
		case types.LanguageDeno:
			// Deno typically doesn't need additional native compilation deps
		case types.LanguageRust:
			deps = append(deps, "build-essential")
		case types.LanguageStaticfile:
			// Static files don't need compilation deps
		case types.LanguageShell:
			// Shell scripts don't need compilation deps
		}
	}

	// Add git if needed
	if result.Evidence.Reason == "git repository" {
		deps = append(deps, "git")
	}

	return deps
}

// needsNativeCompilation checks if native compilation is needed
func (g *ExecutionPlanGenerator) needsNativeCompilation(result *types.DetectResult) bool {
	provider := g.registry.GetProvider(result.Language)
	if provider == nil {
		return false
	}

	return provider.NeedsNativeCompilation(result)
}

// generateCommands generates commands configuration
func (g *ExecutionPlanGenerator) generateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	provider := g.registry.GetProvider(result.Language)
	if provider == nil {
		return types.Commands{}
	}

	return provider.GenerateCommands(result, options)
}

// getPortForResult gets the port for a detection result
func (g *ExecutionPlanGenerator) getPortForResult(result *types.DetectResult) int {
	if port, exists := g.defaultPorts[types.SupportedLanguage(result.Language)]; exists {
		return port
	}
	return DefaultPort
}

// getDefaultVersion gets the default version for a language
func (g *ExecutionPlanGenerator) getDefaultVersion(language types.SupportedLanguage) *string {
	if version, exists := g.defaultVersions[language]; exists {
		return &version
	}
	return nil
}
