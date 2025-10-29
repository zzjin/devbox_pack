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

// ExecutionPlan execution plan structure
type ExecutionPlan struct {
	Provider string `json:"provider"`
	Base     struct {
		Name     string `json:"name"`
		Platform string `json:"platform"`
	} `json:"base"`
	Runtime struct {
		Language string                 `json:"language"`
		Version  string                 `json:"version"`
		Tools    interface{}            `json:"tools"`
		Env      map[string]interface{} `json:"env"`
	} `json:"runtime"`
	Apt      interface{} `json:"apt"`
	AptDeps  interface{} `json:"aptDeps"`
	Dev      *Command    `json:"dev"`
	Build    *Command    `json:"build"`
	Start    *Command    `json:"start"`
	Commands struct {
		Dev   []string `json:"dev"`
		Build []string `json:"build"`
		Start []string `json:"start"`
	} `json:"commands"`
	Evidence struct {
		Files  interface{} `json:"files"`
		Reason string      `json:"reason"`
	} `json:"evidence"`
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

	// Validate base image
	if plan.Base.Name == "" {
		result.Errors = append(result.Errors, "base image name cannot be empty")
		result.Valid = false
	}

	// Validate runtime
	if plan.Runtime.Language == "" {
		result.Errors = append(result.Errors, "runtime language cannot be empty")
		result.Valid = false
	}

	// Validate tools
	if knownTools, exists := v.knownCommands[plan.Provider]; exists {
		if tools, ok := plan.Runtime.Tools.([]interface{}); ok {
			for _, tool := range tools {
				if toolStr, ok := tool.(string); ok {
					if toolStr == "" {
						result.Warnings = append(result.Warnings, "detected empty tool name")
						continue
					}
					found := false
					for _, knownTool := range knownTools {
						if strings.Contains(toolStr, knownTool) {
							found = true
							break
						}
					}
					if !found {
						result.Warnings = append(result.Warnings, fmt.Sprintf("unknown tool: %s (Provider: %s)", toolStr, plan.Provider))
					}
				}
			}
		}
	}

	// Check commands - prioritize Commands field, if empty then check traditional fields
	hasCommands := len(plan.Commands.Dev) > 0 || len(plan.Commands.Build) > 0 || len(plan.Commands.Start) > 0

	if !hasCommands {
		// If Commands field is empty, check traditional fields
		if plan.Dev == nil || plan.Dev.Cmd == "" {
			result.Errors = append(result.Errors, "dev command cannot be empty")
			result.Valid = false
		}

		if plan.Build == nil || plan.Build.Cmd == "" {
			result.Errors = append(result.Errors, "build command cannot be empty")
			result.Valid = false
		}

		if plan.Start == nil || plan.Start.Cmd == "" {
			result.Errors = append(result.Errors, "start command cannot be empty")
			result.Valid = false
		}
	} else {
		// If Commands field has content, validate its completeness
		if len(plan.Commands.Dev) == 0 {
			result.Warnings = append(result.Warnings, "dev command list is empty")
		}
		if len(plan.Commands.Build) == 0 {
			result.Warnings = append(result.Warnings, "build command list is empty")
		}
		if len(plan.Commands.Start) == 0 {
			result.Warnings = append(result.Warnings, "start command list is empty")
		}
	}

	// Check port configuration
	hasPortConfig := false
	if hasCommands {
		// Check port configuration in Commands field
		for _, cmd := range plan.Commands.Start {
			if strings.Contains(cmd, "PORT") || strings.Contains(cmd, "port") {
				hasPortConfig = true
				break
			}
		}
	} else if plan.Start != nil && plan.Start.Cmd != "" {
		// Check port configuration in traditional field
		hasPortConfig = strings.Contains(plan.Start.Cmd, "PORT") || strings.Contains(plan.Start.Cmd, "port")
	}

	if !hasPortConfig && (hasCommands && len(plan.Commands.Start) > 0 || !hasCommands && plan.Start != nil && plan.Start.Cmd != "") {
		result.Warnings = append(result.Warnings, "start command lacks port environment variable")
	}

	// Validate evidence files
	if plan.Evidence.Files == nil {
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
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Node.js project lacks package.json file")
		result.Warnings = append(result.Warnings, "Node.js project lacks lock file")
		return
	}

	// Check if package.json exists
	hasPackageJSON := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok && strings.Contains(fileStr, "package.json") {
				hasPackageJSON = true
				break
			}
		}
	}
	if !hasPackageJSON {
		result.Warnings = append(result.Warnings, "Node.js project lacks package.json file")
	}

	// Check consistency between tools and lock files
	hasLockFile := false
	lockFiles := []string{"package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb"}
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok {
				for _, lockFile := range lockFiles {
					if strings.Contains(fileStr, lockFile) {
						hasLockFile = true
						break
					}
				}
				if hasLockFile {
					break
				}
			}
		}
	}
	if !hasLockFile {
		result.Warnings = append(result.Warnings, "Node.js project lacks lock file")
	}
}

// validatePythonProject validates Python project
func (v *PlanValidator) validatePythonProject(plan *ExecutionPlan, result *ValidationResult) {
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Python project lacks dependency files")
		return
	}

	// Check dependency files
	hasDepsFile := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok {
				if strings.Contains(fileStr, "requirements.txt") ||
					strings.Contains(fileStr, "pyproject.toml") ||
					strings.Contains(fileStr, "Pipfile") ||
					strings.Contains(fileStr, "poetry.lock") {
					hasDepsFile = true
					break
				}
			}
		}
	}
	if !hasDepsFile {
		result.Warnings = append(result.Warnings, "Python project lacks dependency files")
	}
}

// validateJavaProject validates Java project
func (v *PlanValidator) validateJavaProject(plan *ExecutionPlan, result *ValidationResult) {
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Java project lacks build files (pom.xml or build.gradle)")
		return
	}

	// Check build files
	hasBuildFile := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok {
				if strings.Contains(fileStr, "pom.xml") || strings.Contains(fileStr, "build.gradle") {
					hasBuildFile = true
					break
				}
			}
		}
	}
	if !hasBuildFile {
		result.Warnings = append(result.Warnings, "Java project lacks build files (pom.xml or build.gradle)")
	}
}

// validateGoProject validates Go project
func (v *PlanValidator) validateGoProject(plan *ExecutionPlan, result *ValidationResult) {
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Go project lacks go.mod file")
		return
	}

	// Check go.mod
	hasGoMod := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok && strings.Contains(fileStr, "go.mod") {
				hasGoMod = true
				break
			}
		}
	}
	if !hasGoMod {
		result.Warnings = append(result.Warnings, "Go project lacks go.mod file")
	}
}

// validateRustProject validates Rust project
func (v *PlanValidator) validateRustProject(plan *ExecutionPlan, result *ValidationResult) {
	if plan.Evidence.Files == nil {
		result.Warnings = append(result.Warnings, "Rust project lacks Cargo.toml file")
		return
	}

	// Check Cargo.toml
	hasCargoToml := false
	if files, ok := plan.Evidence.Files.([]interface{}); ok {
		for _, file := range files {
			if fileStr, ok := file.(string); ok && strings.Contains(fileStr, "Cargo.toml") {
				hasCargoToml = true
				break
			}
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
