/**
 * DevBox Pack Execution Plan Generator - Ruby Provider
 */

package providers

import (
	"regexp"

	"github.com/labring/devbox-pack/pkg/types"
)

// RubyProvider Ruby project detector
type RubyProvider struct {
	BaseProvider
}

// NewRubyProvider creates Ruby Provider
func NewRubyProvider() *RubyProvider {
	return &RubyProvider{
		BaseProvider: BaseProvider{
			Name:     "ruby",
			Language: "ruby",
			Priority: 65,
		},
	}
}

// GetName gets Provider name
func (p *RubyProvider) GetName() string {
	return p.Name
}

// GetPriority gets Provider priority
func (p *RubyProvider) GetPriority() int {
	return p.Priority
}

// Detect detects Ruby project
func (p *RubyProvider) Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error) {
	indicators := []types.ConfidenceIndicator{
		{Weight: 30, Satisfied: p.HasFile(files, "Gemfile")},
		{Weight: 20, Satisfied: p.HasFile(files, "Gemfile.lock")},
		{Weight: 20, Satisfied: p.HasAnyFile(files, []string{"*.rb"})},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{".ruby-version", ".rvmrc"})},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"config.ru", "Rakefile"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"app/", "lib/"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"spec/", "test/"})},
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.3

	if !detected {
		return p.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Detect version
	version, err := p.detectRubyVersion(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	// Detect framework
	framework, err := p.detectFramework(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	metadata := map[string]interface{}{
		"hasGemfile":     p.HasFile(files, "Gemfile"),
		"hasGemfileLock": p.HasFile(files, "Gemfile.lock"),
		"hasRakefile":    p.HasFile(files, "Rakefile"),
		"hasConfigRu":    p.HasFile(files, "config.ru"),
		"hasRubyVersion": p.HasFile(files, ".ruby-version"),
		"framework":      framework,
	}

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	if p.HasFile(files, "Gemfile") {
		evidenceFiles = append(evidenceFiles, "Gemfile")
	}
	if p.HasFile(files, "Gemfile.lock") {
		evidenceFiles = append(evidenceFiles, "Gemfile.lock")
	}
	if p.HasFile(files, "Rakefile") {
		evidenceFiles = append(evidenceFiles, "Rakefile")
	}
	if p.HasFile(files, "config.ru") {
		evidenceFiles = append(evidenceFiles, "config.ru")
	}
	if p.HasFile(files, ".ruby-version") {
		evidenceFiles = append(evidenceFiles, ".ruby-version")
	}
	if p.HasFile(files, ".rvmrc") {
		evidenceFiles = append(evidenceFiles, ".rvmrc")
	}
	if p.HasFile(files, ".rbenv-version") {
		evidenceFiles = append(evidenceFiles, ".rbenv-version")
	}

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected Ruby project based on: "
	var reasons []string
	if p.HasAnyFile(files, []string{"*.rb"}) {
		reasons = append(reasons, "Ruby source files")
	}
	if p.HasFile(files, "Gemfile") {
		reasons = append(reasons, "Gemfile")
	}
	if p.HasFile(files, "Rakefile") {
		reasons = append(reasons, "Rakefile")
	}
	if p.HasFile(files, "config.ru") {
		reasons = append(reasons, "Rack configuration (config.ru)")
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
		"ruby",
		version,
		framework,
		"bundler",
		"bundler",
		metadata,
		evidence,
	), nil
}

// detectRubyVersion detects Ruby version
func (p *RubyProvider) detectRubyVersion(projectPath string, gitHandler interface{}) (*types.VersionInfo, error) {
	// Read from .ruby-version
	version, err := p.ParseVersionFromText(
		projectPath,
		".ruby-version",
		gitHandler,
		regexp.MustCompile(`^(.+?)(?:\s|$)`),
	)
	if err == nil && version != "" {
		return p.CreateVersionInfo(p.NormalizeVersion(version), ".ruby-version"), nil
	}

	// Read from .rvmrc
	version, err = p.ParseVersionFromText(
		projectPath,
		".rvmrc",
		gitHandler,
		regexp.MustCompile(`ruby-(.+)`),
	)
	if err == nil && version != "" {
		return p.CreateVersionInfo(p.NormalizeVersion(version), ".rvmrc"), nil
	}

	// Read from Gemfile
	gemfileContent, err := p.SafeReadText(projectPath, "Gemfile", gitHandler)
	if err == nil && gemfileContent != "" {
		re := regexp.MustCompile(`ruby\s+['"]([^'"]+)['"]`)
		matches := re.FindStringSubmatch(gemfileContent)
		if len(matches) > 1 {
			return p.CreateVersionInfo(p.NormalizeVersion(matches[1]), "Gemfile"), nil
		}
	}

	// Default version
	return p.CreateVersionInfo("3.2", "default"), nil
}

// detectFramework detects framework
func (p *RubyProvider) detectFramework(projectPath string, gitHandler interface{}) (string, error) {
	gemfileContent, err := p.SafeReadText(projectPath, "Gemfile", gitHandler)
	if err != nil {
		return "", err
	}
	if gemfileContent == "" {
		return "", nil
	}

	frameworkMap := map[string]string{
		"rails":     "Ruby on Rails",
		"sinatra":   "Sinatra",
		"grape":     "Grape",
		"hanami":    "Hanami",
		"roda":      "Roda",
		"cuba":      "Cuba",
		"padrino":   "Padrino",
		"jekyll":    "Jekyll",
		"middleman": "Middleman",
	}

	for gem, framework := range frameworkMap {
		pattern := `gem\s+['"]` + gem + `['"]`
		re := regexp.MustCompile(`(?i)` + pattern)
		if re.MatchString(gemfileContent) {
			return framework, nil
		}
	}

	return "", nil
}

// GenerateCommands generates commands for Ruby project
func (p *RubyProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	// Check file existence
	hasGemfile := p.HasFileInEvidence(result.Evidence.Files, "Gemfile")
	hasConfigRu := p.HasFileInEvidence(result.Evidence.Files, "config.ru")
	hasAppRb := p.HasFileInEvidence(result.Evidence.Files, "app.rb")

	if hasGemfile {
		commands.Dev = []string{
			"bundle install",
		}
		commands.Build = []string{
			"bundle install --without development test",
		}
	}

	// Determine start command
	if hasConfigRu {
		// Rack application
		commands.Dev = append(commands.Dev, "bundle exec rackup -o 0.0.0.0 -p 4567")
		commands.Start = []string{"bundle exec rackup -o 0.0.0.0 -p 4567"}
	} else if hasAppRb {
		commands.Dev = append(commands.Dev, "ruby app.rb")
		commands.Start = []string{"ruby app.rb"}
	} else {
		commands.Dev = append(commands.Dev, "ruby -run -e httpd . -p 4567")
		commands.Start = []string{"ruby -run -e httpd . -p 4567"}
	}

	return commands
}

// NeedsNativeCompilation checks if Ruby project needs native compilation
func (p *RubyProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Ruby projects usually don't need native compilation, unless they have gems with C extensions
	// Most Ruby projects are interpreted
	return false
}
