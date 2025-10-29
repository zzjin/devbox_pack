// Package providers contains the language and framework detection providers
// for the DevBox Pack execution plan generator.
package providers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

// BaseProvider Provider base class
// All language detectors should inherit from this class
type BaseProvider struct {
	Name     string
	Language string
	Priority int
}

// GetName gets Provider name
func (bp *BaseProvider) GetName() string {
	return bp.Name
}

// GetLanguage gets Provider language
func (bp *BaseProvider) GetLanguage() string {
	return bp.Language
}

// GetPriority gets Provider priority
func (bp *BaseProvider) GetPriority() int {
	return bp.Priority
}

// HasFile checks if file exists
func (bp *BaseProvider) HasFile(files []types.FileInfo, fileName string) bool {
	if strings.Contains(fileName, "*") {
		// Improve wildcard matching logic to ensure only file extensions are matched
		pattern := fileName
		if strings.HasPrefix(pattern, "*.") {
			// For *.ext format, ensure matching the entire filename
			ext := pattern[2:] // Remove "*."
			pattern = `.*\.` + regexp.QuoteMeta(ext) + `$`
		} else {
			// For other wildcards, use original logic
			pattern = strings.ReplaceAll(pattern, "*", ".*")
		}

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return false
		}
		for _, file := range files {
			if !file.IsDirectory && regex.MatchString(file.Path) {
				return true
			}
		}
		return false
	}

	for _, file := range files {
		if !file.IsDirectory && file.Path == fileName {
			return true
		}
	}
	return false
}

// HasAnyFile checks if any of multiple files exist (returns true if any exists)
func (bp *BaseProvider) HasAnyFile(files []types.FileInfo, fileNames []string) bool {
	for _, fileName := range fileNames {
		if bp.HasFile(files, fileName) {
			return true
		}
	}
	return false
}

// HasAllFiles checks if all files exist
func (bp *BaseProvider) HasAllFiles(files []types.FileInfo, fileNames []string) bool {
	for _, fileName := range fileNames {
		if !bp.HasFile(files, fileName) {
			return false
		}
	}
	return true
}

// GetMatchingFiles gets list of matching files
func (bp *BaseProvider) GetMatchingFiles(files []types.FileInfo, pattern string) []types.FileInfo {
	var matchingFiles []types.FileInfo

	if strings.Contains(pattern, "*") {
		regexPattern := strings.ReplaceAll(pattern, "*", ".*")
		regex, err := regexp.Compile(regexPattern)
		if err != nil {
			return matchingFiles
		}
		for _, file := range files {
			if !file.IsDirectory && regex.MatchString(file.Path) {
				matchingFiles = append(matchingFiles, file)
			}
		}
	} else {
		for _, file := range files {
			if !file.IsDirectory && file.Path == pattern {
				matchingFiles = append(matchingFiles, file)
			}
		}
	}

	return matchingFiles
}

// ParseVersionFromJSON parses version information from JSON file
func (bp *BaseProvider) ParseVersionFromJSON(
	projectPath string,
	fileName string,
	gitHandler interface{},
	versionField string,
) (string, error) {
	gh := gitHandler.(*git.GitHandler)
	if versionField == "" {
		versionField = "version"
	}

	var content map[string]interface{}
	err := gh.ReadJSONFile(projectPath, fileName, &content)
	if err != nil {
		return "", err
	}

	if version, ok := content[versionField].(string); ok {
		return version, nil
	}

	return "", fmt.Errorf("version field '%s' not found or not a string", versionField)
}

// ParseVersionFromText parses version information from text file
func (bp *BaseProvider) ParseVersionFromText(
	projectPath string,
	fileName string,
	gitHandler interface{},
	pattern *regexp.Regexp,
) (string, error) {
	gh := gitHandler.(*git.GitHandler)
	content, err := gh.ReadFile(projectPath, fileName)
	if err != nil {
		return "", err
	}

	matches := pattern.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("version pattern not found in file")
}

// DetectPackageManager detects package manager
func (bp *BaseProvider) DetectPackageManager(
	files []types.FileInfo,
	lockFiles map[string]string,
) string {
	for lockFile, manager := range lockFiles {
		if bp.HasFile(files, lockFile) {
			return manager
		}
	}
	return ""
}

// CreateDetectResult creates basic detection result
func (bp *BaseProvider) CreateDetectResult(
	detected bool,
	confidence float64,
	language string,
	version *types.VersionInfo,
	framework string,
	packageManager string,
	buildTool string,
	metadata map[string]interface{},
	evidence types.Evidence,
) *types.DetectResult {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	versionStr := ""
	if version != nil {
		versionStr = version.Version
	}

	var packageMgr *types.PackageManager
	if packageManager != "" {
		packageMgr = &types.PackageManager{Name: packageManager}
	}

	var buildTools []string
	if buildTool != "" {
		buildTools = []string{buildTool}
	}

	return &types.DetectResult{
		Matched:        detected,
		Confidence:     confidence,
		Language:       language,
		Version:        versionStr,
		Framework:      framework,
		PackageManager: packageMgr,
		BuildTools:     buildTools,
		Metadata:       metadata,
		Evidence:       evidence,
	}
}

