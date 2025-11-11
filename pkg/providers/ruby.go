/**
 * DevBox Pack Execution Plan Generator - Ruby Provider
 */

package providers

import (
	"regexp"
	"strings"

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
			Priority: 50,
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
	// Check for Rails-specific files first
	isRailsProject := p.HasAnyFile(files, []string{
		"config/application.rb",
		"config/routes.rb",
		"app/models/",
		"app/controllers/",
		"app/views/",
		"bin/rails",
		"config/environments/development.rb",
	})

	// Adjust indicators for Rails detection
	indicators := []types.ConfidenceIndicator{
		{Weight: 30, Satisfied: p.HasFile(files, "Gemfile")},
		{Weight: 20, Satisfied: p.HasFile(files, "Gemfile.lock")},
		{Weight: 25, Satisfied: p.HasAnyFile(files, []string{"*.rb"})}, // Higher weight for .rb files
		{Weight: 15, Satisfied: isRailsProject}, // Higher weight for Rails-specific files
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{".ruby-version", ".rvmrc"})},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"config.ru", "Rakefile"})},
		{Weight: 10, Satisfied: p.HasFile(files, "config/database.yml")},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"app/", "lib/"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"spec/", "test/"})},
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.2 // Lower threshold for projects with only .rb files

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

	// Detect Rails-specific features
	var railsFeatures []string
	if framework == "Rails" {
		railsFeatures = p.detectRailsFeatures(projectPath, files, gitHandler)
	}

	// Detect asset pipeline
	assetPipeline := p.detectAssetPipeline(projectPath, gitHandler)

	metadata := map[string]interface{}{
		"hasGemfile":       p.HasFile(files, "Gemfile"),
		"hasGemfileLock":   p.HasFile(files, "Gemfile.lock"),
		"hasRakefile":      p.HasFile(files, "Rakefile"),
		"hasConfigRu":      p.HasFile(files, "config.ru"),
		"hasRubyVersion":   p.HasFile(files, ".ruby-version"),
		"hasDatabaseYml":   p.HasFile(files, "config/database.yml"),
		"hasRoutesRb":      p.HasFile(files, "config/routes.rb"),
		"hasApplicationRb": p.HasFile(files, "config/application.rb"),
		"isRailsProject":   isRailsProject,
		"framework":        framework,
		"railsFeatures":    railsFeatures,
		"assetPipeline":    assetPipeline,
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

	// Add Rails-specific files
	if isRailsProject {
		if p.HasFile(files, "config/application.rb") {
			evidenceFiles = append(evidenceFiles, "config/application.rb")
		}
		if p.HasFile(files, "config/routes.rb") {
			evidenceFiles = append(evidenceFiles, "config/routes.rb")
		}
		if p.HasFile(files, "config/database.yml") {
			evidenceFiles = append(evidenceFiles, "config/database.yml")
		}
		if p.HasFile(files, "bin/rails") {
			evidenceFiles = append(evidenceFiles, "bin/rails")
		}
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
	if framework == "Rails" {
		reasons = append(reasons, "Ruby on Rails framework")
		if len(railsFeatures) > 0 {
			reasons = append(reasons, "features: "+strings.Join(railsFeatures, ", "))
		}
		if assetPipeline != "" {
			reasons = append(reasons, "asset pipeline: "+assetPipeline)
		}
	}
	if p.HasFile(files, "Rakefile") {
		reasons = append(reasons, "Rakefile")
	}
	if p.HasFile(files, "config.ru") {
		reasons = append(reasons, "Rack configuration (config.ru)")
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
		regexp.MustCompile(`rvm use ([\d\.]+)`),
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

// detectRailsFeatures detects Rails-specific features
func (p *RubyProvider) detectRailsFeatures(projectPath string, files []types.FileInfo, gitHandler interface{}) []string {
	var features []string

	// Check for ActiveRecord
	if p.HasFile(files, "app/models/") || p.HasFile(files, "config/database.yml") {
		features = append(features, "ActiveRecord")
	}

	// Check for ActionCable
	if p.HasFile(files, "app/channels/") {
		features = append(features, "ActionCable")
	}

	// Check for ActionMailbox
	if p.HasFile(files, "app/mailboxes/") {
		features = append(features, "ActionMailbox")
	}

	// Check for ActionText
	if p.HasFile(files, "app/views/action_text/") {
		features = append(features, "ActionText")
	}

	// Check for ActiveStorage
	if p.HasFile(files, "app/models/active_storage/") || p.HasFile(files, "storage/") {
		features = append(features, "ActiveStorage")
	}

	// Check for ActiveJob
	if p.HasFile(files, "app/jobs/") {
		features = append(features, "ActiveJob")
	}

	// Check for API mode
	if p.HasFile(files, "config/routes/api.rb") {
		features = append(features, "API")
	}

	return features
}

// detectAssetPipeline detects which asset pipeline is used
func (p *RubyProvider) detectAssetPipeline(projectPath string, gitHandler interface{}) string {
	// Check for Sprockets
	gemfileContent, err := p.SafeReadText(projectPath, "Gemfile", gitHandler)
	if err == nil && gemfileContent != "" {
		if regexp.MustCompile(`(?i)gem\s+['"]sprockets['"]`).MatchString(gemfileContent) {
			return "Sprockets"
		}
	}

	// Check for Propshaft
	if err == nil && gemfileContent != "" {
		if regexp.MustCompile(`(?i)gem\s+['"]propshaft['"]`).MatchString(gemfileContent) {
			return "Propshaft"
		}
	}

	// Check for jsbundling-rails (modern approach)
	if err == nil && gemfileContent != "" {
		if regexp.MustCompile(`(?i)gem\s+['"]jsbundling-rails['"]`).MatchString(gemfileContent) {
			return "JS Bundling"
		}
	}

	// Check for cssbundling-rails (modern approach)
	if err == nil && gemfileContent != "" {
		if regexp.MustCompile(`(?i)gem\s+['"]cssbundling-rails['"]`).MatchString(gemfileContent) {
			return "CSS Bundling"
		}
	}

	return ""
}

// detectFramework detects framework
func (p *RubyProvider) detectFramework(projectPath string, gitHandler interface{}) (string, error) {
	gemfileContent, err := p.SafeReadText(projectPath, "Gemfile", gitHandler)
	if err != nil {
		// If file doesn't exist, return empty string instead of error
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return "", nil
		}
		return "", err
	}
	if gemfileContent == "" {
		return "", nil
	}

	frameworkMap := map[string]string{
		"rails":     "Rails",
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
	hasRailsApp := p.HasFileInEvidence(result.Evidence.Files, "config/application.rb")

	// Setup commands - install dependencies
	if hasGemfile {
		commands.Setup = []string{"bundle install"}
	}

	// Development and Run commands
	if result.Framework == "Rails" || hasRailsApp {
		// Rails-specific commands
		commands.Setup = append(commands.Setup, "rails db:prepare")
		commands.Dev = []string{"bundle exec rails server -b 0.0.0.0 -p 3000"}
		commands.Run = []string{"bundle exec rails server -b 0.0.0.0 -e production -p 3000"}
		commands.Build = []string{"bundle exec rails assets:precompile"}
	} else if hasConfigRu {
		// Rack application
		commands.Dev = []string{"bundle exec rackup -o 0.0.0.0 -p 4567"}
		commands.Run = []string{"bundle exec rackup -o 0.0.0.0 -p 4567"}
	} else if hasAppRb {
		commands.Dev = []string{"ruby app.rb"}
		commands.Run = []string{"ruby app.rb"}
	} else {
		commands.Dev = []string{"ruby -run -e httpd . -p 4567"}
		commands.Run = []string{"ruby -run -e httpd . -p 4567"}
	}

	return commands
}

// GenerateEnvironment generates environment variables for Ruby project
func (p *RubyProvider) GenerateEnvironment(result *types.DetectResult) map[string]string {
	env := make(map[string]string)

	// Set Ruby specific environment variables
	env["RACK_ENV"] = "production"
	env["RAILS_ENV"] = "production"

	// Rails-specific environment variables
	if result.Framework == "Rails" {
		// Set default Rails port
		env["PORT"] = "3000"

		// Add Rails-specific environment variables
		env["RAILS_MASTER_KEY"] = ""
		env["RAILS_LOG_TO_STDOUT"] = "true"
		env["RAILS_SERVE_STATIC_FILES"] = "true"
		env["BUNDLE_WITHOUT"] = "development:test"

		// Add asset pipeline environment variables
		if result.Metadata != nil {
			if assetPipeline, ok := result.Metadata["assetPipeline"].(string); ok {
				switch assetPipeline {
				case "Sprockets":
					env["RAILS_ENV"] = "production"
					env["RAILS_GROUPS"] = "assets"
				case "Propshaft":
					env["RAILS_ENV"] = "production"
				case "JS Bundling":
					env["RAILS_ENV"] = "production"
				case "CSS Bundling":
					env["RAILS_ENV"] = "production"
				}
			}

			// Add Rails features information
			if railsFeatures, ok := result.Metadata["railsFeatures"].([]string); ok {
				for _, feature := range railsFeatures {
					switch feature {
					case "ActiveRecord":
						env["DATABASE_URL"] = "sqlite3:///db/production.sqlite3"
					case "ActionCable":
						env["ACTION_CABLE_MOUNT_PATH"] = "/cable"
					case "API":
						env["RAILS_API_ONLY"] = "true"
					}
				}
			}
		}
	} else {
		// Set port for non-Rails Ruby applications
		env["PORT"] = "4567"
	}

	// Add Ruby version if available
	if result.Version != "" {
		env["RUBY_VERSION"] = result.Version
	}

	return env
}

// NeedsNativeCompilation checks if Ruby project needs native compilation
func (p *RubyProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Check metadata for native gems flag
	if result.Metadata != nil {
		if hasNativeGems, ok := result.Metadata["hasNativeGems"].(bool); ok {
			return hasNativeGems
		}
	}

	// Ruby projects usually don't need native compilation, unless they have gems with C extensions
	// Most Ruby projects are interpreted
	return false
}
