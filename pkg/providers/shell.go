/**
 * DevBox Pack Execution Plan Generator - Shell Provider
 */

package providers

import (
	"github.com/labring/devbox-pack/pkg/types"
)

// ShellProvider Shell project detector
type ShellProvider struct {
	BaseProvider
}

// NewShellProvider creates Shell Provider
func NewShellProvider() *ShellProvider {
	return &ShellProvider{
		BaseProvider: BaseProvider{
			Name:     "shell",
			Language: "shell",
			Priority: 100,
		},
	}
}

// GetName gets Provider name
func (p *ShellProvider) GetName() string {
	return p.Name
}

// GetPriority gets Provider priority
func (p *ShellProvider) GetPriority() int {
	return p.Priority
}

// Detect detects Shell project
func (p *ShellProvider) Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error) {
	indicators := []types.ConfidenceIndicator{
		{Weight: 30, Satisfied: p.HasAnyFile(files, []string{"*.sh", "*.bash", "*.zsh"})},
		{Weight: 20, Satisfied: p.HasAnyFile(files, []string{"Makefile", "makefile"})},
		{Weight: 15, Satisfied: p.HasAnyFile(files, []string{"*.fish", "*.csh", "*.ksh"})},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"install.sh", "setup.sh", "build.sh"})},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"bin/", "scripts/"})},
		{Weight: 5, Satisfied: p.HasFile(files, "README.md")},
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.3

	if !detected {
		return p.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Shell projects usually don't have specific versions
	version := p.CreateVersionInfo("latest", "default")

	// Detect Shell type
	shellType := p.detectShellType(files)

	// Detect project type
	projectType := p.detectProjectType(files)

	metadata := map[string]interface{}{
		"shellType":   shellType,
		"projectType": projectType,
		"shellFiles":  p.getShellFiles(files),
		"hasShebang":  len(p.getShellFilesWithShebang(files, shellType)) > 0,
	}

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	shellFiles := p.getShellFiles(files)
	evidenceFiles = append(evidenceFiles, shellFiles...)

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected Shell project based on: "
	var reasons []string
	if len(shellFiles) > 0 {
		reasons = append(reasons, "Shell script files")
	}
	if shellType != "" {
		reasons = append(reasons, "shell type: "+shellType)
	}
	if projectType != "" {
		reasons = append(reasons, "project type: "+projectType)
	}
	if len(p.getShellFilesWithShebang(files, shellType)) > 0 {
		reasons = append(reasons, "files with shebang")
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
		"shell",
		version,
		projectType,
		"",
		"",
		metadata,
		evidence,
	), nil
}

// detectShellType detects Shell type
func (p *ShellProvider) detectShellType(files []types.FileInfo) string {
	if p.HasAnyFile(files, []string{"*.bash"}) {
		return "bash"
	}
	if p.HasAnyFile(files, []string{"*.zsh"}) {
		return "zsh"
	}
	if p.HasAnyFile(files, []string{"*.fish"}) {
		return "fish"
	}
	if p.HasAnyFile(files, []string{"*.sh"}) {
		return "sh"
	}
	return "sh"
}

// detectProjectType detects project type
func (p *ShellProvider) detectProjectType(files []types.FileInfo) string {
	if p.HasAnyFile(files, []string{"install.sh", "setup.sh", "installer.sh"}) {
		return "installer"
	}

	if p.HasAnyFile(files, []string{"build.sh", "compile.sh", "make.sh"}) {
		return "build-tool"
	}

	if p.HasAnyFile(files, []string{"deploy.sh", "release.sh", "publish.sh"}) {
		return "deployment"
	}

	if p.HasAnyFile(files, []string{"backup.sh", "restore.sh", "sync.sh"}) {
		return "utility"
	}

	if p.HasAnyFile(files, []string{"start.sh", "stop.sh", "restart.sh", "service.sh"}) {
		return "service"
	}

	if p.HasAnyFile(files, []string{"test.sh", "check.sh", "validate.sh"}) {
		return "testing"
	}

	return "script"
}

// getShellFiles gets all Shell files
func (p *ShellProvider) getShellFiles(files []types.FileInfo) []string {
	shellExtensions := []string{"*.sh", "*.bash", "*.zsh", "*.fish"}
	var shellFiles []string

	for _, extension := range shellExtensions {
		matchingFiles := p.GetMatchingFiles(files, extension)
		for _, f := range matchingFiles {
			shellFiles = append(shellFiles, f.Path)
		}
	}

	return shellFiles
}

// getShellFilesWithShebang gets Shell files containing specific shebang
func (p *ShellProvider) getShellFilesWithShebang(files []types.FileInfo, shellType string) []string {
	// Simplified handling here, should actually read file content to check shebang
	// For performance reasons, only judge by file extension here
	extensionMap := map[string][]string{
		"bash": {"*.bash"},
		"zsh":  {"*.zsh"},
		"fish": {"*.fish"},
		"sh":   {"*.sh"},
	}

	extensions, exists := extensionMap[shellType]
	if !exists {
		return []string{}
	}

	var matchingFiles []string
	for _, extension := range extensions {
		filesMatching := p.GetMatchingFiles(files, extension)
		for _, f := range filesMatching {
			matchingFiles = append(matchingFiles, f.Path)
		}
	}

	return matchingFiles
}

// GenerateCommands generates commands for Shell project
func (p *ShellProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	// Setup commands - make scripts executable
	commands.Setup = []string{"chmod +x *.sh"}

	// Development and Run commands
	commands.Dev = []string{"./start.sh"}
	commands.Run = []string{"./start.sh"}

	return commands
}

// GenerateEnvironment generates environment variables for Shell project
func (p *ShellProvider) GenerateEnvironment(result *types.DetectResult) map[string]string {
	env := make(map[string]string)

	// Set Shell specific environment variables
	env["SHELL_ENV"] = "production"

	// Set common environment variables
	env["PATH"] = "/usr/local/bin:/usr/bin:/bin"

	return env
}

// NeedsNativeCompilation checks if Shell project needs native compilation
func (p *ShellProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Shell scripts don't need compilation
	return false
}
