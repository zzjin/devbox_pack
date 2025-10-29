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
		Provider: bestResult.Language,
		Runtime:  g.generateRuntime(bestResult, options),
		Base:     g.generateBase(bestResult, options),
		Commands: g.generateCommands(bestResult, options),
		Port:     g.getPortForResult(bestResult),
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

// selectBestResult selects the best detection result
func (g *ExecutionPlanGenerator) selectBestResult(results []types.DetectResult) *types.DetectResult {
	if len(results) == 0 {
		return nil
	}

	var bestResult *types.DetectResult
	maxConfidence := 0.0

	for i := range results {
		result := &results[i]
		if result.Confidence > maxConfidence {
			maxConfidence = result.Confidence
			bestResult = result
		}
	}

	return bestResult
}

// generateRuntime generates runtime configuration
func (g *ExecutionPlanGenerator) generateRuntime(result *types.DetectResult, _ types.CLIOptions) types.RuntimeConfig {
	version := g.getDefaultVersion(types.SupportedLanguage(result.Language))

	// If there is version information in the detection result, use the detected version
	if result.Version != "" {
		version = &result.Version
	}

	runtime := types.RuntimeConfig{
		Language: result.Language,
	}

	// Add version only when there is a value
	if version != nil {
		runtime.Version = version
	}

	return runtime
}

// generateBase generates base configuration
func (g *ExecutionPlanGenerator) generateBase(result *types.DetectResult, _ types.CLIOptions) types.BaseConfig {
	base := types.BaseConfig{}

	// Get base image from catalog
	if catalog, exists := g.baseCatalog[types.SupportedLanguage(result.Language)]; exists {
		if image, exists := catalog["image"]; exists {
			base.Name = image
		}
	}

	return base
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
