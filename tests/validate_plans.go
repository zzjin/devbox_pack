// Package main provides validation functionality for DevBox Pack execution plans.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ExecutionPlan execution plan structure (current)
type ExecutionPlan struct {
	Provider    string            `json:"provider"`
	Runtime     RuntimeConfig     `json:"runtime"`
	Environment map[string]string `json:"environment,omitempty"`
	Apt         []string          `json:"apt,omitempty"`
	Commands    Commands          `json:"commands"`
	Port        int               `json:"port"`
	Evidence    Evidence          `json:"evidence,omitempty"`
}

// RuntimeConfig represents the simplified runtime configuration
type RuntimeConfig struct {
	Image     string  `json:"image"`
	Framework *string `json:"framework,omitempty"`
}

// Commands represents the command configuration
type Commands struct {
	Setup []string `json:"setup,omitempty"`
	Dev   []string `json:"dev,omitempty"`
	Build []string `json:"build,omitempty"`
	Run   []string `json:"run,omitempty"`
}

// Evidence represents detection evidence
type Evidence struct {
	Files  []string `json:"files,omitempty"`
	Reason string   `json:"reason,omitempty"`
}

// Command command structure
type Command struct {
	Cmd     string   `json:"cmd"`
	Caches  []string `json:"caches,omitempty"`
	PortEnv string   `json:"portEnv,omitempty"`
	Notes   []string `json:"notes,omitempty"`
}

// ValidationResult validation result
type ValidationResult struct {
	TestCase string         `json:"testCase"`
	Valid    bool           `json:"valid"`
	Errors   []string       `json:"errors"`
	Warnings []string       `json:"warnings"`
	Plan     *ExecutionPlan `json:"plan,omitempty"`
}

// ValidationSummary validation summary
type ValidationSummary struct {
	TotalCases    int                `json:"totalCases"`
	ValidCases    int                `json:"validCases"`
	InvalidCases  int                `json:"invalidCases"`
	SuccessRate   float64            `json:"successRate"`
	Results       []ValidationResult `json:"results"`
	ProviderStats map[string]int     `json:"providerStats"`
}

// PlanValidator validates DevBox Pack execution plans for correctness and completeness.
type PlanValidator struct {
	knownProviders map[string]bool
	knownCommands  map[string][]string
}

// NewPlanValidator creates a new validator
func NewPlanValidator() *PlanValidator {
	return &PlanValidator{
		knownProviders: map[string]bool{
			"node":       true,
			"python":     true,
			"java":       true,
			"go":         true,
			"php":        true,
			"ruby":       true,
			"elixir":     true,
			"deno":       true,
			"rust":       true,
			"staticfile": true,
			"shell":      true,
		},
		knownCommands: map[string][]string{
			"node":   {"npm", "yarn", "pnpm", "bun"},
			"python": {"pip", "poetry", "pipenv", "uv", "pdm"},
			"java":   {"mvn", "gradle"},
			"go":     {"go"},
			"php":    {"composer"},
			"ruby":   {"bundle"},
			"elixir": {"mix"},
			"deno":   {"deno"},
			"rust":   {"cargo"},
		},
	}
}

// ValidatePlan validates a single execution plan
func (v *PlanValidator) ValidatePlan(testCase string, plan *ExecutionPlan) ValidationResult {
	result := ValidationResult{
		TestCase: testCase,
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
		Plan:     plan,
	}

	// Validate Provider
	if plan.Provider == "" {
		result.Errors = append(result.Errors, "provider cannot be empty")
		result.Valid = false
	} else if !v.knownProviders[plan.Provider] {
		result.Errors = append(result.Errors, fmt.Sprintf("unknown Provider: %s", plan.Provider))
		result.Valid = false
	}

	// Validate runtime image
	if plan.Runtime.Image == "" {
		result.Warnings = append(result.Warnings, "runtime image is empty (will be set by generator)")
	}

	// Check commands structure
	if len(plan.Commands.Setup) == 0 {
		result.Warnings = append(result.Warnings, "setup command list is empty")
	}
	if len(plan.Commands.Dev) == 0 {
		result.Warnings = append(result.Warnings, "dev command list is empty")
	}
	if len(plan.Commands.Build) == 0 {
		result.Warnings = append(result.Warnings, "build command list is empty")
	}
	if len(plan.Commands.Run) == 0 {
		result.Warnings = append(result.Warnings, "run command list is empty")
	}

	// Check port configuration
	hasPortConfig := false
	for _, cmd := range plan.Commands.Run {
		if strings.Contains(cmd, "PORT") || strings.Contains(cmd, "port") {
			hasPortConfig = true
			break
		}
	}

	if !hasPortConfig && len(plan.Commands.Run) > 0 {
		result.Warnings = append(result.Warnings, "run command lacks port environment variable")
	}

	// Validate evidence files
	if len(plan.Evidence.Files) == 0 {
		result.Warnings = append(result.Warnings, "no evidence files detected")
	}

	// Validate specific Provider rules
	v.validateProviderSpecific(plan, &result)

	return result
}

