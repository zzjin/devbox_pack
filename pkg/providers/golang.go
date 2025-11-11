/**
 * DevBox Pack Execution Plan Generator - Go Provider
 */

package providers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/labring/devbox-pack/pkg/types"
)

// GoProvider Go project detector
type GoProvider struct {
	BaseProvider
}

// NewGoProvider creates Go Provider
func NewGoProvider() *GoProvider {
	return &GoProvider{
		BaseProvider: BaseProvider{
			Name:     "go",
			Language: "go",
			Priority: 20,
		},
	}
}

// GetName gets Provider name
func (p *GoProvider) GetName() string {
	return p.Name
}

// GetPriority gets Provider priority
func (p *GoProvider) GetPriority() int {
	return p.Priority
}

// Detect detects Go project
func (p *GoProvider) Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error) {
	indicators := []types.ConfidenceIndicator{
		{Weight: 40, Satisfied: p.HasFile(files, "go.mod")},
		{Weight: 35, Satisfied: p.HasFile(files, "go.work")}, // Higher weight for workspaces
		{Weight: 25, Satisfied: p.HasAnyFile(files, []string{"*.go"})},
		{Weight: 15, Satisfied: p.HasFile(files, "go.sum")},
		{Weight: 15, Satisfied: p.HasFile(files, "go.work.sum")},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"main.go", "cmd/"})},
		{Weight: 5, Satisfied: p.HasFile(files, "vendor/")},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"Makefile", "makefile"})},
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.2 // Lower detection threshold to support go.work projects

	if !detected {
		return p.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Detect if this is a workspace
	isWorkspace := p.HasFile(files, "go.work")

	// Detect version
	version, err := p.detectGoVersion(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	// Detect framework
	framework, err := p.detectFramework(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	// Detect workspace modules
	var workspaceModules []string
	if isWorkspace {
		workspaceModules, err = p.detectWorkspaceModules(projectPath, gitHandler)
		if err != nil {
			return nil, err
		}
	}

	metadata := map[string]interface{}{
		"hasGoMod":         p.HasFile(files, "go.mod"),
		"hasGoSum":         p.HasFile(files, "go.sum"),
		"hasMainGo":        p.HasFile(files, "main.go"),
		"hasVendor":        p.HasFile(files, "vendor"),
		"hasGoWork":        isWorkspace,
		"hasGoWorkSum":     p.HasFile(files, "go.work.sum"),
		"hasMakefile":      p.HasAnyFile(files, []string{"Makefile", "makefile"}),
		"framework":        framework,
		"isWorkspace":      isWorkspace,
		"workspaceModules": workspaceModules,
	}

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	if p.HasFile(files, "go.mod") {
		evidenceFiles = append(evidenceFiles, "go.mod")
	}
	if p.HasFile(files, "go.sum") {
		evidenceFiles = append(evidenceFiles, "go.sum")
	}
	if p.HasFile(files, "go.work") {
		evidenceFiles = append(evidenceFiles, "go.work")
	}
	if p.HasFile(files, "go.work.sum") {
		evidenceFiles = append(evidenceFiles, "go.work.sum")
	}
	if p.HasFile(files, "main.go") {
		evidenceFiles = append(evidenceFiles, "main.go")
	}
	if p.HasAnyFile(files, []string{"Makefile", "makefile"}) {
		if p.HasFile(files, "Makefile") {
			evidenceFiles = append(evidenceFiles, "Makefile")
		} else {
			evidenceFiles = append(evidenceFiles, "makefile")
		}
	}

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected Go project based on: "
	var reasons []string
	if p.HasFile(files, "go.mod") {
		reasons = append(reasons, "go.mod")
	}
	if isWorkspace {
		reasons = append(reasons, "go.work workspace")
		if len(workspaceModules) > 0 {
			reasons = append(reasons, fmt.Sprintf("%d modules", len(workspaceModules)))
		}
	}
	if p.HasAnyFile(files, []string{"*.go"}) {
		reasons = append(reasons, "Go source files")
	}
	if framework != "" {
		reasons = append(reasons, "framework: "+framework)
	}
	if len(reasons) > 0 {
		reason += reasons[0]
		for i := 1; i < len(reasons); i++ {
			reason += ", " + reasons[i]
		}
	}
	evidence.Reason = reason

	return p.CreateDetectResult(
		true,
		confidence,
		"go",
		version,
		framework,
		"go",
		"go",
		metadata,
		evidence,
	), nil
}

// detectGoVersion detects Go version
func (p *GoProvider) detectGoVersion(projectPath string, gitHandler interface{}) (*types.VersionInfo, error) {
	// Read from go.work
	goWorkContent, err := p.SafeReadText(projectPath, "go.work", gitHandler)
	if err != nil {
		// If file doesn't exist, don't return error, continue trying other methods
		if !strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return nil, err
		}
	}
	if goWorkContent != "" {
		re := regexp.MustCompile(`^go\s+(.+)$`)
		lines := strings.Split(goWorkContent, "\n")
		for _, line := range lines {
			matches := re.FindStringSubmatch(strings.TrimSpace(line))
			if len(matches) > 1 {
				return p.CreateVersionInfo(matches[1], "go.work"), nil
			}
		}
	}

	// Read from go.mod
	goModContent, err := p.SafeReadText(projectPath, "go.mod", gitHandler)
	if err != nil {
		// If file doesn't exist, don't return error, continue trying other methods
		if !strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return nil, err
		}
	}
	if goModContent != "" {
		re := regexp.MustCompile(`^go\s+(.+)$`)
		lines := strings.Split(goModContent, "\n")
		for _, line := range lines {
			matches := re.FindStringSubmatch(strings.TrimSpace(line))
			if len(matches) > 1 {
				return p.CreateVersionInfo(matches[1], "go.mod"), nil
			}
		}
	}

	// Read from .go-version
	version, err := p.ParseVersionFromText(
		projectPath,
		".go-version",
		gitHandler,
		regexp.MustCompile(`^(.+)$`),
	)
	if err != nil {
		// If file doesn't exist, don't return error, use default version
		if !strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return nil, err
		}
	}
	if version != "" {
		return p.CreateVersionInfo(p.NormalizeVersion(version), ".go-version"), nil
	}

	// Default version
	return p.CreateVersionInfo("latest", "default"), nil
}

