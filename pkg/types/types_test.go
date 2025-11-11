package types

import (
	"testing"
	"time"
)

func TestExecutionPlan_Structure(t *testing.T) {
	plan := ExecutionPlan{
		Provider: "node",
		Runtime: RuntimeConfig{
			Image:     "node:20-alpine",
			Framework: stringPtr("nextjs"),
		},
		Environment: map[string]string{
			"NODE_ENV": "development",
		},
		Apt: []string{"git", "curl"},
		Commands: Commands{
			Setup: []string{"npm install"},
			Dev:   []string{"npm run dev"},
			Build: []string{"npm run build"},
			Run:   []string{"npm start"},
		},
		Port: 3000,
		Evidence: Evidence{
			Files:  []string{"package.json", "package-lock.json"},
			Reason: "Node.js project detected",
		},
	}

	if plan.Provider != "node" {
		t.Errorf("expected provider 'node', got %s", plan.Provider)
	}

	if plan.Runtime.Image != "node:20-alpine" {
		t.Errorf("expected image 'node:20-alpine', got %s", plan.Runtime.Image)
	}

	if plan.Runtime.Framework == nil || *plan.Runtime.Framework != "nextjs" {
		t.Errorf("expected framework 'nextjs', got %v", plan.Runtime.Framework)
	}

	if plan.Port != 3000 {
		t.Errorf("expected port 3000, got %d", plan.Port)
	}

	if len(plan.Commands.Setup) != 1 || plan.Commands.Setup[0] != "npm install" {
		t.Errorf("expected setup command 'npm install', got %v", plan.Commands.Setup)
	}
}

func TestRuntimeConfig_Structure(t *testing.T) {
	config := RuntimeConfig{
		Image:     "python:3.11-slim",
		Framework: stringPtr("django"),
	}

	if config.Image != "python:3.11-slim" {
		t.Errorf("expected image 'python:3.11-slim', got %s", config.Image)
	}

	if config.Framework == nil || *config.Framework != "django" {
		t.Errorf("expected framework 'django', got %v", config.Framework)
	}

	// Test with nil framework
	configNoFramework := RuntimeConfig{
		Image: "node:20-alpine",
	}

	if configNoFramework.Framework != nil {
		t.Errorf("expected nil framework, got %v", configNoFramework.Framework)
	}
}

func TestEvidence_Structure(t *testing.T) {
	evidence := Evidence{
		Files:  []string{"requirements.txt", "manage.py"},
		Reason: "Python Django project detected",
	}

	if len(evidence.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(evidence.Files))
	}

	if evidence.Files[0] != "requirements.txt" {
		t.Errorf("expected first file 'requirements.txt', got %s", evidence.Files[0])
	}

	if evidence.Reason != "Python Django project detected" {
		t.Errorf("expected specific reason, got %s", evidence.Reason)
	}
}

func TestDetectResult_Structure(t *testing.T) {
	packageManager := &PackageManager{
		Name:     "npm",
		LockFile: stringPtr("package-lock.json"),
	}

	result := DetectResult{
		Matched:    true,
		Provider:   stringPtr("node"),
		Confidence: 0.95,
		Evidence: Evidence{
			Files:  []string{"package.json"},
			Reason: "Node.js project",
		},
		Language:       "javascript",
		Framework:      "react",
		Version:        "18.17.0",
		PackageManager: packageManager,
		BuildTools:     []string{"webpack", "babel"},
		Metadata: map[string]interface{}{
			"hasTypeScript": true,
		},
	}

	if !result.Matched {
		t.Error("expected Matched to be true")
	}

	if result.Provider == nil || *result.Provider != "node" {
		t.Errorf("expected provider 'node', got %v", result.Provider)
	}

	if result.Confidence != 0.95 {
		t.Errorf("expected confidence 0.95, got %f", result.Confidence)
	}

	if result.PackageManager.Name != "npm" {
		t.Errorf("expected package manager 'npm', got %s", result.PackageManager.Name)
	}

	if len(result.BuildTools) != 2 {
		t.Errorf("expected 2 build tools, got %d", len(result.BuildTools))
	}
}

func TestCLIOptions_Structure(t *testing.T) {
	options := CLIOptions{
		Repository: "https://github.com/user/repo.git",
		Ref:        stringPtr("main"),
		Subdir:     stringPtr("frontend"),
		Provider:   stringPtr("node"),
		Format:     "json",
		Verbose:    true,
		Offline:    false,
		Platform:   stringPtr("linux/amd64"),
		Base:       stringPtr("alpine"),
		Help:       false,
		Version:    false,
		Quiet:      false,
		Pretty:     true,
	}

	if options.Repository != "https://github.com/user/repo.git" {
		t.Errorf("expected specific repository, got %s", options.Repository)
	}

	if options.Ref == nil || *options.Ref != "main" {
		t.Errorf("expected ref 'main', got %v", options.Ref)
	}

	if !options.Verbose {
		t.Error("expected Verbose to be true")
	}

	if !options.Pretty {
		t.Error("expected Pretty to be true")
	}
}

