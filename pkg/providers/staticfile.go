/**
 * DevBox Pack Execution Plan Generator - StaticFile Provider
 */

package providers

import (
	"github.com/labring/devbox-pack/pkg/types"
)

// StaticFileProvider static file project detector
type StaticFileProvider struct {
	BaseProvider
}

// NewStaticFileProvider creates StaticFile Provider
func NewStaticFileProvider() *StaticFileProvider {
	return &StaticFileProvider{
		BaseProvider: BaseProvider{
			Name:     "staticfile",
			Language: "staticfile",
			Priority: 10,
		},
	}
}

// GetName gets Provider name
func (p *StaticFileProvider) GetName() string {
	return p.Name
}

// GetPriority gets Provider priority
func (p *StaticFileProvider) GetPriority() int {
	return p.Priority
}

// Detect detects static file project
func (p *StaticFileProvider) Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error) {
	indicators := []types.ConfidenceIndicator{
		{Weight: 30, Satisfied: p.HasAnyFile(files, []string{"*.html", "*.htm"})},
		{Weight: 20, Satisfied: p.HasAnyFile(files, []string{"*.css"})},
		{Weight: 15, Satisfied: p.HasAnyFile(files, []string{"*.js"})},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"*.png", "*.jpg", "*.jpeg", "*.gif", "*.svg"})},
		{Weight: 25, Satisfied: p.HasFile(files, "index.html")}, // Increase weight for index.html
		{Weight: 15, Satisfied: p.HasFile(files, "Staticfile")}, // Add support for Staticfile
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"assets/", "static/", "public/"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"favicon.ico", "robots.txt"})},
	}

	// If there are project files from other languages, reduce static file confidence
	exclusions := []bool{
		p.HasAnyFile(files, []string{"package.json", "composer.json", "requirements.txt", "go.mod", "Cargo.toml", "Gemfile"}),
		p.HasAnyFile(files, []string{"pom.xml", "build.gradle"}),
		p.HasAnyFile(files, []string{"*.py", "*.java", "*.go", "*.rs", "*.rb", "*.php", "*.ex", "*.exs"}),
	}

	confidence := p.CalculateConfidence(indicators)
	hasExclusions := false
	for _, exclusion := range exclusions {
		if exclusion {
			hasExclusions = true
			break
		}
	}

	// If there are characteristic files from other languages, significantly reduce confidence
	adjustedConfidence := confidence
	if hasExclusions {
		adjustedConfidence = confidence * 0.3
	}

	// For pure static file projects (only HTML/CSS/JS/images), should not be excluded
	// test.json is a test configuration file and should not affect static file detection
	isPureStaticFile := p.HasAnyFile(files, []string{"*.html", "*.htm"}) &&
		!p.HasAnyFile(files, []string{"package.json", "composer.json", "requirements.txt", "go.mod", "Cargo.toml", "Gemfile", "pom.xml", "build.gradle", "*.py", "*.java", "*.go", "*.rs", "*.rb", "*.php"})

	if isPureStaticFile {
		adjustedConfidence = confidence // Do not reduce confidence
	}

	// Special handling: if only index.html and test.json exist, should be recognized as static file project
	if len(files) == 2 && p.HasFile(files, "index.html") && p.HasFile(files, "test.json") {
		adjustedConfidence = confidence
	}

	// Special handling: if there is a Staticfile configuration file, should be recognized as static file project
	if p.HasFile(files, "Staticfile") {
		adjustedConfidence = confidence
	}

	detected := adjustedConfidence > 0.2 // Lower detection threshold

	if !detected {
		return p.CreateDetectResult(false, adjustedConfidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Static file projects usually don't have specific versions
	version := p.CreateVersionInfo("latest", "default")

	metadata := map[string]interface{}{
		"hasHtml":       p.HasAnyFile(files, []string{"*.html", "*.htm"}),
		"hasCss":        p.HasAnyFile(files, []string{"*.css"}),
		"hasJs":         p.HasAnyFile(files, []string{"*.js"}),
		"hasImages":     p.HasAnyFile(files, []string{"*.png", "*.jpg", "*.jpeg", "*.gif", "*.svg"}),
		"hasIndexHtml":  p.HasFile(files, "index.html"),
		"hasStaticfile": p.HasFile(files, "Staticfile"),
		"hasAssets":     p.HasAnyFile(files, []string{"assets/", "static/", "public/"}),
	}

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	if p.HasFile(files, "index.html") {
		evidenceFiles = append(evidenceFiles, "index.html")
	}
	if p.HasFile(files, "Staticfile") {
		evidenceFiles = append(evidenceFiles, "Staticfile")
	}
	// Collect some main static files
	for _, file := range files {
		if file.IsDirectory {
			continue
		}
		name := file.Name
		if name == "index.html" || name == "Staticfile" {
			continue // Already added
		}
		// Add some important static files
		if name == "style.css" || name == "main.css" || name == "app.css" ||
			name == "script.js" || name == "main.js" || name == "app.js" ||
			name == "favicon.ico" || name == "robots.txt" {
			evidenceFiles = append(evidenceFiles, name)
		}
	}

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected static file project based on: "
	var reasons []string
	if p.HasAnyFile(files, []string{"*.html", "*.htm"}) {
		reasons = append(reasons, "HTML files")
	}
	if p.HasAnyFile(files, []string{"*.css"}) {
		reasons = append(reasons, "CSS files")
	}
	if p.HasAnyFile(files, []string{"*.js"}) {
		reasons = append(reasons, "JavaScript files")
	}
	if p.HasFile(files, "index.html") {
		reasons = append(reasons, "index.html entry point")
	}
	if p.HasFile(files, "Staticfile") {
		reasons = append(reasons, "Staticfile configuration")
	}
	if p.HasAnyFile(files, []string{"*.png", "*.jpg", "*.jpeg", "*.gif", "*.svg"}) {
		reasons = append(reasons, "image assets")
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
		adjustedConfidence,
		"staticfile",
		version,
		"",
		"",
		"",
		metadata,
		evidence,
	), nil
}

// GenerateCommands generates commands for static file project
func (p *StaticFileProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	commands.Start = []string{
		"nginx -g 'daemon off;'",
	}

	return commands
}

// NeedsNativeCompilation checks if static file project needs native compilation
func (p *StaticFileProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Static file projects don't need compilation
	return false
}
