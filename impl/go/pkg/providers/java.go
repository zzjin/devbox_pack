/**
 * DevBox Pack Execution Plan Generator - Java Provider
 */

package providers

import (
	"regexp"
	"strings"

	"github.com/labring/devbox-pack/pkg/types"
)

// JavaProvider Java project detector
type JavaProvider struct {
	BaseProvider
}

// NewJavaProvider creates Java Provider
func NewJavaProvider() *JavaProvider {
	return &JavaProvider{
		BaseProvider: BaseProvider{
			Name:     "java",
			Language: "java",
			Priority: 70,
		},
	}
}

// GetName gets Provider name
func (p *JavaProvider) GetName() string {
	return p.Name
}

// GetPriority gets Provider priority
func (p *JavaProvider) GetPriority() int {
	return p.Priority
}

// Detect detects Java project
func (p *JavaProvider) Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error) {
	indicators := []types.ConfidenceIndicator{
		{Weight: 30, Satisfied: p.HasAnyFile(files, []string{"pom.xml", "build.gradle", "build.gradle.kts"})},
		{Weight: 25, Satisfied: p.HasAnyFile(files, []string{"*.java", "*.kt", "*.scala"})},
		{Weight: 15, Satisfied: p.HasFile(files, "gradle.properties")},
		{Weight: 10, Satisfied: p.HasFile(files, "gradlew")},
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"src/main/java", "src/main/kotlin", "src/main/scala"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"mvnw", "mvnw.cmd"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{".mvn", ".gradle"})},
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.2 // Lower detection threshold

	if !detected {
		return p.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Detect version
	version, err := p.detectJavaVersion(projectPath, gitHandler)
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
		"pom.xml":          "maven",
		"build.gradle":     "gradle",
		"build.gradle.kts": "gradle",
	})
	if packageManager == "" {
		packageManager = "maven"
	}

	metadata := map[string]interface{}{
		"hasPom":           p.HasFile(files, "pom.xml"),
		"hasGradle":        p.HasAnyFile(files, []string{"build.gradle", "build.gradle.kts"}),
		"hasGradleWrapper": p.HasFile(files, "gradlew"),
		"hasMavenWrapper":  p.HasFile(files, "mvnw"),
		"packageManager":   packageManager,
		"framework":        framework,
	}

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	if p.HasFile(files, "pom.xml") {
		evidenceFiles = append(evidenceFiles, "pom.xml")
	}
	if p.HasFile(files, "build.gradle") {
		evidenceFiles = append(evidenceFiles, "build.gradle")
	}
	if p.HasFile(files, "build.gradle.kts") {
		evidenceFiles = append(evidenceFiles, "build.gradle.kts")
	}
	if p.HasFile(files, "gradlew") {
		evidenceFiles = append(evidenceFiles, "gradlew")
	}
	if p.HasFile(files, "mvnw") {
		evidenceFiles = append(evidenceFiles, "mvnw")
	}
	if p.HasFile(files, "gradle.properties") {
		evidenceFiles = append(evidenceFiles, "gradle.properties")
	}
	if p.HasFile(files, "settings.gradle") {
		evidenceFiles = append(evidenceFiles, "settings.gradle")
	}
	if p.HasFile(files, "settings.gradle.kts") {
		evidenceFiles = append(evidenceFiles, "settings.gradle.kts")
	}

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected Java project based on: "
	var reasons []string
	if p.HasAnyFile(files, []string{"*.java"}) {
		reasons = append(reasons, "Java source files")
	}
	if p.HasFile(files, "pom.xml") {
		reasons = append(reasons, "Maven configuration (pom.xml)")
	}
	if p.HasAnyFile(files, []string{"build.gradle", "build.gradle.kts"}) {
		reasons = append(reasons, "Gradle configuration")
	}
	if packageManager != "" {
		reasons = append(reasons, "build tool: "+packageManager)
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
		"java",
		version,
		framework,
		packageManager,
		packageManager,
		metadata,
		evidence,
	), nil
}

