// Package git provides Git repository handling functionality
// for the DevBox Pack execution plan generator.
package git

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/labring/devbox-pack/pkg/types"
)

// Scan configuration constants
const (
	DefaultDepth = 3
	MaxFiles     = 1000
)

// GitHandler Git repository handler
type GitHandler struct {
	tempDirs []string
}

// NewGitHandler creates a new Git handler instance
func NewGitHandler() *GitHandler {
	return &GitHandler{
		tempDirs: make([]string, 0),
	}
}

// execGit executes Git commands
func (g *GitHandler) execGit(args []string, cwd string) (string, error) {
	cmd := exec.Command("git", args...)
	if cwd != "" {
		cmd.Dir = cwd
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", types.NewDevBoxPackError(
			fmt.Sprintf("Git operation failed: %s", string(output)),
			types.ErrorCodeGitError,
			map[string]interface{}{
				"command": strings.Join(args, " "),
				"output":  string(output),
			},
		)
	}

	return strings.TrimSpace(string(output)), nil
}

// PrepareProject prepares project directory (local path or remote repository)
func (g *GitHandler) PrepareProject(repoPath string) (string, error) {
	repo := g.parseRepository(repoPath)

	if repo.IsLocal {
		return g.prepareLocalProject(repo)
	}
	return g.cloneRepository(repo)
}

// parseRepository parses repository path
func (g *GitHandler) parseRepository(repoPath string) *types.GitRepository {
	// Check if it's a local path
	if !strings.HasPrefix(repoPath, "http") && !strings.Contains(repoPath, "@") {
		return &types.GitRepository{
			URL:     repoPath,
			IsLocal: true,
		}
	}

	// Parse remote repository URL
	return &types.GitRepository{
		URL:     repoPath,
		IsLocal: false,
	}
}

// prepareLocalProject prepares local project
func (g *GitHandler) prepareLocalProject(repo *types.GitRepository) (string, error) {
	projectPath := repo.URL

	// Check if path exists
	stat, err := os.Stat(projectPath)
	if err != nil {
		return "", types.NewDevBoxPackError(
			fmt.Sprintf("cannot access local project: %s", err.Error()),
			types.ErrorCodeLocalAccessError,
			nil,
		)
	}

	if !stat.IsDir() {
		return "", types.NewDevBoxPackError(
			fmt.Sprintf("path is not a directory: %s", projectPath),
			types.ErrorCodeInvalidPath,
			nil,
		)
	}

	// Check if it's a Git repository (optional)
	_, err = g.execGit([]string{"rev-parse", "--git-dir"}, projectPath)
	if err != nil {
		// It's okay if it's not a Git repository, continue processing
	}

	return projectPath, nil
}

// cloneRepository clones remote repository
func (g *GitHandler) cloneRepository(repo *types.GitRepository) (string, error) {
	tempDir, err := g.createTempDir()
	if err != nil {
		return "", err
	}

	repoName := g.extractRepoName(repo.URL)
	clonePath := filepath.Join(tempDir, repoName)

	// Build clone command arguments
	cloneArgs := []string{"clone", "--depth", "1", repo.URL, clonePath}

	// Execute clone
	_, err = g.execGit(cloneArgs, "")
	if err != nil {
		g.cleanupTempDir(tempDir)
		return "", types.NewDevBoxPackError(
			fmt.Sprintf("repository clone failed: %s", err.Error()),
			types.ErrorCodeCloneError,
			map[string]interface{}{"url": repo.URL},
		)
	}

	// If ref is specified, switch to corresponding branch/tag
	if repo.Ref != nil {
		// First try to fetch all branches and tags
		_, err = g.execGit([]string{"fetch", "--all", "--tags"}, clonePath)
		if err != nil {
			g.cleanupTempDir(tempDir)
			return "", err
		}

		// Switch to specified ref
		_, err = g.execGit([]string{"checkout", *repo.Ref}, clonePath)
		if err != nil {
			g.cleanupTempDir(tempDir)
			return "", types.NewDevBoxPackError(
				fmt.Sprintf("cannot switch to specified ref: %s", *repo.Ref),
				types.ErrorCodeGitCheckoutError,
				map[string]interface{}{
					"ref":   *repo.Ref,
					"error": err.Error(),
				},
			)
		}
	}

	// If subdirectory is specified, return subdirectory path
	if repo.Subdir != nil {
		subdirPath := filepath.Join(clonePath, *repo.Subdir)
		stat, err := os.Stat(subdirPath)
		if err != nil {
			g.cleanupTempDir(tempDir)
			return "", types.NewDevBoxPackError(
				fmt.Sprintf("cannot access subdirectory: %s", *repo.Subdir),
				types.ErrorCodeSubdirAccessError,
				nil,
			)
		}
		if !stat.IsDir() {
			g.cleanupTempDir(tempDir)
			return "", types.NewDevBoxPackError(
				fmt.Sprintf("subdirectory does not exist: %s", *repo.Subdir),
				types.ErrorCodeSubdirNotFound,
				nil,
			)
		}
		return subdirPath, nil
	}

	return clonePath, nil
}

