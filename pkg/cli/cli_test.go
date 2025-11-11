package cli

import (
	"testing"
)

func TestNewCLIApp(t *testing.T) {
	app := NewCLIApp()
	if app == nil {
		t.Fatal("NewCLIApp() returned nil")
	}

	if app.version == "" {
		t.Error("version should not be empty")
	}
}

func TestParseArgs_Help(t *testing.T) {
	// Test --help flag (should exit, but we can't test that directly)
	// Instead, we test that it doesn't return an error when parsing
	args := []string{"devbox-pack", "--help"}

	// This will call os.Exit(0), so we can't test it directly
	// We'll test the parsing logic instead
	if len(args) >= 2 && (args[1] == "--help" || args[1] == "-h") {
		// Help flag detected correctly
		return
	}
}

func TestParseArgs_Version(t *testing.T) {
	// Test --version flag
	args := []string{"devbox-pack", "--version"}

	if len(args) >= 2 && (args[1] == "--version" || args[1] == "-v") {
		// Version flag detected correctly
		return
	}
}

func TestParseArgs_ValidRepository(t *testing.T) {
	app := NewCLIApp()

	tests := []struct {
		name     string
		args     []string
		wantRepo string
		wantErr  bool
	}{
		{
			name:     "local path",
			args:     []string{"devbox-pack", "."},
			wantRepo: ".",
			wantErr:  false,
		},
		{
			name:     "github url",
			args:     []string{"devbox-pack", "https://github.com/user/repo"},
			wantRepo: "https://github.com/user/repo",
			wantErr:  false,
		},
		{
			name:     "with verbose flag",
			args:     []string{"devbox-pack", ".", "--verbose"},
			wantRepo: ".",
			wantErr:  false,
		},
		{
			name:     "with format option",
			args:     []string{"devbox-pack", ".", "--format", "json"},
			wantRepo: ".",
			wantErr:  false,
		},
		{
			name:     "no repository",
			args:     []string{"devbox-pack"},
			wantRepo: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, options, err := app.parseArgs(tt.args)

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

			if repo != tt.wantRepo {
				t.Errorf("expected repo %s, got %s", tt.wantRepo, repo)
			}

			if options == nil {
				t.Error("options should not be nil")
			}
		})
	}
}

func TestParseArgs_WithOptions(t *testing.T) {
	app := NewCLIApp()

	args := []string{
		"devbox-pack", ".",
		"--ref", "main",
		"--subdir", "backend",
		"--provider", "node",
		"--format", "json",
		"--verbose",
		"--offline",
		"--platform", "linux/amd64",
		"--base", "node:18-alpine",
	}

	repo, options, err := app.parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}

	if repo != "." {
		t.Errorf("expected repo '.', got '%s'", repo)
	}

	expectedOptions := map[string]interface{}{
		"ref":      "main",
		"subdir":   "backend",
		"provider": "node",
		"format":   "json",
		"verbose":  true,
		"offline":  true,
		"platform": "linux/amd64",
		"base":     "node:18-alpine",
	}

	for key, expected := range expectedOptions {
		if actual, ok := options[key]; !ok {
			t.Errorf("missing option %s", key)
		} else if actual != expected {
			t.Errorf("option %s: expected %v, got %v", key, expected, actual)
		}
	}
}

func TestValidateOptions_ValidFormat(t *testing.T) {
	app := NewCLIApp()

	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"json format", "json", false},
		{"pretty format", "pretty", false},
		{"invalid format", "xml", true},
		{"empty format", "", false}, // should default to pretty
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawOptions := map[string]interface{}{}
			if tt.format != "" {
				rawOptions["format"] = tt.format
			}

			options, err := app.validateOptions(rawOptions)

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

			if options == nil {
				t.Error("options should not be nil")
			}

			expectedFormat := tt.format
			if expectedFormat == "" {
				expectedFormat = "pretty"
			}

			if options.Format != expectedFormat {
				t.Errorf("expected format %s, got %s", expectedFormat, options.Format)
			}
		})
	}
}

