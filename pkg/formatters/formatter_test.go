package formatters

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewJSONFormatter(t *testing.T) {
	formatter := NewJSONFormatter()
	if formatter == nil {
		t.Fatal("NewJSONFormatter() returned nil")
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	formatter := NewJSONFormatter()

	plan := &types.ExecutionPlan{
		Provider: "node",
		Runtime: types.RuntimeConfig{
			Image:     "node:20-alpine",
			Framework: func() *string { s := "express"; return &s }(),
		},
		Environment: map[string]string{
			"NODE_ENV": "production",
			"PORT":     "3000",
		},
		Apt: []string{"curl", "git"},
		Commands: types.Commands{
			Setup: []string{"npm install"},
			Dev:   []string{"npm run dev"},
			Build: []string{"npm run build"},
			Run:   []string{"npm start"},
		},
		Port: 3000,
		Evidence: types.Evidence{
			Files:  []string{"package.json"},
			Reason: "Node.js project detected",
		},
	}

	options := &types.CLIOptions{
		Format:  "json",
		Verbose: false,
	}

	result, err := formatter.Format(plan, options)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if result == "" {
		t.Fatal("Format returned empty string")
	}

	// Verify it's valid JSON
	var parsed types.ExecutionPlan
	err = json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify content
	if parsed.Provider != "node" {
		t.Errorf("expected provider 'node', got '%s'", parsed.Provider)
	}
	if parsed.Runtime.Image != "node:20-alpine" {
		t.Errorf("expected image 'node:20-alpine', got '%s'", parsed.Runtime.Image)
	}
}

func TestJSONFormatter_FormatPretty(t *testing.T) {
	formatter := NewJSONFormatter()

	plan := &types.ExecutionPlan{
		Provider: "python",
		Runtime: types.RuntimeConfig{
			Image: "python:3.11-slim",
		},
		Commands: types.Commands{
			Setup: []string{"pip install -r requirements.txt"},
			Run:   []string{"python app.py"},
		},
		Port: 8000,
	}

	options := &types.CLIOptions{
		Format: "json",
		Pretty: true,
	}

	result, err := formatter.Format(plan, options)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Pretty formatted JSON should contain indentation
	if !strings.Contains(result, "  ") {
		t.Error("Pretty formatted JSON should contain indentation")
	}

	// Verify it's still valid JSON
	var parsed types.ExecutionPlan
	err = json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
}

func TestNewPrettyFormatter(t *testing.T) {
	formatter := NewPrettyFormatter()
	if formatter == nil {
		t.Fatal("NewPrettyFormatter() returned nil")
	}
}

func TestPrettyFormatter_Format(t *testing.T) {
	formatter := NewPrettyFormatter()

	plan := &types.ExecutionPlan{
		Provider: "node",
		Runtime: types.RuntimeConfig{
			Image:     "node:20-alpine",
			Framework: func() *string { s := "express"; return &s }(),
		},
		Environment: map[string]string{
			"NODE_ENV": "production",
			"PORT":     "3000",
		},
		Apt: []string{"curl", "git"},
		Commands: types.Commands{
			Setup: []string{"npm install"},
			Dev:   []string{"npm run dev"},
			Build: []string{"npm run build"},
			Run:   []string{"npm start"},
		},
		Port: 3000,
		Evidence: types.Evidence{
			Files:  []string{"package.json"},
			Reason: "Node.js project detected",
		},
	}

	options := &types.CLIOptions{
		Format:  "pretty",
		Verbose: false,
	}

	result, err := formatter.Format(plan, options)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if result == "" {
		t.Fatal("Format returned empty string")
	}

	// Check for expected sections
	expectedSections := []string{
		"DevBox Pack Execution Plan",
		"Runtime Configuration",
		"Provider: node",
		"Framework: express",
		"Base Image",
		"Image: node:20-alpine",
		"System Dependencies",
		"curl",
		"git",
		"Development Commands",
		"npm run dev",
		"Build Commands",
		"npm run build",
		"Production Commands",
		"npm start",
		"Setup Commands",
		"npm install",
		"Environment Variables",
		"NODE_ENV=production",
		"PORT=3000",
		"Port: 3000",
		"Detection Evidence",
		"package.json",
		"Node.js project detected",
	}

	for _, section := range expectedSections {
		if !strings.Contains(result, section) {
			t.Errorf("expected section '%s' not found in output", section)
		}
	}
}

func TestPrettyFormatter_FormatMinimal(t *testing.T) {
	formatter := NewPrettyFormatter()

	// Minimal plan with only required fields
	plan := &types.ExecutionPlan{
		Provider: "staticfile",
		Runtime: types.RuntimeConfig{
			Image: "nginx:alpine",
		},
		Commands: types.Commands{
			Run: []string{"nginx -g 'daemon off;'"},
		},
		Port: 80,
	}

	options := &types.CLIOptions{
		Format: "pretty",
	}

	result, err := formatter.Format(plan, options)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should contain basic sections even for minimal plan
	expectedSections := []string{
		"DevBox Pack Execution Plan",
		"Runtime Configuration",
		"Provider: staticfile",
		"Base Image",
		"Image: nginx:alpine",
		"Production Commands",
		"nginx -g 'daemon off;'",
		"Port: 80",
	}

	for _, section := range expectedSections {
		if !strings.Contains(result, section) {
			t.Errorf("expected section '%s' not found in output", section)
		}
	}

	// Should not contain sections for empty fields
	unexpectedSections := []string{
		"Framework:",
		"System Dependencies",
		"Environment Variables",
		"Detection Evidence",
	}

	for _, section := range unexpectedSections {
		if strings.Contains(result, section) {
			t.Errorf("unexpected section '%s' found in output", section)
		}
	}
}

func TestNewFormatterFactory(t *testing.T) {
	factory := NewFormatterFactory()
	if factory == nil {
		t.Fatal("NewFormatterFactory() returned nil")
	}

	if factory.formatters == nil {
		t.Fatal("formatters map is nil")
	}

	// Check that default formatters are registered
	supportedFormats := factory.GetSupportedFormats()
	expectedFormats := []string{"json", "pretty"}

	for _, expected := range expectedFormats {
		found := false
		for _, format := range supportedFormats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected format '%s' not found in supported formats", expected)
		}
	}
}