// validateProviderSpecific validates specific Provider rules
func (v *PlanValidator) validateProviderSpecific(plan *ExecutionPlan, result *ValidationResult) {
	switch plan.Provider {
	case "node":
		v.validateNodeProject(plan, result)
	case "python":
		v.validatePythonProject(plan, result)
	case "java":
		v.validateJavaProject(plan, result)
	case "go":
		v.validateGoProject(plan, result)
	case "rust":
		v.validateRustProject(plan, result)
	}
}

// validateNodeProject validates Node.js project
func (v *PlanValidator) validateNodeProject(plan *ExecutionPlan, result *ValidationResult) {
	if len(plan.Evidence.Files) == 0 {
		result.Warnings = append(result.Warnings, "Node.js project lacks package.json file")
		result.Warnings = append(result.Warnings, "Node.js project lacks lock file")
		return
	}

	// Check if package.json exists
	hasPackageJSON := false
	for _, file := range plan.Evidence.Files {
		if strings.Contains(file, "package.json") {
			hasPackageJSON = true
			break
		}
	}
	if !hasPackageJSON {
		result.Warnings = append(result.Warnings, "Node.js project lacks package.json file")
	}

	// Check consistency between tools and lock files
	hasLockFile := false
	lockFiles := []string{"package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb"}
	for _, file := range plan.Evidence.Files {
		for _, lockFile := range lockFiles {
			if strings.Contains(file, lockFile) {
				hasLockFile = true
				break
			}
		}
		if hasLockFile {
			break
		}
	}
	if !hasLockFile {
		result.Warnings = append(result.Warnings, "Node.js project lacks lock file")
	}
}

// validatePythonProject validates Python project
func (v *PlanValidator) validatePythonProject(plan *ExecutionPlan, result *ValidationResult) {
	if len(plan.Evidence.Files) == 0 {
		result.Warnings = append(result.Warnings, "Python project lacks dependency files")
		return
	}

	// Check dependency files
	hasDepsFile := false
	for _, file := range plan.Evidence.Files {
		if strings.Contains(file, "requirements.txt") ||
			strings.Contains(file, "pyproject.toml") ||
			strings.Contains(file, "Pipfile") ||
			strings.Contains(file, "poetry.lock") {
			hasDepsFile = true
			break
		}
	}
	if !hasDepsFile {
		result.Warnings = append(result.Warnings, "Python project lacks dependency files")
	}
}

// validateJavaProject validates Java project
func (v *PlanValidator) validateJavaProject(plan *ExecutionPlan, result *ValidationResult) {
	if len(plan.Evidence.Files) == 0 {
		result.Warnings = append(result.Warnings, "Java project lacks build files (pom.xml or build.gradle)")
		return
	}

	// Check build files
	hasBuildFile := false
	for _, file := range plan.Evidence.Files {
		if strings.Contains(file, "pom.xml") || strings.Contains(file, "build.gradle") {
			hasBuildFile = true
			break
		}
	}
	if !hasBuildFile {
		result.Warnings = append(result.Warnings, "Java project lacks build files (pom.xml or build.gradle)")
	}
}

// validateGoProject validates Go project
func (v *PlanValidator) validateGoProject(plan *ExecutionPlan, result *ValidationResult) {
	if len(plan.Evidence.Files) == 0 {
		result.Warnings = append(result.Warnings, "Go project lacks go.mod file")
		return
	}

	// Check go.mod
	hasGoMod := false
	for _, file := range plan.Evidence.Files {
		if strings.Contains(file, "go.mod") {
			hasGoMod = true
			break
		}
	}
	if !hasGoMod {
		result.Warnings = append(result.Warnings, "Go project lacks go.mod file")
	}

	// Check for required environment variables
	if _, exists := plan.Environment["GO_ENV"]; !exists {
		result.Warnings = append(result.Warnings, "Go project missing GO_ENV environment variable")
	}
	if _, exists := plan.Environment["CGO_ENABLED"]; !exists {
		result.Warnings = append(result.Warnings, "Go project missing CGO_ENABLED environment variable")
	}
	if _, exists := plan.Environment["PORT"]; !exists {
		result.Warnings = append(result.Warnings, "Go project missing PORT environment variable")
	}
}

