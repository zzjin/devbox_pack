// Package detector provides interfaces and engine for detecting project types
// and generating execution plans for the DevBox Pack system.
package detector

import (
	"github.com/labring/devbox-pack/pkg/types"
)

// Provider interface definition
type Provider interface {
	GetName() string
	GetLanguage() string
	GetPriority() int
	Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error)
	GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands
	GenerateEnvironment(result *types.DetectResult) map[string]string
	NeedsNativeCompilation(result *types.DetectResult) bool
}