func TestGitRepository_Structure(t *testing.T) {
	repo := GitRepository{
		URL:     "https://github.com/user/repo.git",
		Ref:     stringPtr("develop"),
		Subdir:  stringPtr("backend"),
		IsLocal: false,
	}

	if repo.URL != "https://github.com/user/repo.git" {
		t.Errorf("expected specific URL, got %s", repo.URL)
	}

	if repo.Ref == nil || *repo.Ref != "develop" {
		t.Errorf("expected ref 'develop', got %v", repo.Ref)
	}

	if repo.IsLocal {
		t.Error("expected IsLocal to be false")
	}
}

func TestFileInfo_Structure(t *testing.T) {
	now := time.Now()
	size := int64(1024)
	ext := ".js"

	fileInfo := FileInfo{
		Path:        "src/index.js",
		Name:        "index.js",
		IsDirectory: false,
		Size:        &size,
		Extension:   &ext,
		MTime:       &now,
	}

	if fileInfo.Path != "src/index.js" {
		t.Errorf("expected path 'src/index.js', got %s", fileInfo.Path)
	}

	if fileInfo.Name != "index.js" {
		t.Errorf("expected name 'index.js', got %s", fileInfo.Name)
	}

	if fileInfo.IsDirectory {
		t.Error("expected IsDirectory to be false")
	}

	if fileInfo.Size == nil || *fileInfo.Size != 1024 {
		t.Errorf("expected size 1024, got %v", fileInfo.Size)
	}

	if fileInfo.Extension == nil || *fileInfo.Extension != ".js" {
		t.Errorf("expected extension '.js', got %v", fileInfo.Extension)
	}
}

func TestVersionInfo_Structure(t *testing.T) {
	sourceDetail := "package.json engines field"

	versionInfo := VersionInfo{
		Version:      "18.17.0",
		Source:       "file",
		SourceDetail: &sourceDetail,
	}

	if versionInfo.Version != "18.17.0" {
		t.Errorf("expected version '18.17.0', got %s", versionInfo.Version)
	}

	if versionInfo.Source != "file" {
		t.Errorf("expected source 'file', got %s", versionInfo.Source)
	}

	if versionInfo.SourceDetail == nil || *versionInfo.SourceDetail != sourceDetail {
		t.Errorf("expected source detail '%s', got %v", sourceDetail, versionInfo.SourceDetail)
	}
}

func TestPackageManager_Structure(t *testing.T) {
	lockFile := "yarn.lock"
	configFile := ".yarnrc.yml"

	pm := PackageManager{
		Name:        "yarn",
		LockFile:    &lockFile,
		ConfigFile:  &configFile,
		UseCorepack: true,
	}

	if pm.Name != "yarn" {
		t.Errorf("expected name 'yarn', got %s", pm.Name)
	}

	if pm.LockFile == nil || *pm.LockFile != "yarn.lock" {
		t.Errorf("expected lock file 'yarn.lock', got %v", pm.LockFile)
	}

	if pm.ConfigFile == nil || *pm.ConfigFile != ".yarnrc.yml" {
		t.Errorf("expected config file '.yarnrc.yml', got %v", pm.ConfigFile)
	}

	if !pm.UseCorepack {
		t.Error("expected UseCorepack to be true")
	}
}

func TestDevBoxPackError_Error(t *testing.T) {
	err := &DevBoxPackError{
		Message: "Test error message",
		Code:    "TEST_ERROR",
		Details: map[string]string{"key": "value"},
	}

	errorString := err.Error()
	expectedString := "[TEST_ERROR] Test error message"
	if errorString != expectedString {
		t.Errorf("expected error message '%s', got %s", expectedString, errorString)
	}
}

func TestNewDevBoxPackError(t *testing.T) {
	details := map[string]interface{}{
		"file": "package.json",
		"line": 10,
	}

	err := NewDevBoxPackError("Invalid JSON format", ErrorCodeJSONParseError, details)

	if err == nil {
		t.Fatal("NewDevBoxPackError returned nil")
	}

	if err.Message != "Invalid JSON format" {
		t.Errorf("expected message 'Invalid JSON format', got %s", err.Message)
	}

	if err.Code != ErrorCodeJSONParseError {
		t.Errorf("expected code %s, got %s", ErrorCodeJSONParseError, err.Code)
	}

	if err.Details == nil {
		t.Error("expected Details to be set")
	}
}

