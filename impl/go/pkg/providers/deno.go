/**
 * DevBox Pack Execution Plan Generator - Deno Provider
 */

package providers

import (
	"strings"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

// DenoProvider Deno project detector
type DenoProvider struct {
	BaseProvider
}

// NewDenoProvider creates Deno Provider
func NewDenoProvider() *DenoProvider {
	return &DenoProvider{
		BaseProvider: BaseProvider{
			Name:     "deno",
			Language: "deno",
			Priority: 95,
		},
	}
}

// GetName gets Provider name
func (p *DenoProvider) GetName() string {
	return p.Name
}

// GetPriority gets Provider priority
func (p *DenoProvider) GetPriority() int {
	return p.Priority
}

// Detect detects Deno project
func (p *DenoProvider) Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error) {
	gh := gitHandler.(*git.GitHandler)

	// Check if Staticfile or go.work exists, if so, don't detect as Deno project
	if p.HasFile(files, "Staticfile") || p.HasFile(files, "go.work") {
		return p.CreateDetectResult(false, 0, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	indicators := []types.ConfidenceIndicator{
		{Weight: 40, Satisfied: p.HasAnyFile(files, []string{"deno.json", "deno.jsonc"})},
		{Weight: 25, Satisfied: p.HasAnyFile(files, []string{"*.ts", "*.js"})},
		{Weight: 15, Satisfied: p.HasFile(files, "deno.lock")},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"deps.ts", "mod.ts"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"import_map.json"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"main.ts", "app.ts", "server.ts"})},
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.2 // Lower detection threshold

	if !detected {
		return p.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Detect version
	version, err := p.detectDenoVersion(projectPath, gh)
	if err != nil {
		return nil, err
	}

	// Detect framework
	framework, err := p.detectFramework(projectPath, gh)
	if err != nil {
		return nil, err
	}

	metadata := map[string]interface{}{
		"hasDenoJson":   p.HasFile(files, "deno.json"),
		"hasDenoJsonc":  p.HasFile(files, "deno.jsonc"),
		"hasDenoLock":   p.HasFile(files, "deno.lock"),
		"hasImportMap":  p.HasFile(files, "import_map.json"),
		"hasDepsTs":     p.HasFile(files, "deps.ts"),
		"hasModTs":      p.HasFile(files, "mod.ts"),
		"hasTypeScript": p.HasAnyFile(files, []string{"*.ts"}),
	}

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	if p.HasFile(files, "deno.json") {
		evidenceFiles = append(evidenceFiles, "deno.json")
	}
	if p.HasFile(files, "deno.jsonc") {
		evidenceFiles = append(evidenceFiles, "deno.jsonc")
	}
	if p.HasFile(files, "deno.lock") {
		evidenceFiles = append(evidenceFiles, "deno.lock")
	}
	if p.HasFile(files, "import_map.json") {
		evidenceFiles = append(evidenceFiles, "import_map.json")
	}
	if p.HasFile(files, "deps.ts") {
		evidenceFiles = append(evidenceFiles, "deps.ts")
	}
	if p.HasFile(files, "mod.ts") {
		evidenceFiles = append(evidenceFiles, "mod.ts")
	}

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected Deno project based on: "
	var reasons []string
	if p.HasAnyFile(files, []string{"*.ts", "*.js"}) {
		reasons = append(reasons, "TypeScript/JavaScript source files")
	}
	if p.HasFile(files, "deno.json") || p.HasFile(files, "deno.jsonc") {
		reasons = append(reasons, "Deno configuration file")
	}
	if p.HasFile(files, "deno.lock") {
		reasons = append(reasons, "Deno lock file")
	}
	if p.HasFile(files, "deps.ts") || p.HasFile(files, "mod.ts") {
		reasons = append(reasons, "Deno module files")
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
		"deno",
		version,
		framework,
		"",
		"",
		metadata,
		evidence,
	), nil
}

// detectDenoVersion detects Deno version
func (p *DenoProvider) detectDenoVersion(projectPath string, gitHandler interface{}) (*types.VersionInfo, error) {
	// Read from deno.json
	denoJson, err := p.SafeReadJSON(projectPath, "deno.json", gitHandler)
	if err != nil {
		// If file doesn't exist, try deno.jsonc
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			// Continue trying deno.jsonc
		} else {
			return nil, err
		}
	}
	if denoJson != nil {
		if version, ok := denoJson["version"].(string); ok {
			return p.CreateVersionInfo(version, "deno.json"), nil
		}
	}

	// Read from deno.jsonc
	denoJsonc, err := p.SafeReadJSON(projectPath, "deno.jsonc", gitHandler)
	if err != nil {
		// If file doesn't exist, use default version
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return p.CreateVersionInfo("1.40", "default"), nil
		}
		return nil, err
	}
	if denoJsonc != nil {
		if version, ok := denoJsonc["version"].(string); ok {
			return p.CreateVersionInfo(version, "deno.jsonc"), nil
		}
	}

	// Default version
	return p.CreateVersionInfo("1.40", "default"), nil
}

