package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewGoProvider(t *testing.T) {
	provider := NewGoProvider()
	AssertProviderBasic(t, provider, "go", 20)
}

func TestGoProvider_Detect_NoGoFiles(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	provider := NewGoProvider()

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

func TestGoProvider_Detect_WithGoMod(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	provider := NewGoProvider()

	// Use test data template and modify for specific test
	files := CreateTestFiles(helper, GoTestData.Files)
	files = append(files, types.FileInfo{Path: "main.go", IsDirectory: false})

	result, err := provider.Detect(helper.TempDir, files, helper.GitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	AssertDetectResult(t, result, true, "go")

	if result.Version != "1.21" {
		t.Errorf("expected version '1.21', got %s", result.Version)
	}
}

func TestGoProvider_Detect_WithGoWork(t *testing.T) {
	provider := NewGoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with go.work
	tempDir, err := os.MkdirTemp("", "go-work-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	goWorkContent := `go 1.21

use (
	./module1
	./module2
)
`

	goWorkPath := filepath.Join(tempDir, "go.work")
	if err := os.WriteFile(goWorkPath, []byte(goWorkContent), 0644); err != nil {
		t.Fatalf("failed to write go.work: %v", err)
	}

	files := []types.FileInfo{
		{Path: "go.work", IsDirectory: false},
		{Path: "module1/", IsDirectory: true},
		{Path: "module2/", IsDirectory: true},
	}

	result, err := provider.Detect(tempDir, files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Go workspace with go.work")
	}

	if result.Language != "go" {
		t.Errorf("expected language 'go', got %s", result.Language)
	}
}

func TestGoProvider_Detect_WithGoFiles(t *testing.T) {
	provider := NewGoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "main.go", IsDirectory: false},
		{Path: "utils.go", IsDirectory: false},
		{Path: "go.sum", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Go project with .go files")
	}
}

func TestGoProvider_Detect_WithVendor(t *testing.T) {
	provider := NewGoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "main.go", IsDirectory: false},
		{Path: "vendor/", IsDirectory: true},
		{Path: "vendor/modules.txt", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Go project with vendor directory")
	}
}

func TestGoProvider_GenerateCommands_BasicProject(t *testing.T) {
	provider := NewGoProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "go",
		Framework: "",
		Evidence: types.Evidence{
			Files: []string{"go.mod", "main.go"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Setup) == 0 {
		t.Error("expected setup commands")
	}

	if len(commands.Build) == 0 {
		t.Error("expected build commands")
	}

	if len(commands.Run) == 0 {
		t.Error("expected run commands")
	}

	// Check for Go-specific commands
	foundGoBuild := false
	for _, cmd := range commands.Build {
		if strings.Contains(cmd, "go build") {
			foundGoBuild = true
			break
		}
	}
	if !foundGoBuild {
		t.Error("expected 'go build' command in build commands")
	}

	foundGoRun := false
	for _, cmd := range commands.Run {
		if strings.Contains(cmd, "go run") || strings.Contains(cmd, "./") {
			foundGoRun = true
			break
		}
	}
	if !foundGoRun {
		t.Error("expected 'go run' or executable command in run commands")
	}
}

func TestGoProvider_GenerateCommands_WithFramework(t *testing.T) {
	provider := NewGoProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "go",
		Framework: "Gin",
		Evidence: types.Evidence{
			Files: []string{"go.mod", "main.go"},
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
		if strings.Contains(cmd, "go run") || strings.Contains(cmd, "air") {
			foundFrameworkCmd = true
			break
		}
	}
	if !foundFrameworkCmd {
		t.Error("expected framework-specific command in dev commands")
	}
}

func TestGoProvider_GenerateEnvironment(t *testing.T) {
	provider := NewGoProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "go",
		Version:   "1.21",
		Framework: "Gin",
	}

	env := provider.GenerateEnvironment(result)

	if env["GO_ENV"] != "production" {
		t.Errorf("expected GO_ENV 'production', got %s", env["GO_ENV"])
	}

	if env["PORT"] != "8080" {
		t.Errorf("expected PORT '8080', got %s", env["PORT"])
	}

	if env["CGO_ENABLED"] != "0" {
		t.Errorf("expected CGO_ENABLED '0', got %s", env["CGO_ENABLED"])
	}
}

func TestGoProvider_NeedsNativeCompilation(t *testing.T) {
	provider := NewGoProvider()

	testCases := []struct {
		name     string
		result   *types.DetectResult
		expected bool
	}{
		{
			name: "basic go project",
			result: &types.DetectResult{
				Matched:  true,
				Language: "go",
			},
			expected: false,
		},
		{
			name: "go with cgo",
			result: &types.DetectResult{
				Matched:  true,
				Language: "go",
				Metadata: map[string]interface{}{
					"usesCGO": true,
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

func TestGoProvider_DetectFramework(t *testing.T) {
	provider := NewGoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	testCases := []struct {
		name              string
		goModContent      string
		expectedFramework string
	}{
		{
			name: "Gin framework",
			goModContent: `module example.com/myapp

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)
`,
			expectedFramework: "Gin",
		},
		{
			name: "Echo framework",
			goModContent: `module example.com/myapp

go 1.21

require (
	github.com/labstack/echo/v4 v4.11.1
)
`,
			expectedFramework: "Echo",
		},
		{
			name: "Fiber framework",
			goModContent: `module example.com/myapp

go 1.21

require (
	github.com/gofiber/fiber/v2 v2.48.0
)
`,
			expectedFramework: "Fiber",
		},
		{
			name: "Beego framework",
			goModContent: `module example.com/myapp

go 1.21

require (
	github.com/beego/beego/v2 v2.0.7
)
`,
			expectedFramework: "Beego",
		},
		{
			name: "no framework",
			goModContent: `module example.com/myapp

go 1.21

require (
	github.com/stretchr/testify v1.8.4
)
`,
			expectedFramework: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "go-framework-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Write go.mod
			goModPath := filepath.Join(tempDir, "go.mod")
			if err := os.WriteFile(goModPath, []byte(tc.goModContent), 0644); err != nil {
				t.Fatalf("failed to write go.mod: %v", err)
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

func TestGoProvider_DetectGoVersion(t *testing.T) {
	provider := NewGoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	testCases := []struct {
		name            string
		goModContent    string
		expectedVersion string
		expectedSource  string
	}{
		{
			name: "go.mod with version",
			goModContent: `module example.com/myapp

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)
`,
			expectedVersion: "1.21",
			expectedSource:  "go.mod",
		},
		{
			name: "go.mod with patch version",
			goModContent: `module example.com/myapp

go 1.20.5

require (
	github.com/stretchr/testify v1.8.4
)
`,
			expectedVersion: "1.20.5",
			expectedSource:  "go.mod",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "go-version-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Write go.mod
			goModPath := filepath.Join(tempDir, "go.mod")
			if err := os.WriteFile(goModPath, []byte(tc.goModContent), 0644); err != nil {
				t.Fatalf("failed to write go.mod: %v", err)
			}

			version, err := provider.detectGoVersion(tempDir, gitHandler)
			if err != nil {
				t.Fatalf("detectGoVersion failed: %v", err)
			}

			if version.Version != tc.expectedVersion {
				t.Errorf("expected version %s, got %s", tc.expectedVersion, version.Version)
			}

			if version.Source != tc.expectedSource {
				t.Errorf("expected source %s, got %s", tc.expectedSource, version.Source)
			}
		})
	}
}

func TestGoProvider_DetectGoVersion_FromGoWork(t *testing.T) {
	provider := NewGoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with go.work
	tempDir, err := os.MkdirTemp("", "go-work-version-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create go.work with version
	goWorkContent := `go 1.21

use (
	./module1
	./module2
)
`
	goWorkPath := filepath.Join(tempDir, "go.work")
	if err := os.WriteFile(goWorkPath, []byte(goWorkContent), 0644); err != nil {
		t.Fatalf("failed to write go.work: %v", err)
	}

	version, err := provider.detectGoVersion(tempDir, gitHandler)
	if err != nil {
		t.Fatalf("detectGoVersion failed: %v", err)
	}

	if version.Version != "1.21" {
		t.Errorf("expected version '1.21', got %s", version.Version)
	}

	if version.Source != "go.work" {
		t.Errorf("expected source 'go.work', got %s", version.Source)
	}
}

func TestGoProvider_DetectGoVersion_NoVersionFile(t *testing.T) {
	provider := NewGoProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory without version files
	tempDir, err := os.MkdirTemp("", "go-no-version-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	version, err := provider.detectGoVersion(tempDir, gitHandler)
	if err != nil {
		t.Fatalf("detectGoVersion failed: %v", err)
	}

	if version.Version != "latest" {
		t.Errorf("expected version 'latest', got %s", version.Version)
	}

	if version.Source != "default" {
		t.Errorf("expected source 'default', got %s", version.Source)
	}
}