/**
 * DevBox Pack Execution Plan Generator - PHP Provider
 */

package providers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/labring/devbox-pack/pkg/types"
)

// PHPProvider PHP project detector
type PHPProvider struct {
	BaseProvider
}

// NewPHPProvider creates PHP Provider
func NewPHPProvider() *PHPProvider {
	return &PHPProvider{
		BaseProvider: BaseProvider{
			Name:     "php",
			Language: "php",
			Priority: 10,
		},
	}
}

// GetName gets Provider name
func (p *PHPProvider) GetName() string {
	return p.Name
}

// GetPriority gets Provider priority
func (p *PHPProvider) GetPriority() int {
	return p.Priority
}

// Detect detects PHP project
func (p *PHPProvider) Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error) {
	indicators := []types.ConfidenceIndicator{
		{Weight: 30, Satisfied: p.HasFile(files, "composer.json")},
		{Weight: 25, Satisfied: p.HasAnyFile(files, []string{"*.php"})},
		{Weight: 15, Satisfied: p.HasFile(files, "composer.lock")},
		{Weight: 20, Satisfied: p.HasAnyFile(files, []string{"index.php", "app.php"})},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"vendor/", "autoload.php"})},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{".php-version", "phpunit.xml"})},
		{Weight: 50, Satisfied: p.HasFile(files, "artisan")},                                          // High weight for Laravel artisan
		{Weight: 15, Satisfied: p.HasAnyFile(files, []string{"app/", "config/", "resources/views/"})}, // Laravel directory structure
		{Weight: 10, Satisfied: p.HasFile(files, "wp-config.php")},                                    // WordPress detection
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.2 // Lower detection threshold

	if !detected {
		return p.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Detect version
	version, err := p.detectPHPVersion(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	// Detect framework
	framework, err := p.detectFramework(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	metadata := map[string]interface{}{
		"hasComposerJson": p.HasFile(files, "composer.json"),
		"hasComposerLock": p.HasFile(files, "composer.lock"),
		"hasPHPSrc":       p.HasAnyFile(files, []string{"*.php"}),
		"hasIndex":        p.HasFile(files, "index.php"),
		"hasVendor":       p.HasFile(files, "vendor"),
		"framework":       framework,
	}

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	if p.HasFile(files, "composer.json") {
		evidenceFiles = append(evidenceFiles, "composer.json")
	}
	if p.HasFile(files, "composer.lock") {
		evidenceFiles = append(evidenceFiles, "composer.lock")
	}
	if p.HasFile(files, "index.php") {
		evidenceFiles = append(evidenceFiles, "index.php")
	}
	if p.HasFile(files, ".php-version") {
		evidenceFiles = append(evidenceFiles, ".php-version")
	}
	if p.HasFile(files, "phpunit.xml") {
		evidenceFiles = append(evidenceFiles, "phpunit.xml")
	}
	if p.HasFile(files, "artisan") {
		evidenceFiles = append(evidenceFiles, "artisan")
	}

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected PHP project based on: "
	var reasons []string
	if p.HasAnyFile(files, []string{"*.php"}) {
		reasons = append(reasons, "PHP source files")
	}
	if p.HasFile(files, "composer.json") {
		reasons = append(reasons, "Composer configuration (composer.json)")
	}
	if p.HasFile(files, "composer.lock") {
		reasons = append(reasons, "Composer lock file")
	}
	if p.HasFile(files, "index.php") {
		reasons = append(reasons, "PHP entry point (index.php)")
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
		"php",
		version,
		framework,
		"composer",
		"composer",
		metadata,
		evidence,
	), nil
}

// detectPHPVersion detects PHP version
func (p *PHPProvider) detectPHPVersion(projectPath string, gitHandler interface{}) (*types.VersionInfo, error) {
	// Read from composer.json
	composerJson, err := p.SafeReadJSON(projectPath, "composer.json", gitHandler)
	if err != nil {
		// If file doesn't exist, continue trying other methods
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			// Continue trying .php-version
		} else {
			return nil, err
		}
	}
	if composerJson != nil {
		if require, ok := composerJson["require"].(map[string]interface{}); ok {
			if phpVersion, ok := require["php"].(string); ok {
				return p.CreateVersionInfo(p.normalizePHPVersion(phpVersion), "composer.json require"), nil
			}
		}
	}

	// Read from .php-version
	version, err := p.ParseVersionFromText(
		projectPath,
		".php-version",
		gitHandler,
		regexp.MustCompile(`^(.+)$`),
	)
	if err != nil {
		// If file doesn't exist, use default version
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return p.CreateVersionInfo("8.2", "default"), nil
		}
		return nil, err
	}
	if version != "" {
		return p.CreateVersionInfo(p.NormalizeVersion(version), ".php-version"), nil
	}

	// Default version
	return p.CreateVersionInfo("8.2", "default"), nil
}

