package providers

import (
	"testing"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

// SimpleTestSuite provides basic standardized tests for providers
type SimpleTestSuite struct{}

// NewSimpleTestSuite creates a new simple test suite
func NewSimpleTestSuite() *SimpleTestSuite {
	return &SimpleTestSuite{}
}

// RunBasicProviderTest runs basic provider tests using reflection
func (suite *SimpleTestSuite) RunBasicProviderTest(t *testing.T, provider any) {
	// Get provider methods via reflection - simplified approach
	t.Run("BasicProperties", func(t *testing.T) {
		// This is a simplified test that just checks the provider is not nil
		if provider == nil {
			t.Fatal("provider is nil")
		}
		t.Log("Provider basic properties test passed")
	})
}

// RunDetectionTest runs detection test with given files
func (suite *SimpleTestSuite) RunDetectionTest(t *testing.T, provider any, files []types.FileInfo, expectedMatch bool, expectedLang string) {
	t.Run("Detection", func(t *testing.T) {
		helper := NewTestHelper(t)
		defer helper.Cleanup()

		// For now, we'll use a type assertion approach
		switch p := provider.(type) {
		case interface{ Detect(string, []types.FileInfo, *git.GitHandler) (*types.DetectResult, error) }:
			result, err := p.Detect(helper.TempDir, files, helper.GitHandler)
			if err != nil {
				t.Fatalf("Detect failed: %v", err)
			}
			AssertDetectResult(t, result, expectedMatch, expectedLang)
		default:
			t.Skip("Provider does not implement Detect method")
		}
	})
}

// RunCommandTest runs command generation test
func (suite *SimpleTestSuite) RunCommandTest(t *testing.T, provider any, result *types.DetectResult) {
	t.Run("Commands", func(t *testing.T) {
		switch p := provider.(type) {
		case interface{ GenerateCommands(*types.DetectResult, types.CLIOptions) types.Commands }:
			options := types.CLIOptions{}
			commands := p.GenerateCommands(result, options)

			if len(commands.Run) == 0 && len(commands.Setup) == 0 {
				t.Error("expected at least some commands to be generated")
			}
		default:
			t.Skip("Provider does not implement GenerateCommands method")
		}
	})
}

// RunEnvironmentTest runs environment generation test
func (suite *SimpleTestSuite) RunEnvironmentTest(t *testing.T, provider any, result *types.DetectResult) {
	t.Run("Environment", func(t *testing.T) {
		switch p := provider.(type) {
		case interface{ GenerateEnvironment(*types.DetectResult) map[string]string }:
			env := p.GenerateEnvironment(result)
			if env == nil {
				t.Fatal("GenerateEnvironment returned nil")
			}
		default:
			t.Skip("Provider does not implement GenerateEnvironment method")
		}
	})
}