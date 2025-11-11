package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewPHPProvider(t *testing.T) {
	provider := NewPHPProvider()

	if provider == nil {
		t.Fatal("NewPHPProvider() returned nil")
	}

	if provider.GetName() != "php" {
		t.Errorf("expected name 'php', got %s", provider.GetName())
	}

	if provider.GetPriority() != 10 {
		t.Errorf("expected priority 10, got %d", provider.GetPriority())
	}
}

func TestPHPProvider_Detect_NoPHPFiles(t *testing.T) {
	provider := NewPHPProvider()
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
		t.Error("expected not matched for non-PHP project")
	}
}

func TestPHPProvider_Detect_WithComposer(t *testing.T) {
	provider := NewPHPProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with composer.json
	tempDir, err := os.MkdirTemp("", "php-composer-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	composerJsonContent := `{
    "name": "example/my-app",
    "description": "A sample PHP application",
    "type": "project",
    "require": {
        "php": "^8.1",
        "laravel/framework": "^10.0"
    },
    "require-dev": {
        "phpunit/phpunit": "^10.0"
    },
    "autoload": {
        "psr-4": {
            "App\\": "app/"
        }
    }
}`

	composerJsonPath := filepath.Join(tempDir, "composer.json")
	if err := os.WriteFile(composerJsonPath, []byte(composerJsonContent), 0644); err != nil {
		t.Fatalf("failed to write composer.json: %v", err)
	}

	files := []types.FileInfo{
		{Path: "composer.json", IsDirectory: false},
		{Path: "index.php", IsDirectory: false},
		{Path: "app/", IsDirectory: true},
	}

	result, err := provider.Detect(tempDir, files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for PHP Composer project")
	}

	if result.Language != "php" {
		t.Errorf("expected language 'php', got %s", result.Language)
	}

	if result.PackageManager == nil || result.PackageManager.Name != "composer" {
		t.Errorf("expected package manager 'composer', got %s", func() string {
			if result.PackageManager == nil {
				return "nil"
			}
			return result.PackageManager.Name
		}())
	}
}

func TestPHPProvider_Detect_WithPHPFiles(t *testing.T) {
	provider := NewPHPProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "index.php", IsDirectory: false},
		{Path: "config.php", IsDirectory: false},
		{Path: "vendor/", IsDirectory: true},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for PHP project with .php files")
	}
}

func TestPHPProvider_Detect_WithLaravel(t *testing.T) {
	provider := NewPHPProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "composer.json", IsDirectory: false},
		{Path: "artisan", IsDirectory: false},
		{Path: "app/", IsDirectory: true},
		{Path: "config/", IsDirectory: true},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Laravel project")
	}
}

func TestPHPProvider_Detect_WithWordPress(t *testing.T) {
	provider := NewPHPProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "wp-config.php", IsDirectory: false},
		{Path: "wp-content/", IsDirectory: true},
		{Path: "wp-includes/", IsDirectory: true},
		{Path: "index.php", IsDirectory: false},  // Add index.php to increase confidence
		{Path: "wp-load.php", IsDirectory: false}, // Another WordPress file
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for WordPress project")
	}

	// Note: The current implementation doesn't detect WordPress as a framework
	// This test verifies that WordPress files are detected as a PHP project
	if result.Framework != "" {
		t.Logf("detected framework: %s (WordPress detection not implemented)", result.Framework)
	}
}

func TestPHPProvider_GenerateCommands_BasicProject(t *testing.T) {
	provider := NewPHPProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "php",
		PackageManager: &types.PackageManager{Name: "Composer"},
		Framework:      "",
		Evidence: types.Evidence{
			Files: []string{"composer.json", "index.php"},
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

	// Check for Composer-specific commands
	foundComposerInstall := false
	for _, cmd := range commands.Setup {
		if strings.Contains(cmd, "composer install") {
			foundComposerInstall = true
			break
		}
	}
	if !foundComposerInstall {
		t.Error("expected 'composer install' command in setup commands")
	}

	// Check for PHP server command
	foundPHPServer := false
	for _, cmd := range commands.Run {
		if strings.Contains(cmd, "php -S") {
			foundPHPServer = true
			break
		}
	}
	if !foundPHPServer {
		t.Error("expected 'php -S' command in run commands")
	}
}

func TestPHPProvider_GenerateCommands_LaravelProject(t *testing.T) {
	provider := NewPHPProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "php",
		PackageManager: &types.PackageManager{Name: "Composer"},
		Framework:      "Laravel",
		Evidence: types.Evidence{
			Files: []string{"composer.json", "artisan"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Dev) == 0 {
		t.Error("expected dev commands")
	}

	// Check for Laravel-specific commands
	foundArtisanServe := false
	for _, cmd := range commands.Dev {
		if strings.Contains(cmd, "php artisan serve") {
			foundArtisanServe = true
			break
		}
	}
	if !foundArtisanServe {
		t.Error("expected 'php artisan serve' command in dev commands")
	}
}

func TestPHPProvider_GenerateEnvironment(t *testing.T) {
	provider := NewPHPProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "php",
		Version:        "8.1",
		PackageManager: &types.PackageManager{Name: "Composer"},
		Framework:      "Laravel",
	}

	env := provider.GenerateEnvironment(result)

	if env["PHP_ENV"] != "production" {
		t.Errorf("expected PHP_ENV 'production', got %s", env["PHP_ENV"])
	}

	if env["PORT"] != "8000" {
		t.Errorf("expected PORT '8000', got %s", env["PORT"])
	}

	// Laravel sets APP_ENV to 'local' by default
	if env["APP_ENV"] != "local" {
		t.Errorf("expected APP_ENV 'local', got %s", env["APP_ENV"])
	}
}

func TestPHPProvider_NeedsNativeCompilation(t *testing.T) {
	provider := NewPHPProvider()

	testCases := []struct {
		name     string
		result   *types.DetectResult
		expected bool
	}{
		{
			name: "basic php project",
			result: &types.DetectResult{
				Matched:  true,
				Language: "php",
			},
			expected: false,
		},
		{
			name: "php with extensions",
			result: &types.DetectResult{
				Matched:  true,
				Language: "php",
				Metadata: map[string]interface{}{
					"hasExtensions": true,
				},
			},
			expected: false, // PHP projects don't need native compilation by default
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

func TestPHPProvider_DetectFramework(t *testing.T) {
	provider := NewPHPProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	testCases := []struct {
		name              string
		composerContent   string
		expectedFramework string
	}{
		{
			name: "Laravel framework",
			composerContent: `{
    "name": "example/my-app",
    "require": {
        "php": "^8.1",
        "laravel/framework": "^10.0"
    }
}`,
			expectedFramework: "Laravel",
		},
		{
			name: "Symfony framework",
			composerContent: `{
    "name": "example/my-app",
    "require": {
        "php": "^8.1",
        "symfony/framework-bundle": "^6.0"
    }
}`,
			expectedFramework: "Symfony",
		},
		{
			name: "CodeIgniter framework",
			composerContent: `{
    "name": "example/my-app",
    "require": {
        "php": "^8.1",
        "codeigniter4/framework": "^4.0"
    }
}`,
			expectedFramework: "CodeIgniter",
		},
		{
			name: "no framework",
			composerContent: `{
    "name": "example/my-app",
    "require": {
        "php": "^8.1",
        "monolog/monolog": "^3.0"
    }
}`,
			expectedFramework: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "php-framework-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Write composer.json
			composerPath := filepath.Join(tempDir, "composer.json")
			if err := os.WriteFile(composerPath, []byte(tc.composerContent), 0644); err != nil {
				t.Fatalf("failed to write composer.json: %v", err)
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

func TestPHPProvider_DetectPHPVersion(t *testing.T) {
	provider := NewPHPProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	testCases := []struct {
		name            string
		composerContent string
		expectedVersion string
		expectedSource  string
	}{
		{
			name: "composer.json with PHP version",
			composerContent: `{
    "name": "example/my-app",
    "require": {
        "php": "^8.1"
    }
}`,
			expectedVersion: "8.1",
			expectedSource:  "composer.json require",
		},
		{
			name: "composer.json with specific PHP version",
			composerContent: `{
    "name": "example/my-app",
    "require": {
        "php": ">=8.0.0"
    }
}`,
			expectedVersion: "8.0",
			expectedSource:  "composer.json require",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "php-version-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Write composer.json
			composerPath := filepath.Join(tempDir, "composer.json")
			if err := os.WriteFile(composerPath, []byte(tc.composerContent), 0644); err != nil {
				t.Fatalf("failed to write composer.json: %v", err)
			}

			version, err := provider.detectPHPVersion(tempDir, gitHandler)
			if err != nil {
				t.Fatalf("detectPHPVersion failed: %v", err)
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

func TestPHPProvider_DetectPHPVersion_FromPHPVersion(t *testing.T) {
	provider := NewPHPProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with .php-version
	tempDir, err := os.MkdirTemp("", "php-version-file-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .php-version file
	phpVersionContent := "8.2.0"
	phpVersionPath := filepath.Join(tempDir, ".php-version")
	if err := os.WriteFile(phpVersionPath, []byte(phpVersionContent), 0644); err != nil {
		t.Fatalf("failed to write .php-version: %v", err)
	}

	version, err := provider.detectPHPVersion(tempDir, gitHandler)
	if err != nil {
		t.Fatalf("detectPHPVersion failed: %v", err)
	}

	if version.Version != "8.2.0" {
		t.Errorf("expected version '8.2.0', got %s", version.Version)
	}

	if version.Source != ".php-version" {
		t.Errorf("expected source '.php-version', got %s", version.Source)
	}
}
