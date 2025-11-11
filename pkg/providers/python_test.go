package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewPythonProvider(t *testing.T) {
	provider := NewPythonProvider()

	if provider == nil {
		t.Fatal("NewPythonProvider() returned nil")
	}

	if provider.GetName() != "python" {
		t.Errorf("expected name 'python', got %s", provider.GetName())
	}

	if provider.GetPriority() != 60 {
		t.Errorf("expected priority 60, got %d", provider.GetPriority())
	}
}

func TestPythonProvider_Detect_NoPythonFiles(t *testing.T) {
	provider := NewPythonProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "README.md", IsDirectory: false},
		{Path: "main.js", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if result.Matched {
		t.Error("expected not matched for non-Python project")
	}
}

func TestPythonProvider_Detect_WithRequirements(t *testing.T) {
	provider := NewPythonProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with requirements.txt
	tempDir, err := os.MkdirTemp("", "python-requirements-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	requirementsContent := `Django==4.2.0
requests==2.31.0
pytest==7.4.0
gunicorn==21.2.0
`

	requirementsPath := filepath.Join(tempDir, "requirements.txt")
	if err := os.WriteFile(requirementsPath, []byte(requirementsContent), 0644); err != nil {
		t.Fatalf("failed to write requirements.txt: %v", err)
	}

	files := []types.FileInfo{
		{Path: "requirements.txt", IsDirectory: false},
		{Path: "main.py", IsDirectory: false},
		{Path: "app/", IsDirectory: true},
	}

	result, err := provider.Detect(tempDir, files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Python project with requirements.txt")
	}

	if result.Language != "python" {
		t.Errorf("expected language 'python', got %s", result.Language)
	}

	if result.PackageManager == nil || result.PackageManager.Name != "pip" {
		t.Errorf("expected package manager 'pip', got %s", func() string {
			if result.PackageManager == nil {
				return "nil"
			}
			return result.PackageManager.Name
		}())
	}
}

func TestPythonProvider_Detect_WithPyprojectToml(t *testing.T) {
	provider := NewPythonProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with pyproject.toml
	tempDir, err := os.MkdirTemp("", "python-pyproject-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	pyprojectContent := `[build-system]
requires = ["poetry-core"]
build-backend = "poetry.core.masonry.api"

[tool.poetry]
name = "my-app"
version = "0.1.0"
description = ""
authors = ["Your Name <you@example.com>"]

[tool.poetry.dependencies]
python = "^3.9"
fastapi = "^0.100.0"
uvicorn = "^0.23.0"

[tool.poetry.group.dev.dependencies]
pytest = "^7.4.0"
`

	pyprojectPath := filepath.Join(tempDir, "pyproject.toml")
	if err := os.WriteFile(pyprojectPath, []byte(pyprojectContent), 0644); err != nil {
		t.Fatalf("failed to write pyproject.toml: %v", err)
	}

	files := []types.FileInfo{
		{Path: "pyproject.toml", IsDirectory: false},
		{Path: "main.py", IsDirectory: false},
	}

	result, err := provider.Detect(tempDir, files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Python project with pyproject.toml")
	}

	if result.PackageManager == nil || result.PackageManager.Name != "pip" {
		t.Errorf("expected package manager 'pip', got %s", func() string {
			if result.PackageManager == nil {
				return "nil"
			}
			return result.PackageManager.Name
		}())
	}
}

func TestPythonProvider_Detect_WithPipfile(t *testing.T) {
	provider := NewPythonProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "Pipfile", IsDirectory: false},
		{Path: "Pipfile.lock", IsDirectory: false},
		{Path: "main.py", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Python project with Pipfile")
	}

	if result.PackageManager == nil || result.PackageManager.Name != "pipenv" {
		t.Errorf("expected package manager 'pipenv', got %s", func() string {
			if result.PackageManager == nil {
				return "nil"
			}
			return result.PackageManager.Name
		}())
	}
}

func TestPythonProvider_Detect_WithSetupPy(t *testing.T) {
	provider := NewPythonProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "setup.py", IsDirectory: false},
		{Path: "src/", IsDirectory: true},
		{Path: "src/mypackage/", IsDirectory: true},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Python project with setup.py")
	}
}

func TestPythonProvider_Detect_WithDjango(t *testing.T) {
	provider := NewPythonProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "manage.py", IsDirectory: false},
		{Path: "requirements.txt", IsDirectory: false},
		{Path: "myproject/", IsDirectory: true},
		{Path: "myproject/settings.py", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Django project")
	}
}

func TestPythonProvider_GenerateCommands_BasicProject(t *testing.T) {
	provider := NewPythonProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "python",
		PackageManager: &types.PackageManager{Name: "pip"},
		Framework:      "",
		Evidence: types.Evidence{
			Files: []string{"requirements.txt", "main.py"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Setup) == 0 {
		t.Error("expected setup commands")
	}

	if len(commands.Run) == 0 {
		t.Error("expected run commands")
	}

	// Check for pip install command
	foundPipInstall := false
	for _, cmd := range commands.Setup {
		if strings.Contains(cmd, "pip install") {
			foundPipInstall = true
			break
		}
	}
	if !foundPipInstall {
		t.Error("expected 'pip install' command in setup commands")
	}

	// Check for python run command
	foundPythonRun := false
	for _, cmd := range commands.Run {
		if strings.Contains(cmd, "python") {
			foundPythonRun = true
			break
		}
	}
	if !foundPythonRun {
		t.Error("expected 'python' command in run commands")
	}
}

func TestPythonProvider_GenerateCommands_DjangoProject(t *testing.T) {
	provider := NewPythonProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "python",
		PackageManager: &types.PackageManager{Name: "pip"},
		Framework:      "Django",
		Evidence: types.Evidence{
			Files: []string{"requirements.txt", "manage.py"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Dev) == 0 {
		t.Error("expected dev commands")
	}

	// Check for Django-specific commands
	foundDjangoRunserver := false
	for _, cmd := range commands.Dev {
		if strings.Contains(cmd, "python manage.py runserver") {
			foundDjangoRunserver = true
			break
		}
	}
	if !foundDjangoRunserver {
		t.Error("expected 'python manage.py runserver' command in dev commands")
	}
}

func TestPythonProvider_GenerateCommands_PoetryProject(t *testing.T) {
	provider := NewPythonProvider()

	result := &types.DetectResult{
		Matched:   true,
		Language:  "python",
		Framework: "FastAPI",
		Evidence: types.Evidence{
			Files: []string{"pyproject.toml", "main.py"},
		},
		Metadata: map[string]interface{}{
			"packageManager": "poetry",
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Setup) == 0 {
		t.Error("expected setup commands")
	}

	// Check for Poetry-specific commands
	foundPoetryInstall := false
	for _, cmd := range commands.Setup {
		if strings.Contains(cmd, "poetry install") {
			foundPoetryInstall = true
			break
		}
	}
	if !foundPoetryInstall {
		t.Error("expected 'poetry install' command in setup commands")
	}
}

func TestPythonProvider_GenerateEnvironment(t *testing.T) {
	provider := NewPythonProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "python",
		Version:        "3.9",
		PackageManager: &types.PackageManager{Name: "pip"},
		Framework:      "Django",
	}

	env := provider.GenerateEnvironment(result)

	// Check for Python-specific environment variables
	if env["PYTHONUNBUFFERED"] != "1" {
		t.Errorf("expected PYTHONUNBUFFERED '1', got %s", env["PYTHONUNBUFFERED"])
	}

	if env["PYTHONDONTWRITEBYTECODE"] != "1" {
		t.Errorf("expected PYTHONDONTWRITEBYTECODE '1', got %s", env["PYTHONDONTWRITEBYTECODE"])
	}

	if env["PORT"] != "8000" {
		t.Errorf("expected PORT '8000', got %s", env["PORT"])
	}

	if env["PYTHON_VERSION"] != "3.9" {
		t.Errorf("expected PYTHON_VERSION '3.9', got %s", env["PYTHON_VERSION"])
	}

	// Note: Django-specific environment variables are not set by default
	// This is expected behavior based on the implementation
}

func TestPythonProvider_NeedsNativeCompilation(t *testing.T) {
	provider := NewPythonProvider()

	testCases := []struct {
		name     string
		result   *types.DetectResult
		expected bool
	}{
		{
			name: "basic python project",
			result: &types.DetectResult{
				Matched:  true,
				Language: "python",
			},
			expected: false,
		},
		{
			name: "python with native dependencies",
			result: &types.DetectResult{
				Matched:  true,
				Language: "python",
				Evidence: types.Evidence{
					Files: []string{"requirements.txt", "numpy", "scipy"},
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

func TestPythonProvider_DetectFramework(t *testing.T) {
	provider := NewPythonProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	testCases := []struct {
		name                string
		requirementsContent string
		expectedFramework   string
	}{
		{
			name: "Django framework",
			requirementsContent: `Django==4.2.0
psycopg2-binary==2.9.7
gunicorn==21.2.0
`,
			expectedFramework: "Django",
		},
		{
			name: "Flask framework",
			requirementsContent: `Flask==2.3.2
Werkzeug==2.3.6
Jinja2==3.1.2
`,
			expectedFramework: "Flask",
		},
		{
			name: "FastAPI framework",
			requirementsContent: `fastapi==0.100.0
uvicorn==0.23.0
pydantic==2.0.0
`,
			expectedFramework: "FastAPI",
		},
		{
			name: "no framework",
			requirementsContent: `requests==2.31.0
pytest==7.4.0
numpy==1.24.0
`,
			expectedFramework: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "python-framework-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Write requirements.txt
			requirementsPath := filepath.Join(tempDir, "requirements.txt")
			if err := os.WriteFile(requirementsPath, []byte(tc.requirementsContent), 0644); err != nil {
				t.Fatalf("failed to write requirements.txt: %v", err)
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

func TestPythonProvider_DetectPythonVersion(t *testing.T) {
	provider := NewPythonProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	testCases := []struct {
		name            string
		fileContent     string
		fileName        string
		expectedVersion string
		expectedSource  string
	}{
		{
			name:            ".python-version file",
			fileContent:     "3.9.16",
			fileName:        ".python-version",
			expectedVersion: "3.9",
			expectedSource:  ".python-version",
		},
		{
			name:            "runtime.txt file",
			fileContent:     "python-3.11.4",
			fileName:        "runtime.txt",
			expectedVersion: "3.11",
			expectedSource:  "runtime.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "python-version-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Write version file
			versionPath := filepath.Join(tempDir, tc.fileName)
			if err := os.WriteFile(versionPath, []byte(tc.fileContent), 0644); err != nil {
				t.Fatalf("failed to write %s: %v", tc.fileName, err)
			}

			version, err := provider.detectPythonVersion(tempDir, gitHandler)
			if err != nil {
				t.Fatalf("detectPythonVersion failed: %v", err)
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

func TestPythonProvider_DetectPythonVersion_FromPyprojectToml(t *testing.T) {
	provider := NewPythonProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with pyproject.toml
	tempDir, err := os.MkdirTemp("", "python-pyproject-version-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create pyproject.toml with Python version
	pyprojectContent := `[tool.poetry]
name = "my-app"
version = "0.1.0"

[tool.poetry.dependencies]
python = "^3.10"
fastapi = "^0.100.0"
`
	pyprojectPath := filepath.Join(tempDir, "pyproject.toml")
	if err := os.WriteFile(pyprojectPath, []byte(pyprojectContent), 0644); err != nil {
		t.Fatalf("failed to write pyproject.toml: %v", err)
	}

	version, err := provider.detectPythonVersion(tempDir, gitHandler)
	if err != nil {
		t.Fatalf("detectPythonVersion failed: %v", err)
	}

	if version.Version != "3.10" {
		t.Errorf("expected version '3.10', got %s", version.Version)
	}

	if version.Source != "pyproject.toml" {
		t.Errorf("expected source 'pyproject.toml', got %s", version.Source)
	}
}
