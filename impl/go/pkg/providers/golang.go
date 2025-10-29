/**
 * DevBox Pack Execution Plan Generator - Go Provider
 */

package providers

import (
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
			Priority: 75,
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
		{Weight: 30, Satisfied: p.HasFile(files, "go.work")}, // Increase go.work weight
		{Weight: 25, Satisfied: p.HasAnyFile(files, []string{"*.go"})},
		{Weight: 15, Satisfied: p.HasFile(files, "go.sum")},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"main.go", "cmd/"})},
		{Weight: 5, Satisfied: p.HasFile(files, "vendor/")},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"Makefile", "makefile"})},
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.2 // Lower detection threshold to support go.work projects

	if !detected {
		return p.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

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

	metadata := map[string]interface{}{
		"hasGoMod":    p.HasFile(files, "go.mod"),
		"hasGoSum":    p.HasFile(files, "go.sum"),
		"hasMainGo":   p.HasFile(files, "main.go"),
		"hasVendor":   p.HasFile(files, "vendor"),
		"hasGoWork":   p.HasFile(files, "go.work"),
		"hasMakefile": p.HasAnyFile(files, []string{"Makefile", "makefile"}),
		"framework":   framework,
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
	if p.HasFile(files, "go.work") {
		reasons = append(reasons, "go.work")
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
	return p.CreateVersionInfo("1.21", "default"), nil
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

	commands.Dev = []string{
		"go mod download",
		"go run .",
	}

	commands.Build = []string{
		"go mod download",
		"go build -o app .",
	}

	commands.Start = []string{
		"./app",
	}

	return commands
}

// NeedsNativeCompilation checks if Go project needs native compilation
func (p *GoProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Go projects usually need compilation, but "native compilation" here might refer to CGO or other native dependencies
	// Check if CGO or other dependencies requiring C compiler are used

	// Simplified handling: Go projects usually don't need additional native compilation steps
	// Unless there are special CGO dependencies, but this requires more complex detection logic
	return false
}
