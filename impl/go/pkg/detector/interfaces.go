/**
 * DevBox Pack Execution Plan Generator - Provider Interface Definition
 */

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
	NeedsNativeCompilation(result *types.DetectResult) bool
}