func TestFormatterFactory_GetFormatter(t *testing.T) {
	factory := NewFormatterFactory()

	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"json formatter", "json", false},
		{"pretty formatter", "pretty", false},
		{"invalid formatter", "xml", true},
		{"empty format", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := factory.GetFormatter(tt.format)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				if formatter != nil {
					t.Error("expected nil formatter for error case")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if formatter == nil {
				t.Error("formatter should not be nil")
			}
		})
	}
}

func TestFormatterFactory_Format(t *testing.T) {
	factory := NewFormatterFactory()

	plan := &types.ExecutionPlan{
		Provider: "go",
		Runtime: types.RuntimeConfig{
			Image: "golang:1.21-alpine",
		},
		Commands: types.Commands{
			Setup: []string{"go mod download"},
			Build: []string{"go build -o app"},
			Run:   []string{"./app"},
		},
		Port: 8080,
	}

	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"json format", "json", false},
		{"pretty format", "pretty", false},
		{"invalid format", "xml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := factory.Format(plan, tt.format)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Error("result should not be empty")
			}
		})
	}
}

func TestFormatterFactory_RegisterFormatter(t *testing.T) {
	factory := NewFormatterFactory()

	// Register a custom formatter
	customFormatter := NewJSONFormatter()
	factory.RegisterFormatter("custom", customFormatter)

	// Verify it was registered
	formatter, err := factory.GetFormatter("custom")
	if err != nil {
		t.Fatalf("failed to get custom formatter: %v", err)
	}

	if formatter != customFormatter {
		t.Error("registered formatter does not match")
	}

	// Verify it appears in supported formats
	supportedFormats := factory.GetSupportedFormats()
	found := false
	for _, format := range supportedFormats {
		if format == "custom" {
			found = true
			break
		}
	}
	if !found {
		t.Error("custom formatter not found in supported formats")
	}
}

func TestNewOutputUtils(t *testing.T) {
	utils := NewOutputUtils()
	if utils == nil {
		t.Fatal("NewOutputUtils() returned nil")
	}

	if utils.factory == nil {
		t.Fatal("factory should not be nil")
	}
}

func TestOutputUtils_OutputPlan(t *testing.T) {
	utils := NewOutputUtils()

	plan := &types.ExecutionPlan{
		Provider: "ruby",
		Runtime: types.RuntimeConfig{
			Image: "ruby:3.2-alpine",
		},
		Commands: types.Commands{
			Setup: []string{"bundle install"},
			Run:   []string{"ruby app.rb"},
		},
		Port: 4567,
	}

	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"json output", "json", false},
		{"pretty output", "pretty", false},
		{"invalid format", "xml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &types.CLIOptions{
				Format: tt.format,
			}

			err := utils.OutputPlan(plan, options)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestJSONFormatter_FormatNilPlan(t *testing.T) {
	formatter := NewJSONFormatter()

	_, err := formatter.Format(nil, nil)
	if err == nil {
		t.Error("expected error for nil plan")
	}
}

func TestPrettyFormatter_FormatNilPlan(t *testing.T) {
	formatter := NewPrettyFormatter()

	_, err := formatter.Format(nil, nil)
	if err == nil {
		t.Error("expected error for nil plan")
	}
}