// detectWorkspaceModules detects modules in a Go workspace
func (p *GoProvider) detectWorkspaceModules(projectPath string, gitHandler interface{}) ([]string, error) {
	goWorkContent, err := p.SafeReadText(projectPath, "go.work", gitHandler)
	if err != nil {
		// If file doesn't exist, return empty slice
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return []string{}, nil
		}
		return nil, err
	}
	if goWorkContent == "" {
		return []string{}, nil
	}

	var modules []string
	lines := strings.Split(goWorkContent, "\n")
	inUseBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for use block
		if strings.HasPrefix(line, "use") {
			inUseBlock = true
			// Extract module path from use directive
			modulePath := strings.TrimPrefix(line, "use")
			modulePath = strings.TrimSpace(modulePath)
			// Remove quotes if present
			if strings.HasPrefix(modulePath, `"`) && strings.HasSuffix(modulePath, `"`) {
				modulePath = strings.Trim(modulePath, `"`)
			}
			if modulePath != "" && modulePath != "." {
				modules = append(modules, modulePath)
			}
		} else if inUseBlock && (strings.HasPrefix(line, "}") || line == "") {
			inUseBlock = false
		}
	}

	return modules, nil
}

// detectFramework detects framework
func (p *GoProvider) detectFramework(projectPath string, gitHandler interface{}) (string, error) {
	goModContent, err := p.SafeReadText(projectPath, "go.mod", gitHandler)
	if err != nil {
		// If file doesn't exist, return empty string instead of error
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return "", nil
		}
		return "", err
	}
	if goModContent == "" {
		return "", nil
	}

	frameworkMap := map[string]string{
		"github.com/gin-gonic/gin": "Gin",
		"github.com/gorilla/mux":   "Gorilla Mux",
		"github.com/labstack/echo": "Echo",
		"github.com/gofiber/fiber": "Fiber",
		"github.com/beego/beego":   "Beego",
		"github.com/revel/revel":   "Revel",
		"github.com/astaxie/beego": "Beego",
		"go.uber.org/fx":           "Fx",
		"github.com/spf13/cobra":   "Cobra CLI",
		"github.com/urfave/cli":    "CLI",
		"gorm.io/gorm":             "GORM",
		"github.com/go-gorm/gorm":  "GORM",
	}

	for dependency, framework := range frameworkMap {
		if strings.Contains(goModContent, dependency) {
			return framework, nil
		}
	}

	return "", nil
}

