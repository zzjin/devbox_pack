/**
 * DevBox Pack Execution Plan Generator - Python Provider
 */

package providers

import (
	"regexp"
	"strings"

	"github.com/labring/devbox-pack/pkg/types"
)

// PythonProvider Python project detector
type PythonProvider struct {
	BaseProvider
}

// NewPythonProvider creates Python Provider
func NewPythonProvider() *PythonProvider {
	return &PythonProvider{
		BaseProvider: BaseProvider{
			Name:     "python",
			Language: "python",
			Priority: 80,
		},
	}
}

// GetName gets Provider name
func (p *PythonProvider) GetName() string {
	return p.Name
}

// GetPriority gets Provider priority
func (p *PythonProvider) GetPriority() int {
	return p.Priority
}

// Detect detects Python project
func (p *PythonProvider) Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error) {
	indicators := []types.ConfidenceIndicator{
		{Weight: 30, Satisfied: p.HasAnyFile(files, []string{"requirements.txt", "pyproject.toml", "setup.py", "Pipfile"})},
		{Weight: 25, Satisfied: p.HasAnyFile(files, []string{"*.py"})},
		{Weight: 15, Satisfied: p.HasFile(files, "poetry.lock")},
		{Weight: 15, Satisfied: p.HasFile(files, "Pipfile.lock")},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{".python-version", "runtime.txt"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"manage.py", "app.py", "main.py"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"__pycache__", "*.pyc"})},
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.3

	if !detected {
		return p.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Detect version
	version, err := p.detectPythonVersion(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	// Detect framework
	framework, err := p.detectFramework(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	// Detect package manager
	packageManager := p.DetectPackageManager(files, map[string]string{
		"poetry.lock":      "poetry",
		"Pipfile.lock":     "pipenv",
		"requirements.txt": "pip",
		"pyproject.toml":   "pip",
	})
	if packageManager == "" {
		packageManager = "pip"
	}

	metadata := map[string]interface{}{
		"hasRequirements":  p.HasFile(files, "requirements.txt"),
		"hasPyprojectToml": p.HasFile(files, "pyproject.toml"),
		"hasSetupPy":       p.HasFile(files, "setup.py"),
		"hasPipfile":       p.HasFile(files, "Pipfile"),
		"hasPoetryLock":    p.HasFile(files, "poetry.lock"),
		"packageManager":   packageManager,
		"framework":        framework,
	}

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	if p.HasFile(files, "requirements.txt") {
		evidenceFiles = append(evidenceFiles, "requirements.txt")
	}
	if p.HasFile(files, "pyproject.toml") {
		evidenceFiles = append(evidenceFiles, "pyproject.toml")
	}
	if p.HasFile(files, "setup.py") {
		evidenceFiles = append(evidenceFiles, "setup.py")
	}
	if p.HasFile(files, "Pipfile") {
		evidenceFiles = append(evidenceFiles, "Pipfile")
	}
	if p.HasFile(files, "Pipfile.lock") {
		evidenceFiles = append(evidenceFiles, "Pipfile.lock")
	}
	if p.HasFile(files, "poetry.lock") {
		evidenceFiles = append(evidenceFiles, "poetry.lock")
	}
	if p.HasFile(files, ".python-version") {
		evidenceFiles = append(evidenceFiles, ".python-version")
	}
	if p.HasFile(files, "runtime.txt") {
		evidenceFiles = append(evidenceFiles, "runtime.txt")
	}

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected Python project based on: "
	var reasons []string
	if p.HasAnyFile(files, []string{"*.py"}) {
		reasons = append(reasons, "Python source files")
	}
	if p.HasFile(files, "requirements.txt") {
		reasons = append(reasons, "requirements.txt")
	}
	if p.HasFile(files, "pyproject.toml") {
		reasons = append(reasons, "pyproject.toml")
	}
	if p.HasFile(files, "setup.py") {
		reasons = append(reasons, "setup.py")
	}
	if packageManager != "pip" {
		reasons = append(reasons, packageManager+" configuration")
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
		"python",
		version,
		framework,
		packageManager,
		packageManager,
		metadata,
		evidence,
	), nil
}

// GenerateCommands generates commands for Python project
func (p *PythonProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	// Check different Python project types
	hasRequirements := p.HasFileInEvidence(result.Evidence.Files, "requirements.txt")
	hasPipfile := p.HasFileInEvidence(result.Evidence.Files, "Pipfile")
	hasPyproject := p.HasFileInEvidence(result.Evidence.Files, "pyproject.toml")
	hasApp := p.HasFileInEvidence(result.Evidence.Files, "app.py")
	hasMain := p.HasFileInEvidence(result.Evidence.Files, "main.py")
	hasManage := p.HasFileInEvidence(result.Evidence.Files, "manage.py")

	if hasRequirements {
		// Standard Python project
		commands.Dev = []string{
			"pip install -r requirements.txt",
		}
		commands.Build = []string{
			"pip install -r requirements.txt",
		}
	} else if hasPipfile {
		// Pipenv project
		commands.Dev = []string{
			"pipenv install --dev",
		}
		commands.Build = []string{
			"pipenv install",
		}
	} else if hasPyproject {
		// Poetry or modern Python project
		commands.Dev = []string{
			"pip install -e .",
		}
		commands.Build = []string{
			"pip install .",
		}
	}

	// Determine start command
	if hasApp {
		commands.Dev = append(commands.Dev, "python app.py")
		commands.Start = []string{"python app.py"}
	} else if hasMain {
		commands.Dev = append(commands.Dev, "python main.py")
		commands.Start = []string{"python main.py"}
	} else if hasManage {
		// Django project
		commands.Dev = append(commands.Dev, "python manage.py runserver 0.0.0.0:8000")
		commands.Start = []string{"python manage.py runserver 0.0.0.0:8000"}
	} else {
		// Generic Python startup
		commands.Dev = append(commands.Dev, "python -m http.server 8000")
		commands.Start = []string{"python -m http.server 8000"}
	}

	return commands
}

// NeedsNativeCompilation checks if Python project needs native compilation
func (p *PythonProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Check common packages that need compilation
	compilationIndicators := []string{
		"numpy", "scipy", "pandas", "pillow", "lxml",
		"psycopg2", "mysqlclient", "cryptography",
		"cffi", "cython", "pycrypto", "gevent",
	}

	evidenceStr := strings.Join(result.Evidence.Files, " ")
	for _, indicator := range compilationIndicators {
		if strings.Contains(evidenceStr, indicator) {
			return true
		}
	}
	return false
}

// Helper methods

// detectPythonVersion detects Python version
func (p *PythonProvider) detectPythonVersion(projectPath string, gitHandler interface{}) (*types.VersionInfo, error) {
	// Read from .python-version
	version, err := p.ParseVersionFromText(
		projectPath,
		".python-version",
		gitHandler,
		regexp.MustCompile(`^(.+?)(?:\s|$)`),
	)
	if err == nil && version != "" {
		return p.CreateVersionInfo(p.NormalizeVersion(version), ".python-version"), nil
	}

	// Read from runtime.txt (Heroku)
	version, err = p.ParseVersionFromText(
		projectPath,
		"runtime.txt",
		gitHandler,
		regexp.MustCompile(`python-(.+)$`),
	)
	if err == nil && version != "" {
		return p.CreateVersionInfo(p.NormalizeVersion(version), "runtime.txt"), nil
	}

	// Read from pyproject.toml
	pyprojectContent, err := p.SafeReadText(projectPath, "pyproject.toml", gitHandler)
	if err == nil && pyprojectContent != "" {
		re := regexp.MustCompile(`python\s*=\s*["']([^"']+)["']`)
		matches := re.FindStringSubmatch(pyprojectContent)
		if len(matches) > 1 {
			return p.CreateVersionInfo(p.NormalizeVersion(matches[1]), "pyproject.toml"), nil
		}
	}

	// Read from Pipfile
	pipfileContent, err := p.SafeReadText(projectPath, "Pipfile", gitHandler)
	if err == nil && pipfileContent != "" {
		re := regexp.MustCompile(`python_version\s*=\s*["']([^"']+)["']`)
		matches := re.FindStringSubmatch(pipfileContent)
		if len(matches) > 1 {
			return p.CreateVersionInfo(p.NormalizeVersion(matches[1]), "Pipfile"), nil
		}
	}

	// Default version
	return p.CreateVersionInfo("3.11", "default"), nil
}

// detectFramework detects framework
func (p *PythonProvider) detectFramework(projectPath string, gitHandler interface{}) (string, error) {
	frameworkMap := map[string]string{
		"django":    "Django",
		"flask":     "Flask",
		"fastapi":   "FastAPI",
		"tornado":   "Tornado",
		"pyramid":   "Pyramid",
		"bottle":    "Bottle",
		"sanic":     "Sanic",
		"quart":     "Quart",
		"starlette": "Starlette",
		"streamlit": "Streamlit",
		"dash":      "Dash",
		"jupyter":   "Jupyter",
	}

	// Check requirements.txt
	requirements, err := p.SafeReadText(projectPath, "requirements.txt", gitHandler)
	if err == nil && requirements != "" {
		requirementsLower := strings.ToLower(requirements)
		for pkg, framework := range frameworkMap {
			if strings.Contains(requirementsLower, pkg) {
				return framework, nil
			}
		}
	}

	// Check pyproject.toml
	pyprojectToml, err := p.SafeReadText(projectPath, "pyproject.toml", gitHandler)
	if err == nil && pyprojectToml != "" {
		pyprojectLower := strings.ToLower(pyprojectToml)
		for pkg, framework := range frameworkMap {
			if strings.Contains(pyprojectLower, pkg) {
				return framework, nil
			}
		}
	}

	// Check Pipfile
	pipfile, err := p.SafeReadText(projectPath, "Pipfile", gitHandler)
	if err == nil && pipfile != "" {
		pipfileLower := strings.ToLower(pipfile)
		for pkg, framework := range frameworkMap {
			if strings.Contains(pipfileLower, pkg) {
				return framework, nil
			}
		}
	}

	return "", nil
}