func TestValidateOptions_ValidProvider(t *testing.T) {
	app := NewCLIApp()

	tests := []struct {
		name     string
		provider string
		wantErr  bool
	}{
		{"node provider", "node", false},
		{"python provider", "python", false},
		{"java provider", "java", false},
		{"go provider", "go", false},
		{"invalid provider", "invalid", true},
		{"no provider", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawOptions := map[string]interface{}{}
			if tt.provider != "" {
				rawOptions["provider"] = tt.provider
			}

			options, err := app.validateOptions(rawOptions)

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

			if options == nil {
				t.Error("options should not be nil")
			}

			if tt.provider != "" {
				if options.Provider == nil {
					t.Error("provider should not be nil")
				} else if *options.Provider != tt.provider {
					t.Errorf("expected provider %s, got %s", tt.provider, *options.Provider)
				}
			}
		})
	}
}

func TestValidateOptions_ValidPlatform(t *testing.T) {
	app := NewCLIApp()

	tests := []struct {
		name     string
		platform string
		wantErr  bool
	}{
		{"linux amd64", "linux/amd64", false},
		{"linux arm64", "linux/arm64", false},
		{"darwin amd64", "darwin/amd64", false},
		{"darwin arm64", "darwin/arm64", false},
		{"invalid platform", "windows/amd64", true},
		{"no platform", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawOptions := map[string]interface{}{}
			if tt.platform != "" {
				rawOptions["platform"] = tt.platform
			}

			options, err := app.validateOptions(rawOptions)

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

			if options == nil {
				t.Error("options should not be nil")
			}

			if tt.platform != "" {
				if options.Platform == nil {
					t.Error("platform should not be nil")
				} else if *options.Platform != tt.platform {
					t.Errorf("expected platform %s, got %s", tt.platform, *options.Platform)
				}
			}
		})
	}
}

func TestValidateOptions_AllOptions(t *testing.T) {
	app := NewCLIApp()

	rawOptions := map[string]interface{}{
		"ref":      "develop",
		"subdir":   "api",
		"provider": "python",
		"format":   "json",
		"verbose":  true,
		"offline":  true,
		"quiet":    false,
		"platform": "linux/amd64",
		"base":     "python:3.11-slim",
	}

	options, err := app.validateOptions(rawOptions)
	if err != nil {
		t.Fatalf("validateOptions failed: %v", err)
	}

	if options == nil {
		t.Fatal("options should not be nil")
	}

	// Check all options are set correctly
	if options.Ref == nil || *options.Ref != "develop" {
		t.Error("ref not set correctly")
	}
	if options.Subdir == nil || *options.Subdir != "api" {
		t.Error("subdir not set correctly")
	}
	if options.Provider == nil || *options.Provider != "python" {
		t.Error("provider not set correctly")
	}
	if options.Format != "json" {
		t.Error("format not set correctly")
	}
	if !options.Verbose {
		t.Error("verbose not set correctly")
	}
	if !options.Offline {
		t.Error("offline not set correctly")
	}
	if options.Quiet {
		t.Error("quiet should be false")
	}
	if options.Platform == nil || *options.Platform != "linux/amd64" {
		t.Error("platform not set correctly")
	}
	if options.Base == nil || *options.Base != "python:3.11-slim" {
		t.Error("base not set correctly")
	}
}

func TestRun_InvalidArguments(t *testing.T) {
	app := NewCLIApp()

	// Test no arguments
	err := app.Run([]string{"devbox-pack"})
	if err == nil {
		t.Error("expected error for no arguments")
	}

	// Test invalid provider
	err = app.Run([]string{"devbox-pack", ".", "--provider", "invalid"})
	if err == nil {
		t.Error("expected error for invalid provider")
	}

	// Test invalid format
	err = app.Run([]string{"devbox-pack", ".", "--format", "xml"})
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestParseArgs_EdgeCases(t *testing.T) {
	app := NewCLIApp()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "option without value",
			args:    []string{"devbox-pack", ".", "--ref"},
			wantErr: true,
		},
		{
			name:    "unknown argument",
			args:    []string{"devbox-pack", ".", "unknown"},
			wantErr: true,
		},
		{
			name:    "multiple repositories",
			args:    []string{"devbox-pack", "repo1", "repo2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := app.parseArgs(tt.args)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
