# Testing Guide

## Overview

The DevBox Pack project uses a comprehensive testing strategy that includes unit tests, integration tests, and validation tests to ensure reliable detection and execution plan generation across all supported languages and frameworks.

## Test Structure

### Test Organization

```
pkg/
├── detector/
│   └── engine_test.go          # Detection engine tests
├── generators/
│   └── planner_test.go         # Plan generation tests
├── service/
│   └── devbox_test.go          # Service layer tests
└── providers/
    └── [provider]_test.go      # Provider-specific tests

tests/
└── validate_plans.go           # Plan validation utilities

railpack/
├── integration_tests/
│   └── run_test.go             # Integration tests
└── core/
    └── providers/
        └── [language]/
            └── [language]_test.go  # Language provider tests
```

## Test Categories

### 1. Unit Tests

Unit tests focus on individual components and their specific functionality.

#### Provider Detection Tests

Each provider has comprehensive detection tests that verify:

- **File Pattern Recognition**: Correct identification of language-specific files
- **Confidence Scoring**: Appropriate confidence levels for different scenarios
- **Version Detection**: Accurate parsing of version information from various sources
- **Framework Detection**: Recognition of specific frameworks within languages

**Example Test Cases:**

```go
// Node.js Provider Tests
func TestDetectProject_NodeExpress(t *testing.T) {
    // Tests detection of Express.js applications
    // Verifies package.json parsing and framework identification
}

func TestNodeCorepack(t *testing.T) {
    // Tests Corepack detection and configuration
    // Verifies package manager selection logic
}
```

#### Plan Generation Tests

Plan generation tests ensure correct execution plan creation:

- **Command Generation**: Proper dev, build, and start commands
- **Base Image Selection**: Correct base image based on language and version
- **Environment Configuration**: Appropriate environment variables and tools
- **APT Dependencies**: Required system packages for each provider

**Example Test Cases:**

```go
func TestGeneratePlan_NodeProject(t *testing.T) {
    // Tests complete plan generation for Node.js projects
    // Verifies all plan components are correctly generated
}

func TestGeneratePlan_InvalidInput(t *testing.T) {
    // Tests error handling for invalid inputs
    // Ensures graceful failure scenarios
}
```

### 2. Integration Tests

Integration tests validate the complete workflow using real project examples.

#### Sample Repository Testing

The integration test suite uses a collection of sample repositories to test real-world scenarios:

**Node.js Examples:**
- Next.js applications
- Vite projects
- Nuxt applications
- Create React App projects
- Express.js APIs

**Python Examples:**
- Django applications
- FastAPI projects
- Flask applications
- Poetry-managed projects

**Java Examples:**
- Spring Boot applications
- Maven quickstart projects
- Gradle applications

**Go Examples:**
- Gin web applications
- Standard library HTTP servers
- CLI applications

**Other Languages:**
- PHP Laravel applications
- Ruby on Rails projects
- Rust Axum/Rocket applications
- Deno Fresh applications
- Static HTML/CSS/JS sites

#### Integration Test Execution

```go
func TestExamplesIntegration(t *testing.T) {
    // Runs detection and plan generation on all sample projects
    // Validates output against expected results
    // Ensures consistency across different project types
}
```

### 3. Validation Tests

Validation tests ensure generated plans meet quality standards and requirements.

#### Plan Validation

The validation system checks:

- **Required Fields**: All mandatory plan fields are present
- **Command Validity**: Generated commands are syntactically correct
- **Base Image Consistency**: Base images match detected languages and versions
- **Port Configuration**: Proper port binding and environment variable usage
- **Network Binding**: Commands use `0.0.0.0` for container compatibility

#### Validation Criteria

```go
type ValidationResult struct {
    TestCase string   `json:"testCase"`
    Valid    bool     `json:"valid"`
    Errors   []string `json:"errors"`
    Warnings []string `json:"warnings"`
    Plan     *ExecutionPlan `json:"plan,omitempty"`
}
```

**Validation Rules:**
- Provider must be recognized and valid
- Base image name cannot be empty
- Runtime language must be specified
- Commands must bind to `0.0.0.0` when applicable
- Tools must be from known/supported list
- APT dependencies must be valid package names

## Running Tests

### Prerequisites

```bash
# Install Go (version 1.21 or later)
go version

# Clone the repository
git clone https://github.com/labring/devbox-pack.git
cd devbox-pack

# Install dependencies
go mod download
```

### Unit Tests

Run all unit tests:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./pkg/detector/
go test ./pkg/generators/
go test ./pkg/providers/
```

### Provider-Specific Tests

Test individual providers:

```bash
# Test Node.js provider
go test ./pkg/providers/ -run TestNode

# Test Python provider
go test ./pkg/providers/ -run TestPython

# Test detection engine
go test ./pkg/detector/ -run TestDetect
```

### Integration Tests

Run integration tests with sample repositories:

```bash
# Run integration tests
go test ./railpack/integration_tests/

# Run with verbose output to see detailed results
go test -v ./railpack/integration_tests/
```

### Coverage Analysis

Generate test coverage reports:

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Get coverage summary
go tool cover -func=coverage.out
```