// GenerateCommands generates commands for Go project
func (p *GoProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	// Check if this is a workspace
	isWorkspace := false
	if result.Metadata != nil {
		if workspace, ok := result.Metadata["isWorkspace"].(bool); ok {
			isWorkspace = workspace
		}
	}

	// Setup commands - download dependencies
	if isWorkspace {
		commands.Setup = []string{
			"go work sync",
		}
	} else {
		commands.Setup = []string{
			"go mod download",
		}
	}

	// Development commands
	if result.Framework != "" {
		// For framework projects, try to use hot reload tools
		if isWorkspace {
			commands.Dev = []string{
				"go run ./...",
			}
		} else {
			commands.Dev = []string{
				"go run .",
			}
		}
	} else {
		if isWorkspace {
			commands.Dev = []string{
				"go run ./...",
			}
		} else {
			commands.Dev = []string{
				"go run .",
			}
		}
	}

	// Build commands
	if isWorkspace {
		commands.Build = []string{
			"go work build",
		}
	} else {
		commands.Build = []string{
			"go build -o app .",
		}
	}

	// Run commands
	if isWorkspace {
		commands.Run = []string{
			"go run ./...",
		}
	} else {
		commands.Run = []string{
			"./app",
		}
	}

	return commands
}

// GenerateEnvironment generates environment variables for Go project
func (p *GoProvider) GenerateEnvironment(result *types.DetectResult) map[string]string {
	env := make(map[string]string)

	// Check if this is a workspace
	isWorkspace := false
	if result.Metadata != nil {
		if workspace, ok := result.Metadata["isWorkspace"].(bool); ok {
			isWorkspace = workspace
		}
	}

	// Set Go specific environment variables
	env["GO_ENV"] = "production"
	env["CGO_ENABLED"] = "0"
	env["GOOS"] = "linux"
	env["GOARCH"] = "amd64"

	// Set port for web applications
	env["PORT"] = "8080"

	// Add Go version if available
	if result.Version != "" {
		env["GO_VERSION"] = result.Version
	}

	// Add workspace-specific environment variables
	if isWorkspace {
		env["GOWORK"] = "on"
		// Add workspace modules count if available
		if result.Metadata != nil {
			if modules, ok := result.Metadata["workspaceModules"].([]string); ok {
				env["GO_WORKSPACE_MODULES"] = fmt.Sprintf("%d", len(modules))
			}
		}
	}

	return env
}

// NeedsNativeCompilation checks if Go project needs native compilation
func (p *GoProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Check metadata for CGO usage flag
	if result.Metadata != nil {
		if usesCGO, ok := result.Metadata["usesCGO"].(bool); ok {
			return usesCGO
		}
	}

	// Go projects usually need compilation, but "native compilation" here refers to CGO or other native dependencies
	// Default to false unless CGO is specifically detected
	return false
}