func TestConfidenceIndicator_Structure(t *testing.T) {
	indicator := ConfidenceIndicator{
		Weight:    50,
		Satisfied: true,
	}

	if indicator.Weight != 50 {
		t.Errorf("expected weight 50, got %d", indicator.Weight)
	}

	if !indicator.Satisfied {
		t.Error("expected Satisfied to be true")
	}
}

func TestCommands_Structure(t *testing.T) {
	commands := Commands{
		Setup: []string{"npm install", "npm run prepare"},
		Dev:   []string{"npm run dev"},
		Build: []string{"npm run build", "npm run test"},
		Run:   []string{"npm start"},
	}

	if len(commands.Setup) != 2 {
		t.Errorf("expected 2 setup commands, got %d", len(commands.Setup))
	}

	if commands.Setup[0] != "npm install" {
		t.Errorf("expected first setup command 'npm install', got %s", commands.Setup[0])
	}

	if len(commands.Dev) != 1 {
		t.Errorf("expected 1 dev command, got %d", len(commands.Dev))
	}

	if commands.Dev[0] != "npm run dev" {
		t.Errorf("expected dev command 'npm run dev', got %s", commands.Dev[0])
	}
}

func TestScanOptions_Structure(t *testing.T) {
	options := ScanOptions{
		Depth:    2,
		MaxDepth: 5,
		MaxFiles: 1000,
	}

	if options.Depth != 2 {
		t.Errorf("expected depth 2, got %d", options.Depth)
	}

	if options.MaxDepth != 5 {
		t.Errorf("expected max depth 5, got %d", options.MaxDepth)
	}

	if options.MaxFiles != 1000 {
		t.Errorf("expected max files 1000, got %d", options.MaxFiles)
	}
}

func TestErrorCodes(t *testing.T) {
	expectedCodes := map[string]string{
		"GIT_ERROR":           ErrorCodeGitError,
		"LOCAL_ACCESS_ERROR":  ErrorCodeLocalAccessError,
		"INVALID_PATH":        ErrorCodeInvalidPath,
		"CLONE_ERROR":         ErrorCodeCloneError,
		"GIT_CHECKOUT_ERROR":  ErrorCodeGitCheckoutError,
		"SUBDIR_ACCESS_ERROR": ErrorCodeSubdirAccessError,
		"SUBDIR_NOT_FOUND":    ErrorCodeSubdirNotFound,
		"TEMP_DIR_ERROR":      ErrorCodeTempDirError,
		"FILE_READ_ERROR":     ErrorCodeFileReadError,
		"JSON_PARSE_ERROR":    ErrorCodeJSONParseError,
		"INVALID_FORMAT":      ErrorCodeInvalidFormat,
		"INVALID_PLATFORM":    ErrorCodeInvalidPlatform,
		"INVALID_GIT_URL":     ErrorCodeInvalidGitURL,
		"INVALID_INPUT":       ErrorCodeInvalidInput,
		"SCAN_ERROR":          ErrorCodeScanError,
		"INVALID_PROVIDER":    ErrorCodeInvalidProvider,
		"INVALID_ARGUMENT":    ErrorCodeInvalidArgument,
	}

	for expectedValue, actualConstant := range expectedCodes {
		if actualConstant != expectedValue {
			t.Errorf("expected error code constant to equal %s, got %s", expectedValue, actualConstant)
		}
	}
}

func TestSupportedLanguages(t *testing.T) {
	expectedLanguages := map[SupportedLanguage]string{
		LanguageNode:       "node",
		LanguagePython:     "python",
		LanguageJava:       "java",
		LanguageGo:         "go",
		LanguagePHP:        "php",
		LanguageRuby:       "ruby",
		LanguageDeno:       "deno",
		LanguageRust:       "rust",
		LanguageStaticfile: "staticfile",
		LanguageShell:      "shell",
	}

	for constant, expectedValue := range expectedLanguages {
		if string(constant) != expectedValue {
			t.Errorf("expected language constant %s to equal %s, got %s", constant, expectedValue, string(constant))
		}
	}
}

func TestOutputFormats(t *testing.T) {
	expectedFormats := map[OutputFormat]string{
		OutputFormatJSON:   "json",
		OutputFormatPretty: "pretty",
	}

	for constant, expectedValue := range expectedFormats {
		if string(constant) != expectedValue {
			t.Errorf("expected format constant %s to equal %s, got %s", constant, expectedValue, string(constant))
		}
	}
}

func TestPlatforms(t *testing.T) {
	expectedPlatforms := map[Platform]string{
		PlatformLinuxAMD64:  "linux/amd64",
		PlatformLinuxARM64:  "linux/arm64",
		PlatformDarwinAMD64: "darwin/amd64",
		PlatformDarwinARM64: "darwin/arm64",
	}

	for constant, expectedValue := range expectedPlatforms {
		if string(constant) != expectedValue {
			t.Errorf("expected platform constant %s to equal %s, got %s", constant, expectedValue, string(constant))
		}
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