// detectJavaVersion detects Java version
func (p *JavaProvider) detectJavaVersion(projectPath string, gitHandler interface{}) (*types.VersionInfo, error) {
	// Read from pom.xml
	pomContent, err := p.SafeReadText(projectPath, "pom.xml", gitHandler)
	if err != nil {
		// If file doesn't exist, continue trying other methods
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			// Continue trying build.gradle
		} else {
			return nil, err
		}
	}
	if pomContent != "" {
		// Check maven.compiler.source
		re := regexp.MustCompile(`<maven\.compiler\.source>([^<]+)</maven\.compiler\.source>`)
		matches := re.FindStringSubmatch(pomContent)
		if len(matches) > 1 {
			return p.CreateVersionInfo(matches[1], "pom.xml maven.compiler.source"), nil
		}

		// Check java.version
		re = regexp.MustCompile(`<java\.version>([^<]+)</java\.version>`)
		matches = re.FindStringSubmatch(pomContent)
		if len(matches) > 1 {
			return p.CreateVersionInfo(matches[1], "pom.xml java.version"), nil
		}

		// Check version in properties
		re = regexp.MustCompile(`<source>([^<]+)</source>`)
		matches = re.FindStringSubmatch(pomContent)
		if len(matches) > 1 {
			return p.CreateVersionInfo(matches[1], "pom.xml source"), nil
		}
	}

	// Read from build.gradle
	gradleContent, err := p.SafeReadText(projectPath, "build.gradle", gitHandler)
	if err != nil {
		// If file doesn't exist, continue trying other methods
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			// Continue trying build.gradle.kts
		} else {
			return nil, err
		}
	}
	if gradleContent != "" {
		re := regexp.MustCompile(`sourceCompatibility\s*=\s*['"]?([^'"\s]+)['"]?`)
		matches := re.FindStringSubmatch(gradleContent)
		if len(matches) > 1 {
			return p.CreateVersionInfo(matches[1], "build.gradle sourceCompatibility"), nil
		}

		// Check targetCompatibility
		re = regexp.MustCompile(`targetCompatibility\s*=\s*['"]?([^'"\s]+)['"]?`)
		matches = re.FindStringSubmatch(gradleContent)
		if len(matches) > 1 {
			return p.CreateVersionInfo(matches[1], "build.gradle targetCompatibility"), nil
		}
	}

	// Read from build.gradle.kts
	gradleKtsContent, err := p.SafeReadText(projectPath, "build.gradle.kts", gitHandler)
	if err != nil {
		// If file doesn't exist, use default version
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return p.CreateVersionInfo("17", "default"), nil
		}
		return nil, err
	}
	if gradleKtsContent != "" {
		re := regexp.MustCompile(`sourceCompatibility\s*=\s*JavaVersion\.VERSION_(\d+)`)
		matches := re.FindStringSubmatch(gradleKtsContent)
		if len(matches) > 1 {
			return p.CreateVersionInfo(matches[1], "build.gradle.kts sourceCompatibility"), nil
		}
	}

	// Default version
	return p.CreateVersionInfo("17", "default"), nil
}

// detectFramework detects framework
func (p *JavaProvider) detectFramework(projectPath string, gitHandler interface{}) (string, error) {
	frameworkMap := map[string]string{
		"spring-boot-starter": "Spring Boot",
		"spring-webmvc":       "Spring MVC",
		"spring-core":         "Spring Framework",
		"quarkus":             "Quarkus",
		"micronaut":           "Micronaut",
		"vertx":               "Vert.x",
		"dropwizard":          "Dropwizard",
		"spark-core":          "Spark Java",
		"jersey":              "Jersey",
		"struts":              "Struts",
		"wicket":              "Apache Wicket",
		"vaadin":              "Vaadin",
	}

	// Check dependencies in pom.xml
	pomContent, err := p.SafeReadText(projectPath, "pom.xml", gitHandler)
	if err != nil {
		// If file doesn't exist, continue trying other methods
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			// Continue trying build.gradle
		} else {
			return "", err
		}
	}
	if pomContent != "" {
		for dependency, framework := range frameworkMap {
			if strings.Contains(pomContent, dependency) {
				return framework, nil
			}
		}
	}

	// Check dependencies in build.gradle
	gradleContent, err := p.SafeReadText(projectPath, "build.gradle", gitHandler)
	if err != nil {
		// If file doesn't exist, return empty string instead of error
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return "", nil
		}
		return "", err
	}
	if gradleContent != "" {
		for dependency, framework := range frameworkMap {
			if strings.Contains(gradleContent, dependency) {
				return framework, nil
			}
		}
	}

	return "", nil
}

// GenerateCommands generates commands for Java project
func (p *JavaProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	// Check if pom.xml exists
	hasPom := p.HasFileInEvidence(result.Evidence.Files, "pom.xml")

	if hasPom {
		// Maven project
		commands.Dev = []string{
			"mvn clean compile",
			"mvn spring-boot:run",
		}
		commands.Build = []string{
			"mvn clean package",
		}
		commands.Start = []string{
			"java -jar target/*.jar",
		}
	} else {
		// Check if Gradle files exist
		hasGradle := p.HasFileInEvidence(result.Evidence.Files, "build.gradle") || p.HasFileInEvidence(result.Evidence.Files, "build.gradle.kts")

		if hasGradle {
			// Gradle project
			commands.Dev = []string{
				"./gradlew build",
				"./gradlew bootRun",
			}
			commands.Build = []string{
				"./gradlew build",
			}
			commands.Start = []string{
				"java -jar build/libs/*.jar",
			}
		}
	}

	return commands
}

// NeedsNativeCompilation checks if Java project needs native compilation
func (p *JavaProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Java projects usually don't need native compilation, unless using GraalVM Native Image
	// Check if there are GraalVM Native Image related configurations

	// Check if files exist
	if p.HasFileInEvidence(result.Evidence.Files, "pom.xml") ||
		p.HasFileInEvidence(result.Evidence.Files, "build.gradle") ||
		p.HasFileInEvidence(result.Evidence.Files, "build.gradle.kts") {
		// Here we can further check file content, but for simplicity, we assume most Java projects don't need native compilation
		return false
	}

	return false
}