// validateRustProject validates Rust project
func (v *PlanValidator) validateRustProject(plan *ExecutionPlan, result *ValidationResult) {
	if len(plan.Evidence.Files) == 0 {
		result.Warnings = append(result.Warnings, "Rust project lacks Cargo.toml file")
		return
	}

	// Check Cargo.toml
	hasCargoToml := false
	for _, file := range plan.Evidence.Files {
		if strings.Contains(file, "Cargo.toml") {
			hasCargoToml = true
			break
		}
	}
	if !hasCargoToml {
		result.Warnings = append(result.Warnings, "Rust project lacks Cargo.toml file")
	}
}

// Main function
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run validate_plans.go <plans_directory>")
		os.Exit(1)
	}

	plansDir := os.Args[1]
	validator := NewPlanValidator()

	// Get all plan files
	planFiles, err := filepath.Glob(filepath.Join(plansDir, "*.json"))
	if err != nil {
		fmt.Printf("Error: Unable to read plan files: %v\n", err)
		os.Exit(1)
	}

	if len(planFiles) == 0 {
		fmt.Printf("Error: No JSON files found in directory %s\n", plansDir)
		os.Exit(1)
	}

	fmt.Printf("Found %d execution plan files\n", len(planFiles))

	var results []ValidationResult
	providerStats := make(map[string]int)

	// Validate each plan file
	for _, planFile := range planFiles {
		testCase := strings.TrimSuffix(filepath.Base(planFile), ".json")

		// Read plan file
		data, err := os.ReadFile(planFile)
		if err != nil {
			result := ValidationResult{
				TestCase: testCase,
				Valid:    false,
				Errors:   []string{fmt.Sprintf("Unable to read file: %v", err)},
			}
			results = append(results, result)
			continue
		}

		// Parse JSON
		var plan ExecutionPlan
		if err := json.Unmarshal(data, &plan); err != nil {
			result := ValidationResult{
				TestCase: testCase,
				Valid:    false,
				Errors:   []string{fmt.Sprintf("JSON parsing error: %v", err)},
			}
			results = append(results, result)
			continue
		}

		// Validate plan
		result := validator.ValidatePlan(testCase, &plan)
		results = append(results, result)

		// Count Provider statistics
		if result.Valid {
			providerStats[plan.Provider]++
		}
	}

	// Generate summary
	validCases := 0
	for _, result := range results {
		if result.Valid {
			validCases++
		}
	}

	summary := ValidationSummary{
		TotalCases:    len(results),
		ValidCases:    validCases,
		InvalidCases:  len(results) - validCases,
		SuccessRate:   float64(validCases) / float64(len(results)) * 100,
		Results:       results,
		ProviderStats: providerStats,
	}

	// Output results
	fmt.Printf("\n=== Validation Results Summary ===\n")
	fmt.Printf("Total plans: %d\n", summary.TotalCases)
	fmt.Printf("Valid plans: %d\n", summary.ValidCases)
	fmt.Printf("Invalid plans: %d\n", summary.InvalidCases)
	fmt.Printf("Success rate: %.2f%%\n", summary.SuccessRate)

	fmt.Printf("\n=== Provider Statistics ===\n")
	var providers []string
	for provider := range providerStats {
		providers = append(providers, provider)
	}
	sort.Strings(providers)

	for _, provider := range providers {
		fmt.Printf("- %s: %d valid plans\n", provider, providerStats[provider])
	}

	// Output detailed results
	fmt.Printf("\n=== Detailed Validation Results ===\n")
	for _, result := range results {
		status := "✓"
		if !result.Valid {
			status = "✗"
		}

		fmt.Printf("%s %s", status, result.TestCase)
		if result.Plan != nil {
			fmt.Printf(" (Provider: %s)", result.Plan.Provider)
		}
		fmt.Printf("\n")

		if len(result.Errors) > 0 {
			for _, err := range result.Errors {
				fmt.Printf("  Error: %s\n", err)
			}
		}

		if len(result.Warnings) > 0 {
			for _, warning := range result.Warnings {
				fmt.Printf("  Warning: %s\n", warning)
			}
		}
	}

	// Save detailed report
	outputDir := filepath.Dir(plansDir)
	reportFile := filepath.Join(outputDir, "validation_report.json")

	reportData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Printf("Error: Unable to generate report: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(reportFile, reportData, 0600); err != nil {
		fmt.Printf("Error: Unable to save report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDetailed report saved to: %s\n", reportFile)

	// Exit with non-zero status if there are invalid plans
	if summary.InvalidCases > 0 {
		os.Exit(1)
	}
}
