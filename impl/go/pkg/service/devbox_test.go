package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewDevBoxPack(t *testing.T) {
	devbox := NewDevBoxPack()
	if devbox == nil {
		t.Fatal("NewDevBoxPack() returned nil")
	}

	if devbox.gitHandler == nil {
		t.Fatal("gitHandler is nil")
	}

	if devbox.detectionEngine == nil {
		t.Fatal("detectionEngine is nil")
	}

	if devbox.planGenerator == nil {
		t.Fatal("planGenerator is nil")
	}

	if devbox.outputUtils == nil {
		t.Fatal("outputUtils is nil")
	}
}

func TestGeneratePlan_InvalidPath(t *testing.T) {
	devbox := NewDevBoxPack()
	options := &types.CLIOptions{
		Verbose: false,
		Format:  "json",
	}

	// Test with non-existent path
	_, err := devbox.GeneratePlan("/non/existent/path", options)
	if err == nil {
		t.Fatal("Expected error for non-existent path, but got nil")
	}
}

func TestGeneratePlan_ValidPath(t *testing.T) {
	devbox := NewDevBoxPack()
	options := &types.CLIOptions{
		Verbose: false,
		Format:  "json",
	}

	// Create a temporary directory with a simple Node.js project
	tmpDir, err := os.MkdirTemp("", "devbox-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a package.json file
	packageJSON := `{
  "name": "test-project",
  "version": "1.0.0",
  "main": "index.js",
  "scripts": {
    "start": "node index.js"
  }
}`
	err = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Create an index.js file
	indexJS := `console.log("Hello, World!");`
	err = os.WriteFile(filepath.Join(tmpDir, "index.js"), []byte(indexJS), 0644)
	if err != nil {
		t.Fatalf("Failed to create index.js: %v", err)
	}

	// Test plan generation
	plan, err := devbox.GeneratePlan(tmpDir, options)
	if err != nil {
		t.Fatalf("GeneratePlan failed: %v", err)
	}

	if plan == nil {
		t.Fatal("Generated plan is nil")
	}

	// Basic validation of the plan
	if plan.Runtime.Language == "" {
		t.Error("Plan language is empty")
	}

	if plan.Provider == "" {
		t.Error("Plan provider is empty")
	}
}

func TestRun_InvalidPath(t *testing.T) {
	devbox := NewDevBoxPack()
	options := &types.CLIOptions{
		Verbose: false,
		Format:  "json",
	}

	// Test with non-existent path
	err := devbox.Run("/non/existent/path", options)
	if err == nil {
		t.Fatal("Expected error for non-existent path, but got nil")
	}
}