// detectFramework detects framework
func (p *DenoProvider) detectFramework(projectPath string, gitHandler interface{}) (string, error) {
	// Check dependencies in deno.json
	denoJson, err := p.SafeReadJSON(projectPath, "deno.json", gitHandler)
	if err != nil {
		// If file doesn't exist, continue trying other methods
		if !strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return "", err
		}
	}
	if denoJson != nil {
		if imports, ok := denoJson["imports"].(map[string]interface{}); ok {
			frameworkMap := map[string]string{
				"fresh": "Fresh",
				"oak":   "Oak",
				"hono":  "Hono",
				"aleph": "Aleph.js",
				"ultra": "Ultra",
			}

			for framework, name := range frameworkMap {
				for _, importURL := range imports {
					if importURLStr, ok := importURL.(string); ok {
						if strings.Contains(importURLStr, framework) {
							return name, nil
						}
					}
				}
			}
		}
	}

	// Check dependencies in deno.jsonc
	denoJsonc, err := p.SafeReadJSON(projectPath, "deno.jsonc", gitHandler)
	if err != nil {
		// If file doesn't exist, continue trying other methods
		if !strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return "", err
		}
	}
	if denoJsonc != nil {
		if imports, ok := denoJsonc["imports"].(map[string]interface{}); ok {
			frameworkMap := map[string]string{
				"fresh": "Fresh",
				"oak":   "Oak",
				"hono":  "Hono",
				"aleph": "Aleph.js",
				"ultra": "Ultra",
			}

			for framework, name := range frameworkMap {
				for _, importURL := range imports {
					if importURLStr, ok := importURL.(string); ok {
						if strings.Contains(importURLStr, framework) {
							return name, nil
						}
					}
				}
			}
		}
	}

	// Check deps.ts file
	depsContent, err := p.SafeReadText(projectPath, "deps.ts", gitHandler)
	if err != nil {
		// If file doesn't exist, return empty string instead of error
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return "", nil
		}
		return "", err
	}
	if depsContent != "" {
		frameworkMap := map[string]string{
			"fresh": "Fresh",
			"oak":   "Oak",
			"hono":  "Hono",
			"aleph": "Aleph.js",
			"ultra": "Ultra",
		}

		for framework, name := range frameworkMap {
			if strings.Contains(depsContent, framework) {
				return name, nil
			}
		}
	}

	return "", nil
}

// GenerateCommands generates Deno project commands
func (p *DenoProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	// Determine main file
	mainFile := "main.ts"
	if p.HasFileInEvidence(result.Evidence.Files, "mod.ts") {
		mainFile = "mod.ts"
	} else if p.HasFileInEvidence(result.Evidence.Files, "index.ts") {
		mainFile = "index.ts"
	} else if p.HasFileInEvidence(result.Evidence.Files, "app.ts") {
		mainFile = "app.ts"
	}

	commands.Dev = []string{
		"deno run --allow-all " + mainFile,
	}

	commands.Start = []string{
		"deno run --allow-all " + mainFile,
	}

	return commands
}

// NeedsNativeCompilation checks if Deno project needs native compilation
func (p *DenoProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Deno projects usually don't need native compilation, they are runtime interpreted
	return false
}