## Test Data and Fixtures

### Sample Projects

The test suite includes sample projects for each supported language:

```
railpack/examples/
├── node-next/              # Next.js application
├── node-vite/              # Vite project
├── python-django/          # Django application
├── python-fastapi/         # FastAPI project
├── java-spring-boot/       # Spring Boot application
├── go-gin/                 # Gin web application
├── php-laravel/            # Laravel application
├── ruby-rails/             # Rails application
├── rust-axum/              # Rust Axum application
├── deno-fresh/             # Deno Fresh application
└── staticfile-spa/         # Static SPA
```

### Test Fixtures

Unit tests use minimal fixtures to test specific scenarios:

```go
// Example fixture for Node.js testing
var nodePackageJson = `{
  "name": "test-app",
  "version": "1.0.0",
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "next": "^13.0.0",
    "react": "^18.0.0"
  }
}`
```

## Writing New Tests

### Provider Tests

When adding a new provider, include these test cases:

1. **Detection Tests**:
   ```go
   func TestNewProviderDetect(t *testing.T) {
       // Test file pattern recognition
       // Test confidence scoring
       // Test version detection
       // Test framework identification
   }
   ```

2. **Command Generation Tests**:
   ```go
   func TestNewProviderCommands(t *testing.T) {
       // Test dev command generation
       // Test build command generation
       // Test start command generation
       // Test environment variable handling
   }
   ```

3. **Edge Case Tests**:
   ```go
   func TestNewProviderEdgeCases(t *testing.T) {
       // Test missing configuration files
       // Test invalid version specifications
       // Test conflicting framework indicators
   }
   ```

### Integration Tests

Add new sample projects to the examples directory:

1. Create project directory: `railpack/examples/[language]-[framework]/`
2. Include all necessary configuration files
3. Add to integration test suite
4. Document expected behavior

### Test Utilities

Use provided test utilities for consistency:

```go
// Create test file info
func createTestFileInfo(path string, isDir bool) types.FileInfo {
    return types.FileInfo{
        Path:  path,
        IsDir: isDir,
        Size:  100,
    }
}

// Create test detection result
func createTestDetectResult(language string, confidence float64) *types.DetectResult {
    return &types.DetectResult{
        Matched:    true,
        Language:   language,
        Confidence: confidence,
        Evidence: types.Evidence{
            Files:  []string{"package.json"},
            Reason: "Test detection",
        },
    }
}
```

## Continuous Integration

### GitHub Actions

The project uses GitHub Actions for automated testing:

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Run Tests
        run: go test -v ./...
      
      - name: Run Integration Tests
        run: go test -v ./railpack/integration_tests/
      
      - name: Generate Coverage
        run: go test -coverprofile=coverage.out ./...
      
      - name: Upload Coverage
        uses: codecov/codecov-action@v3
```

### Test Requirements

All pull requests must:

1. **Pass All Tests**: Unit and integration tests must pass
2. **Maintain Coverage**: Test coverage should not decrease
3. **Include New Tests**: New features require corresponding tests
4. **Follow Conventions**: Tests should follow established patterns

## Performance Testing

### Benchmark Tests

Performance-critical components include benchmark tests:

```go
func BenchmarkDetectionEngine(b *testing.B) {
    engine := detector.NewDetectionEngine()
    files := createLargeFileSet()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        engine.DetectProject("/test/path", files, nil, nil)
    }
}
```

### Load Testing

Test detection performance with large repositories:

- **File Count**: Test with repositories containing 1000+ files
- **Directory Depth**: Test with deeply nested directory structures
- **Concurrent Detection**: Test multiple simultaneous detections

## Debugging Tests

### Verbose Output

Enable verbose test output for debugging:

```bash
# Run with verbose output
go test -v ./pkg/detector/

# Run specific test with debugging
go test -v -run TestDetectProject_NodeExpress ./pkg/detector/
```

### Test Debugging

Use Go's testing tools for debugging:

```go
func TestDebugExample(t *testing.T) {
    // Use t.Logf for debug output
    t.Logf("Debug info: %+v", result)
    
    // Use t.Helper() in utility functions
    t.Helper()
    
    // Use testify for better assertions
    assert.Equal(t, expected, actual)
    require.NotNil(t, result)
}
```

## Best Practices

### Test Organization

1. **Group Related Tests**: Use subtests for related scenarios
2. **Clear Test Names**: Use descriptive test function names
3. **Isolated Tests**: Each test should be independent
4. **Cleanup**: Always clean up test resources

### Test Data

1. **Minimal Fixtures**: Use the smallest possible test data
2. **Realistic Examples**: Test data should reflect real-world usage
3. **Edge Cases**: Include boundary conditions and error scenarios
4. **Version Compatibility**: Test with multiple language versions

### Assertions

1. **Specific Assertions**: Test specific behaviors, not implementation details
2. **Error Testing**: Always test error conditions
3. **Complete Validation**: Verify all aspects of generated plans
4. **Confidence Levels**: Validate detection confidence scores

This comprehensive testing approach ensures the reliability and accuracy of the DevBox Pack system across all supported languages and deployment scenarios.