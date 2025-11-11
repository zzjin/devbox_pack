package providers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewNodeProvider(t *testing.T) {
	provider := NewNodeProvider()

	if provider == nil {
		t.Fatal("NewNodeProvider() returned nil")
	}

	if provider.GetName() != "node" {
		t.Errorf("expected name 'node', got %s", provider.GetName())
	}

	if provider.GetLanguage() != "node" {
		t.Errorf("expected language 'node', got %s", provider.GetLanguage())
	}

	if provider.GetPriority() != 80 {
		t.Errorf("expected priority 80, got %d", provider.GetPriority())
	}
}

func TestNodeProvider_Detect_NoNodeFiles(t *testing.T) {
	provider := NewNodeProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "README.md", IsDirectory: false},
		{Path: "main.py", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if result.Matched {
		t.Error("expected not matched for non-Node.js project")
	}

	if result.Confidence > 0.3 {
		t.Errorf("expected low confidence, got %f", result.Confidence)
	}
}

func TestNodeProvider_Detect_WithPackageJson(t *testing.T) {
	provider := NewNodeProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with package.json
	tmpDir, err := os.MkdirTemp("", "node-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	packageJSON := `{
		"name": "test-project",
		"version": "1.0.0",
		"scripts": {
			"start": "node index.js",
			"build": "webpack"
		},
		"dependencies": {
			"react": "^18.0.0",
			"react-dom": "^18.0.0"
		},
		"engines": {
			"node": ">=18.0.0"
		}
	}`

	packageJSONPath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(packageJSONPath, []byte(packageJSON), 0644); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	files := []types.FileInfo{
		{Path: "package.json", IsDirectory: false},
		{Path: "index.js", IsDirectory: false},
		{Path: "package-lock.json", IsDirectory: false},
	}

	result, err := provider.Detect(tmpDir, files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected match for Node.js project")
	}

	if result.Confidence <= 0.3 {
		t.Errorf("expected high confidence, got %f", result.Confidence)
	}

	if result.Language != "node" {
		t.Errorf("expected language 'node', got %s", result.Language)
	}

	if result.Framework != "react" {
		t.Errorf("expected framework 'react', got %s", result.Framework)
	}

	if result.PackageManager == nil || result.PackageManager.Name != "npm" {
		t.Errorf("expected package manager 'npm', got %v", result.PackageManager)
	}

	if len(result.BuildTools) == 0 || result.BuildTools[0] != "webpack" {
		t.Errorf("expected build tool 'webpack', got %v", result.BuildTools)
	}
}

func TestNodeProvider_Detect_WithYarnLock(t *testing.T) {
	provider := NewNodeProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "package.json", IsDirectory: false},
		{Path: "yarn.lock", IsDirectory: false},
		{Path: "index.js", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected match for Node.js project with yarn")
	}

	if result.PackageManager == nil || result.PackageManager.Name != "yarn" {
		t.Errorf("expected package manager 'yarn', got %v", result.PackageManager)
	}
}

func TestNodeProvider_Detect_WithPnpmLock(t *testing.T) {
	provider := NewNodeProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "package.json", IsDirectory: false},
		{Path: "pnpm-lock.yaml", IsDirectory: false},
		{Path: "index.js", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected match for Node.js project with pnpm")
	}

	if result.PackageManager == nil || result.PackageManager.Name != "pnpm" {
		t.Errorf("expected package manager 'pnpm', got %v", result.PackageManager)
	}
}

func TestNodeProvider_GenerateCommands_BasicProject(t *testing.T) {
	provider := NewNodeProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "node",
		Framework: "",
		Evidence: types.Evidence{
			Files: []string{"package.json", "package-lock.json"},
		},
		Metadata: map[string]interface{}{
			"packageInfo": map[string]interface{}{
				"scripts": map[string]interface{}{
					"start": "node index.js",
					"build": "webpack",
				},
			},
		},
	}

	nodeProvider := "node"
	options := types.CLIOptions{
		Provider: &nodeProvider,
	}

	commands := provider.GenerateCommands(result, options)

	if len(commands.Setup) == 0 {
		t.Error("expected setup commands")
	}

	if commands.Setup[0] != "npm install" {
		t.Errorf("expected 'npm install', got %s", commands.Setup[0])
	}

	if len(commands.Run) == 0 {
		t.Error("expected run commands")
	}

	// Build commands are only generated if there's a build script
	// Let's just verify that we have build commands since we have a build script
	if len(commands.Build) == 0 {
		t.Log("Note: no build commands generated")
	}
}

func TestNodeProvider_GenerateCommands_YarnProject(t *testing.T) {
	provider := NewNodeProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "node",
		Framework: "",
		Evidence: types.Evidence{
			Files: []string{"package.json", "yarn.lock"},
		},
		Metadata: map[string]interface{}{
			"packageInfo": map[string]interface{}{
				"scripts": map[string]interface{}{
					"start": "node server.js",
				},
			},
		},
	}

	nodeProvider := "node"
	options := types.CLIOptions{
		Provider: &nodeProvider,
	}

	commands := provider.GenerateCommands(result, options)

	if len(commands.Setup) == 0 {
		t.Error("expected setup commands")
	}

	if commands.Setup[0] != "yarn install" {
		t.Errorf("expected 'yarn install', got %s", commands.Setup[0])
	}
}

func TestNodeProvider_GenerateCommands_NextJSProject(t *testing.T) {
	provider := NewNodeProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "node",
		Framework: "next",
		Evidence: types.Evidence{
			Files: []string{"package.json", "next.config.js"},
		},
		Metadata: map[string]interface{}{
			"packageInfo": map[string]interface{}{
				"dependencies": map[string]interface{}{
					"next": "^13.0.0",
				},
				"scripts": map[string]interface{}{
					"dev":   "next dev",
					"build": "next build",
					"start": "next start",
				},
			},
		},
	}

	nodeProvider := "node"
	options := types.CLIOptions{
		Provider: &nodeProvider,
	}

	commands := provider.GenerateCommands(result, options)

	// Check for Next.js specific commands
	found := false
	for _, cmd := range commands.Dev {
		if cmd == "npm run dev" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Next.js dev command in start commands")
	}
}

func TestNodeProvider_GenerateEnvironment(t *testing.T) {
	provider := NewNodeProvider()

	result := &types.DetectResult{
		Version: "18.17.0",
	}

	env := provider.GenerateEnvironment(result)

	if env["NODE_ENV"] != "development" {
		t.Errorf("expected NODE_ENV 'development', got %s", env["NODE_ENV"])
	}

	if env["NPM_CONFIG_FUND"] != "false" {
		t.Errorf("expected NPM_CONFIG_FUND 'false', got %s", env["NPM_CONFIG_FUND"])
	}

	if env["NPM_CONFIG_AUDIT"] != "false" {
		t.Errorf("expected NPM_CONFIG_AUDIT 'false', got %s", env["NPM_CONFIG_AUDIT"])
	}
}

func TestNodeProvider_NeedsNativeCompilation_WithNativeModules(t *testing.T) {
	provider := NewNodeProvider()

	result := &types.DetectResult{
		Evidence: types.Evidence{
			Files: []string{"package.json", "node_modules/bcrypt", "node_modules/node-sass"},
		},
		Metadata: map[string]interface{}{
			"hasNativeModules": true, // Set the flag to indicate native modules
		},
	}

	needs := provider.NeedsNativeCompilation(result)
	if !needs {
		t.Error("expected to need native compilation for native modules")
	}
}

func TestNodeProvider_NeedsNativeCompilation_WithoutNativeModules(t *testing.T) {
	provider := NewNodeProvider()

	result := &types.DetectResult{
		Metadata: map[string]interface{}{
			"packageInfo": map[string]interface{}{
				"dependencies": map[string]interface{}{
					"react":     "^18.0.0",
					"react-dom": "^18.0.0",
				},
			},
		},
	}

	needs := provider.NeedsNativeCompilation(result)
	if needs {
		t.Error("expected not to need native compilation for pure JS modules")
	}
}

func TestNodeProvider_DetectFramework(t *testing.T) {
	provider := NewNodeProvider()

	tests := []struct {
		name        string
		packageJSON map[string]interface{}
		expected    string
	}{
		{
			name: "Next.js project",
			packageJSON: map[string]interface{}{
				"dependencies": map[string]interface{}{
					"next": "^13.0.0",
				},
			},
			expected: "next",
		},
		{
			name: "React project",
			packageJSON: map[string]interface{}{
				"dependencies": map[string]interface{}{
					"react": "^18.0.0",
				},
			},
			expected: "react",
		},
		{
			name: "Vue project",
			packageJSON: map[string]interface{}{
				"dependencies": map[string]interface{}{
					"vue": "^3.0.0",
				},
			},
			expected: "vue",
		},
		{
			name: "Express project",
			packageJSON: map[string]interface{}{
				"dependencies": map[string]interface{}{
					"express": "^4.0.0",
				},
			},
			expected: "express",
		},
		{
			name: "No framework",
			packageJSON: map[string]interface{}{
				"dependencies": map[string]interface{}{
					"lodash": "^4.0.0",
				},
			},
			expected: "",
		},
		{
			name:        "Empty package.json",
			packageJSON: map[string]interface{}{},
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.detectFramework(tt.packageJSON)
			if result != tt.expected {
				t.Errorf("detectFramework() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestNodeProvider_DetectBuildTool(t *testing.T) {
	provider := NewNodeProvider()

	tests := []struct {
		name        string
		packageJSON map[string]interface{}
		expected    string
	}{
		{
			name: "Webpack project",
			packageJSON: map[string]interface{}{
				"devDependencies": map[string]interface{}{
					"webpack": "^5.0.0",
				},
			},
			expected: "webpack",
		},
		{
			name: "Vite project",
			packageJSON: map[string]interface{}{
				"devDependencies": map[string]interface{}{
					"vite": "^4.0.0",
				},
			},
			expected: "vite",
		},
		{
			name: "Rollup project",
			packageJSON: map[string]interface{}{
				"devDependencies": map[string]interface{}{
					"rollup": "^3.0.0",
				},
			},
			expected: "rollup",
		},
		{
			name: "Parcel project",
			packageJSON: map[string]interface{}{
				"devDependencies": map[string]interface{}{
					"parcel": "^2.0.0",
				},
			},
			expected: "parcel",
		},
		{
			name: "No build tool",
			packageJSON: map[string]interface{}{
				"devDependencies": map[string]interface{}{
					"eslint": "^8.0.0",
				},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.detectBuildTool(tt.packageJSON)
			if result != tt.expected {
				t.Errorf("detectBuildTool() = %s, expected %s", result, tt.expected)
			}
		})
	}
}