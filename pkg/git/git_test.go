package git

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewGitHandler(t *testing.T) {
	handler := NewGitHandler()
	if handler == nil {
		t.Fatal("NewGitHandler() returned nil")
	}

	if handler.tempDirs == nil {
		t.Fatal("tempDirs should be initialized")
	}

	if len(handler.tempDirs) != 0 {
		t.Error("tempDirs should be empty initially")
	}
}

func TestParseRepository_Local(t *testing.T) {
	handler := NewGitHandler()

	tests := []struct {
		name     string
		repoPath string
		wantURL  string
		isLocal  bool
	}{
		{
			name:     "current directory",
			repoPath: ".",
			wantURL:  ".",
			isLocal:  true,
		},
		{
			name:     "absolute path",
			repoPath: "/home/user/project",
			wantURL:  "/home/user/project",
			isLocal:  true,
		},
		{
			name:     "relative path",
			repoPath: "../project",
			wantURL:  "../project",
			isLocal:  true,
		},
		{
			name:     "github https url",
			repoPath: "https://github.com/user/repo.git",
			wantURL:  "https://github.com/user/repo.git",
			isLocal:  false,
		},
		{
			name:     "github ssh url",
			repoPath: "git@github.com:user/repo.git",
			wantURL:  "git@github.com:user/repo.git",
			isLocal:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := handler.parseRepository(tt.repoPath)

			if repo == nil {
				t.Fatal("parseRepository returned nil")
			}

			if repo.URL != tt.wantURL {
				t.Errorf("expected URL %s, got %s", tt.wantURL, repo.URL)
			}

			if repo.IsLocal != tt.isLocal {
				t.Errorf("expected IsLocal %v, got %v", tt.isLocal, repo.IsLocal)
			}
		})
	}
}

func TestPrepareLocalProject_ValidPath(t *testing.T) {
	handler := NewGitHandler()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repo := &types.GitRepository{
		URL:     tmpDir,
		IsLocal: true,
	}

	projectPath, err := handler.prepareLocalProject(repo)
	if err != nil {
		t.Fatalf("prepareLocalProject failed: %v", err)
	}

	if projectPath != tmpDir {
		t.Errorf("expected project path %s, got %s", tmpDir, projectPath)
	}
}

func TestPrepareLocalProject_InvalidPath(t *testing.T) {
	handler := NewGitHandler()

	repo := &types.GitRepository{
		URL:     "/non/existent/path",
		IsLocal: true,
	}

	_, err := handler.prepareLocalProject(repo)
	if err == nil {
		t.Error("expected error for non-existent path")
	}

	// Check that it's a DevBoxPackError
	if devboxErr, ok := err.(*types.DevBoxPackError); ok {
		if devboxErr.Code != types.ErrorCodeLocalAccessError {
			t.Errorf("expected error code %s, got %s", types.ErrorCodeLocalAccessError, devboxErr.Code)
		}
	} else {
		t.Error("expected DevBoxPackError")
	}
}

func TestPrepareLocalProject_FilePath(t *testing.T) {
	handler := NewGitHandler()

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "git-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	repo := &types.GitRepository{
		URL:     tmpFile.Name(),
		IsLocal: true,
	}

	_, err = handler.prepareLocalProject(repo)
	if err == nil {
		t.Error("expected error for file path instead of directory")
	}
}

func TestExtractRepoName(t *testing.T) {
	handler := NewGitHandler()

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "github https url",
			url:      "https://github.com/user/repo.git",
			expected: "repo",
		},
		{
			name:     "github https url without .git",
			url:      "https://github.com/user/repo",
			expected: "repo",
		},
		{
			name:     "github ssh url",
			url:      "git@github.com:user/repo.git",
			expected: "repo",
		},
		{
			name:     "gitlab url",
			url:      "https://gitlab.com/user/project.git",
			expected: "project",
		},
		{
			name:     "simple name",
			url:      "myproject",
			expected: "myproject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.extractRepoName(tt.url)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCreateTempDir(t *testing.T) {
	handler := NewGitHandler()

	tempDir, err := handler.createTempDir()
	if err != nil {
		t.Fatalf("createTempDir failed: %v", err)
	}

	// Check that directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("temp directory was not created")
	}

	// Check that it's tracked
	if len(handler.tempDirs) != 1 {
		t.Errorf("expected 1 temp dir tracked, got %d", len(handler.tempDirs))
	}

	if handler.tempDirs[0] != tempDir {
		t.Error("temp dir not tracked correctly")
	}

	// Cleanup
	handler.Cleanup()
}

