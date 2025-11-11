package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewRubyProvider(t *testing.T) {
	provider := NewRubyProvider()

	if provider == nil {
		t.Fatal("NewRubyProvider() returned nil")
	}

	if provider.GetName() != "ruby" {
		t.Errorf("expected name 'ruby', got %s", provider.GetName())
	}

	if provider.GetPriority() != 50 {
		t.Errorf("expected priority 50, got %d", provider.GetPriority())
	}
}

func TestRubyProvider_Detect_NoRubyFiles(t *testing.T) {
	provider := NewRubyProvider()
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
		t.Error("expected not matched for non-Ruby project")
	}
}

func TestRubyProvider_Detect_WithGemfile(t *testing.T) {
	provider := NewRubyProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with Gemfile
	tempDir, err := os.MkdirTemp("", "ruby-gemfile-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gemfileContent := `source 'https://rubygems.org'

ruby '3.1.0'

gem 'rails', '~> 7.0.0'
gem 'sqlite3', '~> 1.4'
gem 'puma', '~> 5.0'
gem 'bootsnap', '>= 1.4.4', require: false

group :development, :test do
  gem 'byebug', platforms: [:mri, :mingw, :x64_mingw]
  gem 'rspec-rails'
end
`

	gemfilePath := filepath.Join(tempDir, "Gemfile")
	if err := os.WriteFile(gemfilePath, []byte(gemfileContent), 0644); err != nil {
		t.Fatalf("failed to write Gemfile: %v", err)
	}

	files := []types.FileInfo{
		{Path: "Gemfile", IsDirectory: false},
		{Path: "app.rb", IsDirectory: false},
		{Path: "config/", IsDirectory: true},
	}

	result, err := provider.Detect(tempDir, files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Ruby project with Gemfile")
	}

	if result.Language != "ruby" {
		t.Errorf("expected language 'ruby', got %s", result.Language)
	}

	if result.PackageManager == nil || result.PackageManager.Name != "bundler" {
		t.Errorf("expected package manager 'bundler', got %s", func() string {
			if result.PackageManager == nil {
				return "nil"
			}
			return result.PackageManager.Name
		}())
	}
}

func TestRubyProvider_Detect_WithRubyFiles(t *testing.T) {
	provider := NewRubyProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "main.rb", IsDirectory: false},
		{Path: "lib/", IsDirectory: true},
		{Path: "lib/utils.rb", IsDirectory: false},
		{Path: "spec/", IsDirectory: true},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Ruby project with .rb files")
	}
}

func TestRubyProvider_Detect_WithRails(t *testing.T) {
	provider := NewRubyProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "Gemfile", IsDirectory: false},
		{Path: "config.ru", IsDirectory: false},
		{Path: "app/", IsDirectory: true},
		{Path: "config/", IsDirectory: true},
		{Path: "config/application.rb", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Rails project")
	}
}

func TestRubyProvider_Detect_WithSinatra(t *testing.T) {
	provider := NewRubyProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	files := []types.FileInfo{
		{Path: "Gemfile", IsDirectory: false},
		{Path: "config.ru", IsDirectory: false},
		{Path: "app.rb", IsDirectory: false},
	}

	result, err := provider.Detect("", files, gitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !result.Matched {
		t.Error("expected matched for Sinatra project")
	}
}

func TestRubyProvider_GenerateCommands_BasicProject(t *testing.T) {
	provider := NewRubyProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "ruby",
		PackageManager: &types.PackageManager{Name: "bundler"},
		Framework:      "",
		Evidence: types.Evidence{
			Files: []string{"Gemfile", "main.rb"},
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

	// Check for bundle install command
	foundBundleInstall := false
	for _, cmd := range commands.Setup {
		if strings.Contains(cmd, "bundle install") {
			foundBundleInstall = true
			break
		}
	}
	if !foundBundleInstall {
		t.Error("expected 'bundle install' command in setup commands")
	}

	// Check for ruby run command
	foundRubyRun := false
	for _, cmd := range commands.Run {
		if strings.Contains(cmd, "ruby") {
			foundRubyRun = true
			break
		}
	}
	if !foundRubyRun {
		t.Error("expected 'ruby' command in run commands")
	}
}

func TestRubyProvider_GenerateCommands_RailsProject(t *testing.T) {
	provider := NewRubyProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "ruby",
		PackageManager: &types.PackageManager{Name: "bundler"},
		Framework:      "Rails",
		Evidence: types.Evidence{
			Files: []string{"Gemfile", "config.ru", "config/application.rb"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Dev) == 0 {
		t.Error("expected dev commands")
	}

	// Check for Rails-specific commands
	foundRailsServer := false
	for _, cmd := range commands.Dev {
		if strings.Contains(cmd, "rails server") || strings.Contains(cmd, "bundle exec rails s") {
			foundRailsServer = true
			break
		}
	}
	if !foundRailsServer {
		t.Error("expected Rails server command in dev commands")
	}
}

func TestRubyProvider_GenerateCommands_SinatraProject(t *testing.T) {
	provider := NewRubyProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "ruby",
		PackageManager: &types.PackageManager{Name: "bundler"},
		Framework:      "Sinatra",
		Evidence: types.Evidence{
			Files: []string{"Gemfile", "config.ru", "app.rb"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Dev) == 0 {
		t.Error("expected dev commands")
	}

	// Check for Sinatra-specific commands
	foundSinatraRun := false
	for _, cmd := range commands.Dev {
		if strings.Contains(cmd, "rackup") || strings.Contains(cmd, "bundle exec rackup") {
			foundSinatraRun = true
			break
		}
	}
	if !foundSinatraRun {
		t.Error("expected Sinatra/Rack command in dev commands")
	}
}

func TestRubyProvider_GenerateEnvironment(t *testing.T) {
	provider := NewRubyProvider()

	result := &types.DetectResult{
		Matched:        true,
		Language:       "ruby",
		Version:        "3.1.0",
		PackageManager: &types.PackageManager{Name: "bundler"},
		Framework:      "Rails",
	}

	env := provider.GenerateEnvironment(result)

	if env["RACK_ENV"] != "production" {
		t.Errorf("expected RACK_ENV 'production', got %s", env["RACK_ENV"])
	}

	if env["RAILS_ENV"] != "production" {
		t.Errorf("expected RAILS_ENV 'production', got %s", env["RAILS_ENV"])
	}

	if env["PORT"] != "3000" {
		t.Errorf("expected PORT '3000', got %s", env["PORT"])
	}

	if env["BUNDLE_WITHOUT"] != "development:test" {
		t.Errorf("expected BUNDLE_WITHOUT 'development:test', got %s", env["BUNDLE_WITHOUT"])
	}
}

func TestRubyProvider_NeedsNativeCompilation(t *testing.T) {
	provider := NewRubyProvider()

	testCases := []struct {
		name     string
		result   *types.DetectResult
		expected bool
	}{
		{
			name: "basic ruby project",
			result: &types.DetectResult{
				Matched:  true,
				Language: "ruby",
			},
			expected: false,
		},
		{
			name: "ruby with native gems",
			result: &types.DetectResult{
				Matched:  true,
				Language: "ruby",
				Metadata: map[string]interface{}{
					"hasNativeGems": true,
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

func TestRubyProvider_DetectFramework(t *testing.T) {
	provider := NewRubyProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	testCases := []struct {
		name              string
		gemfileContent    string
		expectedFramework string
	}{
		{
			name: "Rails framework",
			gemfileContent: `source 'https://rubygems.org'

ruby '3.1.0'

gem 'rails', '~> 7.0.0'
gem 'sqlite3', '~> 1.4'
`,
			expectedFramework: "Rails",
		},
		{
			name: "Sinatra framework",
			gemfileContent: `source 'https://rubygems.org'

ruby '3.1.0'

gem 'sinatra', '~> 3.0'
gem 'thin'
`,
			expectedFramework: "Sinatra",
		},
		{
			name: "Hanami framework",
			gemfileContent: `source 'https://rubygems.org'

ruby '3.1.0'

gem 'hanami', '~> 2.0'
gem 'pg'
`,
			expectedFramework: "Hanami",
		},
		{
			name: "no framework",
			gemfileContent: `source 'https://rubygems.org'

ruby '3.1.0'

gem 'nokogiri'
gem 'httparty'
`,
			expectedFramework: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "ruby-framework-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Write Gemfile
			gemfilePath := filepath.Join(tempDir, "Gemfile")
			if err := os.WriteFile(gemfilePath, []byte(tc.gemfileContent), 0644); err != nil {
				t.Fatalf("failed to write Gemfile: %v", err)
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

func TestRubyProvider_DetectRubyVersion(t *testing.T) {
	provider := NewRubyProvider()
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
			name:            ".ruby-version file",
			fileContent:     "3.1.0",
			fileName:        ".ruby-version",
			expectedVersion: "3.1.0",
			expectedSource:  ".ruby-version",
		},
		{
			name:            ".rvmrc file",
			fileContent:     "#!/usr/bin/env bash\nrvm use 3.0.0@myproject --create",
			fileName:        ".rvmrc",
			expectedVersion: "3.0.0",
			expectedSource:  ".rvmrc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "ruby-version-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Write version file
			versionPath := filepath.Join(tempDir, tc.fileName)
			if err := os.WriteFile(versionPath, []byte(tc.fileContent), 0644); err != nil {
				t.Fatalf("failed to write %s: %v", tc.fileName, err)
			}

			version, err := provider.detectRubyVersion(tempDir, gitHandler)
			if err != nil {
				t.Fatalf("detectRubyVersion failed: %v", err)
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

func TestRubyProvider_DetectRubyVersion_FromGemfile(t *testing.T) {
	provider := NewRubyProvider()
	gitHandler := git.NewGitHandler()
	defer gitHandler.Cleanup()

	// Create temporary directory with Gemfile
	tempDir, err := os.MkdirTemp("", "ruby-gemfile-version-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create Gemfile with Ruby version
	gemfileContent := `source 'https://rubygems.org'

ruby '3.2.0'

gem 'rails', '~> 7.0.0'
`
	gemfilePath := filepath.Join(tempDir, "Gemfile")
	if err := os.WriteFile(gemfilePath, []byte(gemfileContent), 0644); err != nil {
		t.Fatalf("failed to write Gemfile: %v", err)
	}

	version, err := provider.detectRubyVersion(tempDir, gitHandler)
	if err != nil {
		t.Fatalf("detectRubyVersion failed: %v", err)
	}

	if version.Version != "3.2.0" {
		t.Errorf("expected version '3.2.0', got %s", version.Version)
	}

	if version.Source != "Gemfile" {
		t.Errorf("expected source 'Gemfile', got %s", version.Source)
	}
}