// CreateVersionInfo creates version information
func (bp *BaseProvider) CreateVersionInfo(version string, source string) *types.VersionInfo {
	return &types.VersionInfo{
		Version: version,
		Source:  source,
	}
}

// NormalizeVersion normalizes version number
func (bp *BaseProvider) NormalizeVersion(version string) string {
	// Remove prefix characters (like v, ^, ~, >=, etc.)
	re := regexp.MustCompile(`^[v^~>=<]+`)
	cleaned := re.ReplaceAllString(version, "")

	// Extract main version number (x.y.z format)
	versionRe := regexp.MustCompile(`^(\d+)(?:\.(\d+))?(?:\.(\d+))?`)
	matches := versionRe.FindStringSubmatch(cleaned)
	if len(matches) > 1 {
		major := matches[1]
		minor := "0"
		patch := "0"

		if len(matches) > 2 && matches[2] != "" {
			minor = matches[2]
		}
		if len(matches) > 3 && matches[3] != "" {
			patch = matches[3]
		}

		return fmt.Sprintf("%s.%s.%s", major, minor, patch)
	}

	return cleaned
}

// IsVersionCompatible checks if version meets minimum requirements
func (bp *BaseProvider) IsVersionCompatible(version string, minVersion string) bool {
	normalize := func(v string) []int {
		normalized := bp.NormalizeVersion(v)
		parts := strings.Split(normalized, ".")
		result := make([]int, 3)

		for i := 0; i < 3 && i < len(parts); i++ {
			if num, err := strconv.Atoi(parts[i]); err == nil {
				result[i] = num
			}
		}

		return result
	}

	current := normalize(version)
	minimum := normalize(minVersion)

	for i := 0; i < 3; i++ {
		if current[i] > minimum[i] {
			return true
		}
		if current[i] < minimum[i] {
			return false
		}
	}

	return true // versions are equal
}

// DetectFrameworkFromDependencies detects framework from dependencies
func (bp *BaseProvider) DetectFrameworkFromDependencies(
	dependencies map[string]interface{},
	frameworkMap map[string]string,
) string {
	allDeps := make(map[string]interface{})

	// Merge all dependencies
	if deps, ok := dependencies["dependencies"].(map[string]interface{}); ok {
		for k, v := range deps {
			allDeps[k] = v
		}
	}
	if devDeps, ok := dependencies["devDependencies"].(map[string]interface{}); ok {
		for k, v := range devDeps {
			allDeps[k] = v
		}
	}
	if peerDeps, ok := dependencies["peerDependencies"].(map[string]interface{}); ok {
		for k, v := range peerDeps {
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

// CalculateConfidence calculates confidence
func (bp *BaseProvider) CalculateConfidence(indicators []types.ConfidenceIndicator) float64 {
	totalWeight := 0
	satisfiedWeight := 0

	for _, indicator := range indicators {
		totalWeight += indicator.Weight
		if indicator.Satisfied {
			satisfiedWeight += indicator.Weight
		}
	}

	if totalWeight > 0 {
		return float64(satisfiedWeight) / float64(totalWeight)
	}

	return 0.0
}

// SafeReadJSON safely reads JSON file
func (bp *BaseProvider) SafeReadJSON(
	projectPath string,
	fileName string,
	gitHandler interface{},
) (map[string]interface{}, error) {
	gh := gitHandler.(*git.GitHandler)
	var content map[string]interface{}
	err := gh.ReadJSONFile(projectPath, fileName, &content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// SafeReadJSONC safely reads JSONC file (JSON with comments support)
func (bp *BaseProvider) SafeReadJSONC(
	projectPath string,
	fileName string,
	gitHandler interface{},
) (map[string]interface{}, error) {
	gh := gitHandler.(*git.GitHandler)
	var content map[string]interface{}
	err := gh.ReadJSONCFile(projectPath, fileName, &content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// HasFileInEvidence checks if Evidence.Files string array contains specified file
func (bp *BaseProvider) HasFileInEvidence(evidenceFiles []string, fileName string) bool {
	for _, file := range evidenceFiles {
		if file == fileName {
			return true
		}
	}
	return false
}

// SafeReadText safely reads text file
func (bp *BaseProvider) SafeReadText(
	projectPath string,
	fileName string,
	gitHandler interface{},
) (string, error) {
	gh := gitHandler.(*git.GitHandler)
	content, err := gh.ReadFile(projectPath, fileName)
	if err != nil {
		return "", err
	}
	return content, nil
}
