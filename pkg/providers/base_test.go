package providers

import (
	"regexp"
	"testing"

	"github.com/labring/devbox-pack/pkg/types"
)

func TestBaseProvider_GetName(t *testing.T) {
	provider := &BaseProvider{
		Name:     "test",
		Language: "test-lang",
		Priority: 50,
	}

	if provider.GetName() != "test" {
		t.Errorf("expected name 'test', got %s", provider.GetName())
	}
}

func TestBaseProvider_GetLanguage(t *testing.T) {
	provider := &BaseProvider{
		Name:     "test",
		Language: "test-lang",
		Priority: 50,
	}

	if provider.GetLanguage() != "test-lang" {
		t.Errorf("expected language 'test-lang', got %s", provider.GetLanguage())
	}
}

func TestBaseProvider_GetPriority(t *testing.T) {
	provider := &BaseProvider{
		Name:     "test",
		Language: "test-lang",
		Priority: 50,
	}

	if provider.GetPriority() != 50 {
		t.Errorf("expected priority 50, got %d", provider.GetPriority())
	}
}

func TestBaseProvider_HasFile(t *testing.T) {
	provider := &BaseProvider{}

	files := []types.FileInfo{
		{Path: "package.json", IsDirectory: false},
		{Path: "src/main.js", IsDirectory: false},
		{Path: "test.txt", IsDirectory: false},
		{Path: "config.yaml", IsDirectory: false},
		{Path: "docs", IsDirectory: true},
	}

	testCases := []struct {
		name     string
		fileName string
		expected bool
	}{
		{"exact match", "package.json", true},
		{"not found", "missing.txt", false},
		{"wildcard js files", "*.js", true},
		{"wildcard txt files", "*.txt", true},
		{"wildcard no match", "*.py", false},
		{"directory", "docs", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := provider.HasFile(files, tc.fileName)
			if result != tc.expected {
				t.Errorf("HasFile(%s) = %v, expected %v", tc.fileName, result, tc.expected)
			}
		})
	}
}

func TestBaseProvider_HasAnyFile(t *testing.T) {
	provider := &BaseProvider{}

	files := []types.FileInfo{
		{Path: "package.json", IsDirectory: false},
		{Path: "main.py", IsDirectory: false},
	}

	testCases := []struct {
		name      string
		fileNames []string
		expected  bool
	}{
		{"has one", []string{"package.json", "missing.txt"}, true},
		{"has none", []string{"missing1.txt", "missing2.txt"}, false},
		{"has all", []string{"package.json", "main.py"}, true},
		{"empty list", []string{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := provider.HasAnyFile(files, tc.fileNames)
			if result != tc.expected {
				t.Errorf("HasAnyFile(%v) = %v, expected %v", tc.fileNames, result, tc.expected)
			}
		})
	}
}

func TestBaseProvider_HasAllFiles(t *testing.T) {
	provider := &BaseProvider{}

	files := []types.FileInfo{
		{Path: "package.json", IsDirectory: false},
		{Path: "main.py", IsDirectory: false},
	}

	testCases := []struct {
		name      string
		fileNames []string
		expected  bool
	}{
		{"has all", []string{"package.json", "main.py"}, true},
		{"missing one", []string{"package.json", "missing.txt"}, false},
		{"has none", []string{"missing1.txt", "missing2.txt"}, false},
		{"empty list", []string{}, true},
		{"single match", []string{"package.json"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := provider.HasAllFiles(files, tc.fileNames)
			if result != tc.expected {
				t.Errorf("HasAllFiles(%v) = %v, expected %v", tc.fileNames, result, tc.expected)
			}
		})
	}
}

func TestBaseProvider_GetMatchingFiles(t *testing.T) {
	provider := &BaseProvider{}

	files := []types.FileInfo{
		{Path: "main.js", IsDirectory: false},
		{Path: "test.js", IsDirectory: false},
		{Path: "config.json", IsDirectory: false},
		{Path: "README.md", IsDirectory: false},
	}

	testCases := []struct {
		name     string
		pattern  string
		expected int
	}{
		{"js files", "*.js", 2},
		{"json files", "*.json", 1},
		{"md files", "*.md", 1},
		{"no match", "*.py", 0},
		{"exact match", "main.js", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := provider.GetMatchingFiles(files, tc.pattern)
			if len(result) != tc.expected {
				t.Errorf("GetMatchingFiles(%s) returned %d files, expected %d", tc.pattern, len(result), tc.expected)
			}
		})
	}
}