// detectFramework detects framework
func (p *PHPProvider) detectFramework(projectPath string, gitHandler interface{}) (string, error) {
	composerJson, err := p.SafeReadJSON(projectPath, "composer.json", gitHandler)
	if err != nil {
		// If file doesn't exist, return empty string instead of error
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return "", nil
		}
		return "", err
	}
	if composerJson == nil {
		return "", nil
	}

	frameworkMap := map[string]string{
		"laravel/framework":           "Laravel",
		"laravel/ui":                  "Laravel",
		"laravel/sanctum":             "Laravel",
		"laravel/passport":            "Laravel",
		"laravel/horizon":             "Laravel",
		"laravel/telescope":           "Laravel",
		"laravel/vapor":               "Laravel",
		"inertiajs/inertia-laravel":   "Laravel", // Laravel + Inertia.js
		"symfony/symfony":             "Symfony",
		"symfony/framework-bundle":    "Symfony",
		"codeigniter4/framework":      "CodeIgniter",
		"cakephp/cakephp":             "CakePHP",
		"yiisoft/yii2":                "Yii2",
		"zendframework/zendframework": "Zend Framework",
		"laminas/laminas-mvc":         "Laminas",
		"phalcon/cphalcon":            "Phalcon",
		"slim/slim":                   "Slim Framework",
		"doctrine/orm":                "Doctrine ORM",
		"twig/twig":                   "Twig",
	}

	return p.detectFrameworkFromComposerDependencies(composerJson, frameworkMap), nil
}

// detectFrameworkFromComposerDependencies detects framework from Composer dependencies
func (p *PHPProvider) detectFrameworkFromComposerDependencies(
	composerJson map[string]interface{},
	frameworkMap map[string]string,
) string {
	allDeps := make(map[string]interface{})

	// Merge require and require-dev dependencies (Composer format)
	if require, ok := composerJson["require"].(map[string]interface{}); ok {
		for k, v := range require {
			allDeps[k] = v
		}
	}
	if requireDev, ok := composerJson["require-dev"].(map[string]interface{}); ok {
		for k, v := range requireDev {
			allDeps[k] = v
		}
	}

	for depName, framework := range frameworkMap {
		if _, exists := allDeps[depName]; exists {
			return framework
		}
	}

	return ""
}

// normalizePHPVersion normalizes PHP version for base catalog lookup
func (p *PHPProvider) normalizePHPVersion(version string) string {
	// Remove prefix characters (like v, ^, ~, >=, etc.)
	re := regexp.MustCompile(`^[v^~>=<]+`)
	cleaned := re.ReplaceAllString(version, "")

	// Extract major.minor version for base catalog lookup
	versionRe := regexp.MustCompile(`^(\d+)(?:\.(\d+))?`)
	matches := versionRe.FindStringSubmatch(cleaned)
	if len(matches) > 1 {
		major := matches[1]
		minor := "0"

		if len(matches) > 2 && matches[2] != "" {
			minor = matches[2]
		}

		return fmt.Sprintf("%s.%s", major, minor)
	}

	return cleaned
}

// GenerateCommands generates commands for PHP project
func (p *PHPProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	// Check if composer.json exists
	hasComposer := p.HasFileInEvidence(result.Evidence.Files, "composer.json")
	hasIndex := p.HasFileInEvidence(result.Evidence.Files, "index.php")
	isLaravel := result.Framework == "Laravel"

	// Setup commands - install dependencies
	if hasComposer {
		commands.Setup = []string{"composer install"}
		if isLaravel {
			commands.Setup = append(commands.Setup, "php artisan key:generate", "php artisan config:cache")
		}
		commands.Build = []string{"composer install --no-dev --optimize-autoloader"}
		if isLaravel {
			commands.Build = append(commands.Build, "php artisan config:cache", "php artisan route:cache", "php artisan view:cache")
		}
	}

	// Development and Run commands
	if isLaravel {
		// Laravel-specific commands
		commands.Dev = []string{"php artisan serve"}
		commands.Run = []string{"php artisan serve"}
	} else if hasIndex {
		commands.Dev = []string{"php -S 0.0.0.0:8000 index.php"}
		commands.Run = []string{"php -S 0.0.0.0:8000 index.php"}
	} else {
		commands.Dev = []string{"php -S 0.0.0.0:8000"}
		commands.Run = []string{"php -S 0.0.0.0:8000"}
	}

	return commands
}

// GenerateEnvironment generates environment variables for PHP project
func (p *PHPProvider) GenerateEnvironment(result *types.DetectResult) map[string]string {
	env := make(map[string]string)

	// Set PHP specific environment variables
	env["PHP_ENV"] = "production"

	// Laravel-specific environment variables
	if result.Framework == "Laravel" {
		env["APP_ENV"] = "local"
		env["APP_DEBUG"] = "true"
		env["APP_KEY"] = "" // Will be generated by artisan key:generate
		env["LARAVEL_ENV"] = "local"
		env["PORT"] = "8000" // Laravel artisan serve default
	} else {
		// Set port for generic PHP web applications
		env["PORT"] = "8000"
	}

	// Add PHP version if available
	if result.Version != "" {
		env["PHP_VERSION"] = result.Version
	}

	return env
}

// NeedsNativeCompilation checks if PHP project needs native compilation
func (p *PHPProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// PHP projects usually don't need native compilation, unless there are special extensions
	// Most PHP projects are interpreted
	return false
}