func TestScanProject(t *testing.T) {
	handler := NewGitHandler()

	// Create a temporary directory with some files
	tmpDir, err := os.MkdirTemp("", "scan-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{
		"package.json",
		"index.js",
		"README.md",
		"src/app.js",
		"src/utils.js",
		"node_modules/package/index.js", // Should be ignored
		".git/config",                   // Should be ignored
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tmpDir, file)
		dir := filepath.Dir(filePath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}

		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", filePath, err)
		}
	}

	options := &types.ScanOptions{
		Depth:    2,
		MaxDepth: 3,
		MaxFiles: 100,
	}

	files, err := handler.ScanProject(tmpDir, options)
	if err != nil {
		t.Fatalf("ScanProject failed: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("no files found")
	}

	// Check that important files are found
	foundFiles := make(map[string]bool)
	for _, file := range files {
		foundFiles[file.Name] = true
	}

	expectedFiles := []string{"package.json", "index.js", "README.md", "app.js", "utils.js"}
	for _, expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("expected file %s not found", expected)
		}
	}

	// Check that ignored files are not found
	ignoredFiles := []string{"config"} // from .git/config
	for _, ignored := range ignoredFiles {
		if foundFiles[ignored] {
			t.Errorf("ignored file %s was found", ignored)
		}
	}
}