func TestBaseProvider_ParseVersionFromJSON(t *testing.T) {
	provider := &BaseProvider{}
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// Create test JSON file
	jsonContent := `{
		"version": "1.2.3",
		"engines": {
			"node": ">=18.0.0"
		}
	}`

	helper.WriteFile("package.json", jsonContent)

	testCases := []struct {
		name         string
		versionField string
		expected     string
		expectError  bool
	}{
		{"version field", "version", "1.2.3", false},
		{"nested field", "engines.node", ">=18.0.0", false},
		{"missing field", "missing", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := provider.ParseVersionFromJSON(helper.TempDir, "package.json", helper.GitHandler, tc.versionField)

			if tc.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tc.expected {
				t.Errorf("expected version %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestBaseProvider_ParseVersionFromText(t *testing.T) {
	provider := &BaseProvider{}
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// Create test text file
	textContent := "v18.17.0\n"
	helper.WriteFile(".nvmrc", textContent)

	pattern := regexp.MustCompile(`^v?(.+)$`)

	result, err := provider.ParseVersionFromText(helper.TempDir, ".nvmrc", helper.GitHandler, pattern)
	if err != nil {
		t.Fatalf("ParseVersionFromText failed: %v", err)
	}

	if result != "18.17.0" {
		t.Errorf("expected version '18.17.0', got '%s'", result)
	}
}

func TestBaseProvider_DetectPackageManager(t *testing.T) {
	provider := &BaseProvider{}

	lockFiles := map[string]string{
		"package-lock.json": "npm",
		"yarn.lock":         "yarn",
		"pnpm-lock.yaml":    "pnpm",
	}

	testCases := []struct {
		name     string
		files    []types.FileInfo
		expected string
	}{
		{
			name: "npm lock file",
			files: []types.FileInfo{
				{Path: "package.json", IsDirectory: false},
				{Path: "package-lock.json", IsDirectory: false},
			},
			expected: "npm",
		},
		{
			name: "yarn lock file",
			files: []types.FileInfo{
				{Path: "package.json", IsDirectory: false},
				{Path: "yarn.lock", IsDirectory: false},
			},
			expected: "yarn",
		},
		{
			name: "pnpm lock file",
			files: []types.FileInfo{
				{Path: "package.json", IsDirectory: false},
				{Path: "pnpm-lock.yaml", IsDirectory: false},
			},
			expected: "pnpm",
		},
		{
			name: "no lock file",
			files: []types.FileInfo{
				{Path: "package.json", IsDirectory: false},
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := provider.DetectPackageManager(tc.files, lockFiles)
			if result != tc.expected {
				t.Errorf("expected package manager %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestBaseProvider_CreateDetectResult(t *testing.T) {
	provider := &BaseProvider{}

	version := &types.VersionInfo{Version: "1.0.0", Source: "test"}
	metadata := map[string]interface{}{"test": "value"}
	evidence := types.Evidence{Files: []string{"test.txt"}, Reason: "test reason"}

	result := provider.CreateDetectResult(
		true,
		0.8,
		"test-lang",
		version,
		"test-framework",
		"test-pm",
		"test-build",
		metadata,
		evidence,
	)

	if !result.Matched {
		t.Error("expected Matched to be true")
	}

	if result.Confidence != 0.8 {
		t.Errorf("expected confidence 0.8, got %f", result.Confidence)
	}

	if result.Language != "test-lang" {
		t.Errorf("expected language 'test-lang', got %s", result.Language)
	}

	if result.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %s", result.Version)
	}

	if result.Framework != "test-framework" {
		t.Errorf("expected framework 'test-framework', got %s", result.Framework)
	}

	if result.PackageManager == nil || result.PackageManager.Name != "test-pm" {
		t.Errorf("expected package manager 'test-pm', got %v", result.PackageManager)
	}

	if len(result.BuildTools) == 0 || result.BuildTools[0] != "test-build" {
		t.Errorf("expected build tool 'test-build', got %v", result.BuildTools)
	}
}

func TestBaseProvider_CreateVersionInfo(t *testing.T) {
	provider := &BaseProvider{}

	version := provider.CreateVersionInfo("1.2.3", "package.json")

	if version.Version != "1.2.3" {
		t.Errorf("expected version '1.2.3', got %s", version.Version)
	}

	if version.Source != "package.json" {
		t.Errorf("expected source 'package.json', got %s", version.Source)
	}
}