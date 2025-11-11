package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

// TestHelper provides common utilities for provider tests
type TestHelper struct {
	T         *testing.T
	GitHandler *git.GitHandler
	TempDir   string
}

// NewTestHelper creates a new test helper with temporary directory and git handler
func NewTestHelper(t *testing.T) *TestHelper {
	tempDir, err := os.MkdirTemp("", "provider-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	return &TestHelper{
		T:         t,
		GitHandler: git.NewGitHandler(),
		TempDir:   tempDir,
	}
}

// Cleanup cleans up temporary resources
func (h *TestHelper) Cleanup() {
	if h.GitHandler != nil {
		h.GitHandler.Cleanup()
	}
	if h.TempDir != "" {
		os.RemoveAll(h.TempDir)
	}
}

// WriteFile writes a test file to the temporary directory
func (h *TestHelper) WriteFile(filename, content string) string {
	filePath := filepath.Join(h.TempDir, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		h.T.Fatalf("failed to write %s: %v", filename, err)
	}
	return filePath
}

// CreateTempDir creates a temporary subdirectory
func (h *TestHelper) CreateTempDir(name string) string {
	dirPath := filepath.Join(h.TempDir, name)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		h.T.Fatalf("failed to create temp dir %s: %v", name, err)
	}
	return dirPath
}

// AssertProviderBasic asserts basic provider properties
func AssertProviderBasic(t *testing.T, provider interface{ GetName() string; GetPriority() int }, expectedName string, expectedPriority int) {
	if provider == nil {
		t.Fatal("provider is nil")
	}

	if provider.GetName() != expectedName {
		t.Errorf("expected name '%s', got '%s'", expectedName, provider.GetName())
	}

	if provider.GetPriority() != expectedPriority {
		t.Errorf("expected priority %d, got %d", expectedPriority, provider.GetPriority())
	}
}

// AssertDetectResult asserts detection result properties
func AssertDetectResult(t *testing.T, result *types.DetectResult, expectedMatched bool, expectedLanguage string) {
	if result == nil {
		t.Fatal("detect result is nil")
		return
	}

	if result.Matched != expectedMatched {
		t.Errorf("expected matched %v, got %v", expectedMatched, result.Matched)
	}

	if expectedMatched && result.Language != expectedLanguage {
		t.Errorf("expected language '%s', got '%s'", expectedLanguage, result.Language)
	}
}

// AssertPackageManager asserts package manager properties
func AssertPackageManager(t *testing.T, pm *types.PackageManager, expectedName string) {
	if pm == nil {
		if expectedName != "" {
			t.Errorf("expected package manager '%s', got nil", expectedName)
		}
		return
	}

	if pm.Name != expectedName {
		t.Errorf("expected package manager '%s', got '%s'", expectedName, pm.Name)
	}
}

// AssertCommandContains asserts that a command list contains expected text
func AssertCommandContains(t *testing.T, commands []string, expectedText string) bool {
	for _, cmd := range commands {
		if strings.Contains(cmd, expectedText) {
			return true
		}
	}
	t.Errorf("expected command containing '%s' in %v", expectedText, commands)
	return false
}

// AssertCommandNotContains asserts that a command list does not contain text
func AssertCommandNotContains(t *testing.T, commands []string, forbiddenText string) {
	for _, cmd := range commands {
		if strings.Contains(cmd, forbiddenText) {
			t.Errorf("found unexpected command containing '%s': %s", forbiddenText, cmd)
			return
		}
	}
}

// AssertEnvironmentVar asserts environment variable value
func AssertEnvironmentVar(t *testing.T, env map[string]string, key, expectedValue string) {
	actualValue, exists := env[key]
	if !exists {
		t.Errorf("expected environment variable '%s' to be set", key)
		return
	}

	if actualValue != expectedValue {
		t.Errorf("expected %s='%s', got '%s'", key, expectedValue, actualValue)
	}
}

// AssertEnvironmentVarExists asserts that an environment variable exists (regardless of value)
func AssertEnvironmentVarExists(t *testing.T, env map[string]string, key string) {
	if _, exists := env[key]; !exists {
		t.Errorf("expected environment variable '%s' to be set", key)
	}
}

// ProviderTestCase represents a test case for provider detection
type ProviderTestCase struct {
	Name            string
	Files           []types.FileInfo
	ExpectedMatch   bool
	ExpectedLang    string
	ExpectedPM      string
	ExpectedVersion string
	MinConfidence   float64
}

// RunProviderTestCases runs multiple provider test cases
func RunProviderTestCases(t *testing.T, provider interface{ Detect(string, []types.FileInfo, *git.GitHandler) (*types.DetectResult, error) }, testCases []ProviderTestCase) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			helper := NewTestHelper(t)
			defer helper.Cleanup()

			result, err := provider.Detect(helper.TempDir, tc.Files, helper.GitHandler)
			if err != nil {
				t.Fatalf("Detect failed: %v", err)
			}

			AssertDetectResult(t, result, tc.ExpectedMatch, tc.ExpectedLang)

			if tc.ExpectedMatch {
				if tc.ExpectedPM != "" {
					AssertPackageManager(t, result.PackageManager, tc.ExpectedPM)
				}

				if tc.ExpectedVersion != "" && result.Version != tc.ExpectedVersion {
					t.Errorf("expected version '%s', got '%s'", tc.ExpectedVersion, result.Version)
				}

				if tc.MinConfidence > 0 && result.Confidence < tc.MinConfidence {
					t.Errorf("expected confidence >= %f, got %f", tc.MinConfidence, result.Confidence)
				}
			}
		})
	}
}

// CreateTestFiles creates test file structure from file map
func CreateTestFiles(helper *TestHelper, files map[string]string) []types.FileInfo {
	var fileInfos []types.FileInfo

	for path, content := range files {
		if strings.HasSuffix(path, "/") {
			// Create directory
			helper.CreateTempDir(strings.TrimSuffix(path, "/"))
			fileInfos = append(fileInfos, types.FileInfo{Path: path, IsDirectory: true})
		} else {
			// Create file
			helper.WriteFile(path, content)
			fileInfos = append(fileInfos, types.FileInfo{Path: path, IsDirectory: false})
		}
	}

	return fileInfos
}

// AssertNeedsNativeCompilation asserts native compilation requirements
func AssertNeedsNativeCompilation(t *testing.T, provider interface{ NeedsNativeCompilation(*types.DetectResult) bool }, result *types.DetectResult, expected bool) {
	actual := provider.NeedsNativeCompilation(result)
	if actual != expected {
		t.Errorf("expected NeedsNativeCompilation %v, got %v", expected, actual)
	}
}