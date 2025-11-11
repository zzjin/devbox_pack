package generators

import (
	"testing"

	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewExecutionPlanGenerator(t *testing.T) {
	generator := NewExecutionPlanGenerator()
	if generator == nil {
		t.Fatal("NewExecutionPlanGenerator() returned nil")
	}
}

func TestGeneratePlan_NodeProject(t *testing.T) {
	generator := NewExecutionPlanGenerator()

	detectResults := []types.DetectResult{
		{
			Matched:    true,
			Language:   "node",
			Framework:  "express",
			Confidence: 0.9,
			Evidence: types.Evidence{
				Files: []string{"package.json"},
			},
		},
	}

	options := types.CLIOptions{
		Verbose: true,
		Format:  "json",
	}

	plan, err := generator.GeneratePlan(detectResults, options)
	if err != nil {
		t.Fatalf("GeneratePlan failed: %v", err)
	}

	if plan == nil {
		t.Fatal("GeneratePlan returned nil plan")
	}

	if plan.Provider != "node" {
		t.Errorf("expected provider 'node', got '%s'", plan.Provider)
	}

	if plan.Provider != "node" {
		t.Errorf("expected provider 'node', got '%s'", plan.Provider)
	}

	if len(plan.Commands.Setup) == 0 && len(plan.Commands.Dev) == 0 {
		t.Error("plan has no build or dev commands")
	}
}

func TestGeneratePlan_PythonProject(t *testing.T) {
	generator := NewExecutionPlanGenerator()

	detectResults := []types.DetectResult{
		{
			Matched:    true,
			Language:   "python",
			Framework:  "",
			Confidence: 0.8,
			Evidence: types.Evidence{
				Files: []string{"main.py", "requirements.txt"},
			},
		},
	}

	options := types.CLIOptions{
		Verbose: true,
		Format:  "json",
	}

	plan, err := generator.GeneratePlan(detectResults, options)
	if err != nil {
		t.Fatalf("GeneratePlan failed: %v", err)
	}

	if plan == nil {
		t.Fatal("GeneratePlan returned nil plan")
	}

	if plan.Provider != "python" {
		t.Errorf("expected provider 'python', got '%s'", plan.Provider)
	}

	if plan.Provider != "python" {
		t.Errorf("expected provider 'python', got '%s'", plan.Provider)
	}
}

func TestGeneratePlan_StaticProject(t *testing.T) {
	generator := NewExecutionPlanGenerator()

	detectResults := []types.DetectResult{
		{
			Matched:    true,
			Language:   "staticfile",
			Framework:  "",
			Confidence: 0.7,
			Evidence: types.Evidence{
				Files: []string{"index.html"},
			},
		},
	}

	options := types.CLIOptions{
		Verbose: true,
		Format:  "json",
	}

	plan, err := generator.GeneratePlan(detectResults, options)
	if err != nil {
		t.Fatalf("GeneratePlan failed: %v", err)
	}

	if plan == nil {
		t.Fatal("GeneratePlan returned nil plan")
	}

	if plan.Provider != "staticfile" {
		t.Errorf("expected provider 'staticfile', got '%s'", plan.Provider)
	}

	if plan.Provider != "staticfile" {
		t.Errorf("expected provider 'staticfile', got '%s'", plan.Provider)
	}
}

func TestGeneratePlan_InvalidInput(t *testing.T) {
	generator := NewExecutionPlanGenerator()

	// Test with empty results
	_, err := generator.GeneratePlan([]types.DetectResult{}, types.CLIOptions{})
	if err == nil {
		t.Error("expected error for empty detection results")
	}

	// Test with nil results
	_, err = generator.GeneratePlan(nil, types.CLIOptions{})
	if err == nil {
		t.Error("expected error for nil detection results")
	}
}