// extractRepoName extracts repository name from URL
func (g *GitHandler) extractRepoName(url string) string {
	// Handle SSH format: git@github.com:user/repo.git
	if strings.Contains(url, "@") && strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		if len(parts) > 0 {
			repoPath := parts[len(parts)-1]
			return strings.TrimSuffix(filepath.Base(repoPath), ".git")
		}
	}

	// Handle HTTPS format: https://github.com/user/repo.git
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		repoName := parts[len(parts)-1]
		return strings.TrimSuffix(repoName, ".git")
	}

	// If parsing fails, use default name
	return "repo"
}

// createTempDir creates temporary directory
func (g *GitHandler) createTempDir() (string, error) {
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("devbox-pack-%d-%d", time.Now().UnixNano(), os.Getpid()))
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		return "", types.NewDevBoxPackError(
			fmt.Sprintf("failed to create temporary directory: %s", err.Error()),
			types.ErrorCodeTempDirError,
			nil,
		)
	}
	g.tempDirs = append(g.tempDirs, tempDir)
	return tempDir, nil
}

// ScanProject scans project files
func (g *GitHandler) ScanProject(projectPath string, options *types.ScanOptions) ([]*types.FileInfo, error) {
	if options == nil {
		options = &types.ScanOptions{
			Depth:    DefaultDepth,
			MaxFiles: MaxFiles,
		}
	}

	files := make([]*types.FileInfo, 0)
	err := g.scanDirectory(projectPath, "", &files, 0, options.Depth, options.MaxFiles)
	if err != nil {
		return nil, types.NewDevBoxPackError(
			fmt.Sprintf("Project scan failed: %s", err.Error()),
			types.ErrorCodeScanError,
			nil,
		)
	}

	return files, nil
}

// scanDirectory recursively scans directory
func (g *GitHandler) scanDirectory(basePath, currentPath string, files *[]*types.FileInfo, currentDepth, maxDepth, maxFiles int) error {
	if len(*files) >= maxFiles || currentDepth > maxDepth {
		return nil
	}

	fullPath := filepath.Join(basePath, currentPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		// Ignore inaccessible directories
		return nil
	}

	for _, entry := range entries {
		if len(*files) >= maxFiles {
			break
		}

		entryPath := currentPath
		if entryPath != "" {
			entryPath = filepath.Join(entryPath, entry.Name())
		} else {
			entryPath = entry.Name()
		}

		if entry.IsDir() {
			// Skip directories that should be ignored
			if g.shouldIgnoreDirectory(entry.Name()) {
				continue
			}

			// Recursively scan subdirectories
			err = g.scanDirectory(basePath, entryPath, files, currentDepth+1, maxDepth, maxFiles)
			if err != nil {
				return err
			}
		} else {
			// Skip hidden files unless they are important configuration files
			if strings.HasPrefix(entry.Name(), ".") && !g.isImportantDotFile(entry.Name()) {
				continue
			}

			entryFullPath := filepath.Join(basePath, entryPath)
			stat, err := os.Stat(entryFullPath)
			if err != nil {
				// Ignore inaccessible files
				continue
			}

			size := stat.Size()
			ext := g.getFileExtension(entry.Name())
			*files = append(*files, &types.FileInfo{
				Path:        entryPath,
				Name:        entry.Name(),
				Size:        &size,
				IsDirectory: false,
				Extension:   &ext,
			})
		}
	}

	return nil
}