func TestFileExists(t *testing.T) {
	handler := NewGitHandler()

	// Create a temporary directory with a test file
	tmpDir, err := os.MkdirTemp("", "file-exists-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := "test.txt"
	testFilePath := filepath.Join(tmpDir, testFile)
	if err := os.WriteFile(testFilePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test existing file
	if !handler.FileExists(tmpDir, testFile) {
		t.Error("FileExists should return true for existing file")
	}

	// Test non-existing file
	if handler.FileExists(tmpDir, "nonexistent.txt") {
		t.Error("FileExists should return false for non-existing file")
	}
}

func TestReadFile(t *testing.T) {
	handler := NewGitHandler()

	// Create a temporary directory with a test file
	tmpDir, err := os.MkdirTemp("", "read-file-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := "test.txt"
	testContent := "Hello, World!"
	testFilePath := filepath.Join(tmpDir, testFile)
	if err := os.WriteFile(testFilePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test reading existing file
	content, err := handler.ReadFile(tmpDir, testFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if content != testContent {
		t.Errorf("expected content %s, got %s", testContent, content)
	}

	// Test reading non-existing file
	_, err = handler.ReadFile(tmpDir, "nonexistent.txt")
	if err == nil {
		t.Error("expected error for non-existing file")
	}
}

func TestReadJSONFile(t *testing.T) {
	handler := NewGitHandler()

	// Create a temporary directory with a JSON file
	tmpDir, err := os.MkdirTemp("", "read-json-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := "package.json"
	testData := map[string]interface{}{
		"name":    "test-project",
		"version": "1.0.0",
		"scripts": map[string]interface{}{
			"start": "node index.js",
		},
	}

	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	testFilePath := filepath.Join(tmpDir, testFile)
	if err := os.WriteFile(testFilePath, jsonData, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test reading JSON file
	var result map[string]interface{}
	err = handler.ReadJSONFile(tmpDir, testFile, &result)
	if err != nil {
		t.Fatalf("ReadJSONFile failed: %v", err)
	}

	if result["name"] != "test-project" {
		t.Errorf("expected name 'test-project', got %v", result["name"])
	}

	if result["version"] != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %v", result["version"])
	}

	// Test reading non-existing file
	err = handler.ReadJSONFile(tmpDir, "nonexistent.json", &result)
	if err == nil {
		t.Error("expected error for non-existing file")
	}
}

func TestReadJSONCFile(t *testing.T) {
	handler := NewGitHandler()

	// Create a temporary directory with a JSONC file
	tmpDir, err := os.MkdirTemp("", "read-jsonc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := "tsconfig.json"
	// JSONC with comments
	testContent := `{
  // TypeScript configuration
  "compilerOptions": {
    "target": "es2020",
    "module": "commonjs"
  },
  /* Multi-line comment */
  "include": ["src/**/*"]
}`

	testFilePath := filepath.Join(tmpDir, testFile)
	if err := os.WriteFile(testFilePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test reading JSONC file
	var result map[string]interface{}
	err = handler.ReadJSONCFile(tmpDir, testFile, &result)
	if err != nil {
		t.Fatalf("ReadJSONCFile failed: %v", err)
	}

	compilerOptions, ok := result["compilerOptions"].(map[string]interface{})
	if !ok {
		t.Fatal("compilerOptions should be a map")
	}

	if compilerOptions["target"] != "es2020" {
		t.Errorf("expected target 'es2020', got %v", compilerOptions["target"])
	}

	if compilerOptions["module"] != "commonjs" {
		t.Errorf("expected module 'commonjs', got %v", compilerOptions["module"])
	}
}

func TestGetFileExtension(t *testing.T) {
	handler := NewGitHandler()

	tests := []struct {
		filename string
		expected string
	}{
		{"file.txt", ".txt"},
		{"package.json", ".json"},
		{"app.js", ".js"},
		{"style.css", ".css"},
		{"README.md", ".md"},
		{"Dockerfile", ""},
		{"file.tar.gz", ".gz"},
		{"", ""},
		{"file", ""},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := handler.getFileExtension(tt.filename)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIsImportantDotFile(t *testing.T) {
	handler := NewGitHandler()

	tests := []struct {
		name     string
		expected bool
	}{
		{".gitignore", true},
		{".env", true},
		{".nvmrc", true},
		{".node-version", true},
		{".python-version", true},
		{".ruby-version", true},
		{".go-version", true},
		{".eslintrc.js", true},
		{".prettierrc", true},
		{".babelrc", true},
		{".DS_Store", false},
		{".git", false},
		{".svn", false},
		{".random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.isImportantDotFile(tt.name)
			if result != tt.expected {
				t.Errorf("expected %v for %s, got %v", tt.expected, tt.name, result)
			}
		})
	}
}

func TestShouldIgnoreDirectory(t *testing.T) {
	handler := NewGitHandler()

	tests := []struct {
		name     string
		expected bool
	}{
		{"node_modules", true},
		{".git", true},
		{".svn", true},
		{"vendor", true},
		{"target", true},
		{"build", true},
		{"dist", true},
		{".next", true},
		{"__pycache__", true},
		{"src", false},
		{"lib", false},
		{"public", false},
		{"assets", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.shouldIgnoreDirectory(tt.name)
			if result != tt.expected {
				t.Errorf("expected %v for %s, got %v", tt.expected, tt.name, result)
			}
		})
	}
}

func TestCleanup(t *testing.T) {
	handler := NewGitHandler()

	// Create some temp directories
	tempDir1, err := handler.createTempDir()
	if err != nil {
		t.Fatalf("Failed to create temp dir 1: %v", err)
	}

	tempDir2, err := handler.createTempDir()
	if err != nil {
		t.Fatalf("Failed to create temp dir 2: %v", err)
	}

	// Verify directories exist
	if _, err := os.Stat(tempDir1); os.IsNotExist(err) {
		t.Error("temp dir 1 should exist")
	}
	if _, err := os.Stat(tempDir2); os.IsNotExist(err) {
		t.Error("temp dir 2 should exist")
	}

	// Cleanup
	err = handler.Cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Verify directories are removed
	if _, err := os.Stat(tempDir1); !os.IsNotExist(err) {
		t.Error("temp dir 1 should be removed")
	}
	if _, err := os.Stat(tempDir2); !os.IsNotExist(err) {
		t.Error("temp dir 2 should be removed")
	}

	// Verify tempDirs is cleared
	if len(handler.tempDirs) != 0 {
		t.Error("tempDirs should be empty after cleanup")
	}
}

func TestWalkFiles(t *testing.T) {
	handler := NewGitHandler()

	// Create a temporary directory with some files
	tmpDir, err := os.MkdirTemp("", "walk-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{
		"file1.txt",
		"file2.js",
		"subdir/file3.json",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tmpDir, file)
		dir := filepath.Dir(filePath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}

		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", filePath, err)
		}
	}

	// Walk files and collect them
	var foundFiles []string
	err = handler.WalkFiles(tmpDir, func(path string, info os.FileInfo) error {
		if !info.IsDir() {
			foundFiles = append(foundFiles, path)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("WalkFiles failed: %v", err)
	}

	if len(foundFiles) != len(testFiles) {
		t.Errorf("expected %d files, got %d", len(testFiles), len(foundFiles))
	}

	// Check that all test files were found
	foundMap := make(map[string]bool)
	for _, file := range foundFiles {
		foundMap[file] = true
	}

	for _, expected := range testFiles {
		if !foundMap[expected] {
			t.Errorf("expected file %s not found", expected)
		}
	}
}
