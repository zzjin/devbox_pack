package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewDenoProvider(t *testing.T) {
	provider := NewDenoProvider()
	AssertProviderBasic(t, provider, "deno", 70)
}

func TestDenoProvider_Detect_NoDenoFiles(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	provider := NewDenoProvider()

	files := []types.FileInfo{
		{Path: "README.md", IsDirectory: false},
		{Path: "main.py", IsDirectory: false},
	}

	result, err := provider.Detect(helper.TempDir, files, helper.GitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	AssertDetectResult(t, result, false, "")
}

func TestDenoProvider_Detect_WithStaticfile(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "Staticfile", IsDirectory: false},
		{Path: "deno.json", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if result.Matched {
		t.Error("expected not matched when Staticfile exists")
	}
}

func TestDenoProvider_Detect_WithGoWork(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "go.work", IsDirectory: false},
		{Path: "deno.json", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if result.Matched {
		t.Error("expected not matched when go.work exists")
	}
}

func TestDenoProvider_Detect_WithDenoJson(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with deno.json
	tempDir, err := os.MkdirTemp("", "deno-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	denoJsonContent := `{
		"tasks": {
			"dev": "deno run --watch main.ts",
			"start": "deno run main.ts"
		},
		"imports": {
			"std/": "https://deno.land/std@0.200.0/"
		}
	}`

	denoJsonPath := filepath.Join(tempDir, "deno.json")
	if err := os.WriteFile(denoJsonPath, []byte(denoJsonContent), 0644); err != nil {
		t.Fatalf("failed to write deno.json: %v", err)
	}

	files := []types.FileInfo{
		{Path: "deno.json", IsDirectory: false},
		{Path: "main.ts", IsDirectory: false},
	}

	result, err := provider.Detect(tempDir, files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Deno project with deno.json")
	}

	if result.Language != "deno" {
		t.Errorf("expected language 'deno', got %s", result.Language)
	}
}

func TestDenoProvider_Detect_WithDenoJsonc(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "deno.jsonc", IsDirectory: false},
		{Path: "main.ts", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Deno project with deno.jsonc")
	}
}

func TestDenoProvider_Detect_WithDenoLock(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "deno.lock", IsDirectory: false},
		{Path: "main.ts", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Deno project with deno.lock")
	}
}

func TestDenoProvider_Detect_WithTypeScriptFiles(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "main.ts", IsDirectory: false},
		{Path: "utils.ts", IsDirectory: false},
		{Path: "mod.ts", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Deno project with TypeScript files")
	}
}

func TestDenoProvider_Detect_WithJavaScriptFiles(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "main.js", IsDirectory: false},
		{Path: "utils.js", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	// Should have low confidence for JS files without Deno-specific files
	if result.Confidence > 0.3 {
		t.Errorf("expected low confidence for JS files, got %f", result.Confidence)
	}
}

func TestDenoProvider_GenerateCommands_BasicProject(t *testing.T) {
	provider := NewDenoProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "deno",
		Framework: "",
		Evidence: types.Evidence{
			Files: []string{"deno.json", "main.ts"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Setup) == 0 {
		t.Error("expected setup commands")
	}

	if len(commands.Dev) == 0 {
		t.Error("expected dev commands")
	}

	if len(commands.Run) == 0 {
		t.Error("expected run commands")
	}

	// Check for Deno-specific commands
	foundDenoRun := false
	for _, cmd := range commands.Run {
		if strings.Contains(cmd, "deno run") {
			foundDenoRun = true
			break
		}
	}
	if !foundDenoRun {
		t.Error("expected 'deno run' command in run commands")
	}
}

func TestDenoProvider_GenerateCommands_WithFramework(t *testing.T) {
	provider := NewDenoProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "deno",
		Framework: "Fresh",
		Evidence: types.Evidence{
			Files: []string{"deno.json", "main.ts"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Dev) == 0 {
		t.Error("expected dev commands")
	}

	// Check for framework-specific commands
	foundFrameworkCmd := false
	for _, cmd := range commands.Dev {
		if strings.Contains(cmd, "deno task") || strings.Contains(cmd, "fresh") {
			foundFrameworkCmd = true
			break
		}
	}
	if !foundFrameworkCmd {
		t.Error("expected framework-specific command in dev commands")
	}
}

func TestDenoProvider_GenerateEnvironment(t *testing.T) {
	provider := NewDenoProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "deno",
		Version:   "1.37.0",
		Framework: "Fresh",
	}

	env := provider.GenerateEnvironment(result)

	if env["DENO_ENV"] != "production" {
		t.Errorf("expected DENO_ENV 'production', got %s", env["DENO_ENV"])
	}

	if env["PORT"] != "8000" {
		t.Errorf("expected PORT '8000', got %s", env["PORT"])
	}

	if env["DENO_VERSION"] != "1.37.0" {
		t.Errorf("expected DENO_VERSION '1.37.0', got %s", env["DENO_VERSION"])
	}
}

func TestDenoProvider_NeedsNativeCompilation(t *testing.T) {
	provider := NewDenoProvider()

	testCases := []struct {
		name     string
		result   *types.DetectResult
		expected bool
	}{
		{
			name: "basic deno project",
			result: &types.DetectResult{
				Matched:  true,
				Language: "deno",
			},
			expected: false,
		},
		{
			name: "deno with native modules",
			result: &types.DetectResult{
				Matched:  true,
				Language: "deno",
				Metadata: map[string]interface{}{
					"hasNativeModules": true,
				},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := provider.NeedsNativeCompilation(tc.result)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDenoProvider_DetectFramework(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	testCases := []struct {
		name              string
		denoJsonContent   string
		expectedFramework string
	}{
		{
			name: "Fresh framework",
			denoJsonContent: `{
				"imports": {
					"fresh": "https://deno.land/x/fresh@1.4.0/mod.ts"
				}
			}`,
			expectedFramework: "Fresh",
		},
		{
			name: "Oak framework",
			denoJsonContent: `{
				"imports": {
					"oak": "https://deno.land/x/oak@v12.6.0/mod.ts"
				}
			}`,
			expectedFramework: "Oak",
		},
		{
			name: "Hono framework",
			denoJsonContent: `{
				"imports": {
					"hono": "https://deno.land/x/hono@v3.5.0/mod.ts"
				}
			}`,
			expectedFramework: "Hono",
		},
		{
			name: "no framework",
			denoJsonContent: `{
				"imports": {
					"std/": "https://deno.land/std@0.200.0/"
				}
			}`,
			expectedFramework: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "deno-framework-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Write deno.json
			denoJsonPath := filepath.Join(tempDir, "deno.json")
			if err := os.WriteFile(denoJsonPath, []byte(tc.denoJsonContent), 0644); err != nil {
				t.Fatalf("failed to write deno.json: %v", err)
			}

			framework, err := provider.detectFramework(tempDir, gitHandler)
			if err != nil {
				t.Fatalf("detectFramework failed: %v", err)
			}

			if framework != tc.expectedFramework {
				t.Errorf("expected framework %s, got %s", tc.expectedFramework, framework)
			}
		})
	}
}

func TestDenoProvider_DetectDenoVersion(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with .dvmrc
	tempDir, err := os.MkdirTemp("", "deno-version-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .dvmrc file
	dvmrcContent := "1.37.0"
	dvmrcPath := filepath.Join(tempDir, ".dvmrc")
	if err := os.WriteFile(dvmrcPath, []byte(dvmrcContent), 0644); err != nil {
		t.Fatalf("failed to write .dvmrc: %v", err)
	}

	version, err := provider.detectDenoVersion(tempDir, gitHandler)
	if err != nil {
		t.Fatalf("detectDenoVersion failed: %v", err)
	}

	if version.Version != "1.37.0" {
		t.Errorf("expected version '1.37.0', got %s", version.Version)
	}

	if version.Source != ".dvmrc" {
		t.Errorf("expected source '.dvmrc', got %s", version.Source)
	}
}

func TestDenoProvider_DetectDenoVersion_FromDenoJson(t *testing.T) {
	provider := NewDenoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with deno.json
	tempDir, err := os.MkdirTemp("", "deno-version-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create deno.json with version
	denoJsonContent := `{
		"version": "1.36.0",
		"tasks": {
			"start": "deno run main.ts"
		}
	}`
	denoJsonPath := filepath.Join(tempDir, "deno.json")
	if err := os.WriteFile(denoJsonPath, []byte(denoJsonContent), 0644); err != nil {
		t.Fatalf("failed to write deno.json: %v", err)
	}

	version, err := provider.detectDenoVersion(tempDir, gitHandler)
	if err != nil {
		t.Fatalf("detectDenoVersion failed: %v", err)
	}

	if version.Version != "1.36.0" {
		t.Errorf("expected version '1.36.0', got %s", version.Version)
	}

	if version.Source != "deno.json" {
		t.Errorf("expected source 'deno.json', got %s", version.Source)
	}
}