// getFileExtension gets file extension
func (g *GitHandler) getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		return strings.ToLower(ext[1:]) // Remove dot and convert to lowercase
	}
	return ""
}

// isImportantDotFile checks if it's an important dot file
func (g *GitHandler) isImportantDotFile(name string) bool {
	importantDotFiles := []string{".gitignore", ".dockerignore", ".env", ".env.example", ".nvmrc", ".python-version"}
	for _, file := range importantDotFiles {
		if name == file {
			return true
		}
	}
	return false
}

// shouldIgnoreDirectory checks if directory should be ignored
func (g *GitHandler) shouldIgnoreDirectory(name string) bool {
	ignoredDirs := []string{"node_modules", ".git", "dist", "build", "__pycache__", ".pytest_cache", "target", "vendor"}
	for _, dir := range ignoredDirs {
		if name == dir {
			return true
		}
	}
	return false
}

// FileExists checks if file exists
func (g *GitHandler) FileExists(projectPath, filePath string) bool {
	fullPath := filepath.Join(projectPath, filePath)
	_, err := os.Stat(fullPath)
	return err == nil
}

// ReadFile reads file content
func (g *GitHandler) ReadFile(projectPath, filePath string) (string, error) {
	fullPath := filepath.Join(projectPath, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", types.NewDevBoxPackError(
			fmt.Sprintf("Failed to read file: %s", filePath),
			types.ErrorCodeFileReadError,
			map[string]interface{}{
				"path":  filePath,
				"error": err.Error(),
			},
		)
	}
	return string(content), nil
}

// ReadJSONFile reads JSON file
func (g *GitHandler) ReadJSONFile(projectPath, filePath string, v interface{}) error {
	content, err := g.ReadFile(projectPath, filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(content), v)
	if err != nil {
		return types.NewDevBoxPackError(
			fmt.Sprintf("Failed to parse JSON file: %s", filePath),
			types.ErrorCodeJSONParseError,
			map[string]interface{}{
				"path":  filePath,
				"error": err.Error(),
			},
		)
	}

	return nil
}

// ReadJSONCFile reads JSONC file (JSON with comments support)
func (g *GitHandler) ReadJSONCFile(projectPath, filePath string, v interface{}) error {
	content, err := g.ReadFile(projectPath, filePath)
	if err != nil {
		return err
	}

	// Simple JSONC processing: remove single-line comments
	lines := strings.Split(content, "\n")
	var cleanedLines []string
	for _, line := range lines {
		// Remove // comments
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		cleanedLines = append(cleanedLines, line)
	}
	cleanedContent := strings.Join(cleanedLines, "\n")

	err = json.Unmarshal([]byte(cleanedContent), v)
	if err != nil {
		return types.NewDevBoxPackError(
			fmt.Sprintf("Failed to parse JSONC file: %s", filePath),
			types.ErrorCodeJSONParseError,
			map[string]interface{}{
				"path":  filePath,
				"error": err.Error(),
			},
		)
	}

	return nil
}

// Cleanup cleans up temporary directories
func (g *GitHandler) Cleanup() error {
	for _, tempDir := range g.tempDirs {
		g.cleanupTempDir(tempDir)
	}
	g.tempDirs = g.tempDirs[:0]
	return nil
}

// cleanupTempDir cleans up single temporary directory
func (g *GitHandler) cleanupTempDir(tempDir string) {
	_ = os.RemoveAll(tempDir) // Ignore cleanup errors
}

// WalkFiles walks through all files in directory
func (g *GitHandler) WalkFiles(projectPath string, walkFn func(path string, info fs.FileInfo) error) error {
	return filepath.Walk(projectPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip directories that should be ignored
			if g.shouldIgnoreDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip hidden files unless they are important configuration files
		if strings.HasPrefix(info.Name(), ".") && !g.isImportantDotFile(info.Name()) {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			return err
		}

		return walkFn(relPath, info)
	})
}
