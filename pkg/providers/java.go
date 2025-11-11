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
			Priority: 30,
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

	// Detect build tool
	buildTool := p.DetectPackageManager(files, map[string]string{
		"pom.xml":          "Maven",
		"build.gradle":     "Gradle",
		"build.gradle.kts": "Gradle",
	})
	if buildTool == "" {
		buildTool = "Maven"
	}

	// For package manager, use lowercase
	packageManager := strings.ToLower(buildTool)

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
		buildTool,
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
	// Priority-based framework detection - Spring Boot has highest priority
	frameworkMap := map[string]string{
		"spring-boot-starter":           "Spring Boot",
		"spring-boot-starter-web":       "Spring Boot Web",
		"spring-boot-starter-data-jpa":  "Spring Boot Data JPA",
		"spring-boot-starter-security":  "Spring Boot Security",
		"spring-boot":                   "Spring Boot",
		"spring-webmvc":                 "Spring MVC",
		"spring-web":                    "Spring Web",
		"spring-core":                   "Spring Framework",
		"quarkus":                       "Quarkus",
		"quarkus-resteasy":              "Quarkus RESTEasy",
		"micronaut":                     "Micronaut",
		"micronaut-http-server":         "Micronaut HTTP",
		"vertx":                         "Vert.x",
		"vertx-web":                     "Vert.x Web",
		"dropwizard":                    "Dropwizard",
		"dropwizard-core":               "Dropwizard Core",
		"spark-core":                    "Spark Java",
		"jersey":                        "Jersey",
		"struts":                        "Struts",
		"wicket":                        "Apache Wicket",
		"vaadin":                        "Vaadin",
		"jakarta.servlet":               "Jakarta Servlet",
		"jakarta.ws.rs":                 "Jakarta JAX-RS",
		"jakarta.persistence":           "Jakarta Persistence",
	}

	// Check dependencies in pom.xml with priority ordering
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
		// Check for Spring Boot parent first (highest priority)
		if strings.Contains(pomContent, "<parent>") && strings.Contains(pomContent, "spring-boot-starter-parent") {
			return "Spring Boot", nil
		}

		// Check dependencies in priority order
		for dependency, framework := range frameworkMap {
			if strings.Contains(pomContent, dependency) {
				return framework, nil
			}
		}
	}

	// Check dependencies in build.gradle with priority ordering
	gradleContent, err := p.SafeReadText(projectPath, "build.gradle", gitHandler)
	if err != nil {
		// If file doesn't exist, continue trying build.gradle.kts
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			// Continue trying build.gradle.kts
		} else {
			return "", err
		}
	}
	if gradleContent != "" {
		// Check for Spring Boot plugin first (highest priority)
		if strings.Contains(gradleContent, "org.springframework.boot") {
			return "Spring Boot", nil
		}

		// Check dependencies in priority order
		for dependency, framework := range frameworkMap {
			if strings.Contains(gradleContent, dependency) {
				return framework, nil
			}
		}
	}

	// Check dependencies in build.gradle.kts
	gradleKtsContent, err := p.SafeReadText(projectPath, "build.gradle.kts", gitHandler)
	if err != nil {
		// If file doesn't exist, return empty string instead of error
		if strings.Contains(err.Error(), "FILE_READ_ERROR") {
			return "", nil
		}
		return "", err
	}
	if gradleKtsContent != "" {
		// Check for Spring Boot plugin first (highest priority)
		if strings.Contains(gradleKtsContent, "org.springframework.boot") {
			return "Spring Boot", nil
		}

		// Check dependencies in priority order
		for dependency, framework := range frameworkMap {
			if strings.Contains(gradleKtsContent, dependency) {
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

	// Check if Gradle files exist
	hasGradle := p.HasFileInEvidence(result.Evidence.Files, "build.gradle") || p.HasFileInEvidence(result.Evidence.Files, "build.gradle.kts")

	// Check for Spring Boot specifically
	isSpringBoot := strings.Contains(result.Framework, "Spring Boot")

	if hasPom {
		// Maven project with framework-specific commands
		if isSpringBoot {
			commands.Setup = []string{"mvn clean compile"}
			commands.Dev = []string{"mvn spring-boot:run"}
			commands.Build = []string{"mvn clean package -DskipTests"}
			commands.Run = []string{"java -jar target/*.jar"}
		} else if result.Framework == "Quarkus" {
			commands.Setup = []string{"mvn clean compile"}
			commands.Dev = []string{"mvn quarkus:dev"}
			commands.Build = []string{"mvn clean package -DskipTests"}
			commands.Run = []string{"java -jar target/quarkus-app/quarkus-app.jar"}
		} else if result.Framework == "Micronaut" {
			commands.Setup = []string{"mvn clean compile"}
			commands.Dev = []string{"mvn mn:run"}
			commands.Build = []string{"mvn clean package -DskipTests"}
			commands.Run = []string{"java -jar target/*.jar"}
		} else {
			// Generic Maven project
			commands.Setup = []string{"mvn clean compile"}
			commands.Dev = []string{"mvn exec:java"}
			commands.Build = []string{"mvn clean package -DskipTests"}
			commands.Run = []string{"java -jar target/*.jar"}
		}
	} else if hasGradle {
		// Gradle project with framework-specific commands
		if isSpringBoot {
			commands.Setup = []string{"./gradlew compileJava"}
			commands.Dev = []string{"./gradlew bootRun"}
			commands.Build = []string{"./gradlew build -x test"}
			commands.Run = []string{"java -jar build/libs/*.jar"}
		} else if result.Framework == "Quarkus" {
			commands.Setup = []string{"./gradlew compileJava"}
			commands.Dev = []string{"./gradlew quarkusDev"}
			commands.Build = []string{"./gradlew build -x test"}
			commands.Run = []string{"java -jar build/quarkus-app/quarkus-app.jar"}
		} else if result.Framework == "Micronaut" {
			commands.Setup = []string{"./gradlew compileJava"}
			commands.Dev = []string{"./gradlew run"}
			commands.Build = []string{"./gradlew build -x test"}
			commands.Run = []string{"java -jar build/libs/*.jar"}
		} else {
			// Generic Gradle project
			commands.Setup = []string{"./gradlew compileJava"}
			commands.Dev = []string{"./gradlew run"}
			commands.Build = []string{"./gradlew build -x test"}
			commands.Run = []string{"java -jar build/libs/*.jar"}
		}
	}

	return commands
}

// GenerateEnvironment generates environment variables for Java project
func (p *JavaProvider) GenerateEnvironment(result *types.DetectResult) map[string]string {
	env := make(map[string]string)

	// Set Java specific environment variables
	env["JAVA_OPTS"] = "-Xmx512m"

	// Set port for web applications
	env["PORT"] = "8080"
	env["SERVER_PORT"] = "8080"

	// Add Java version if available
	if result.Version != "" {
		env["JAVA_VERSION"] = result.Version
	}

	// Add framework-specific environment variables
	if strings.Contains(result.Framework, "Spring Boot") {
		env["SPRING_PROFILES_ACTIVE"] = "production"
		env["SPRING_JPA_HIBERNATE_DDL_AUTO"] = "update"
		env["SPRING_DATASOURCE_URL"] = "jdbc:h2:mem:testdb"
		env["SPRING_H2_CONSOLE_ENABLED"] = "false"
		env["SPRING_OUTPUT_ANSI_ENABLED"] = "always"
	} else if result.Framework == "Quarkus" {
		env["QUARKUS_PROFILE"] = "prod"
		env["QUARKUS_DATASOURCE_URL"] = "jdbc:h2:mem:testdb"
		env["QUARKUS_HTTP_PORT"] = "8080"
	} else if result.Framework == "Micronaut" {
		env["MICRONAUT_ENVIRONMENTS"] = "prod"
		env["MICRONAUT_SERVER_PORT"] = "8080"
	} else if result.Framework == "Vert.x" {
		env["VERTX_OPTIONS"] = "{}"
	}

	// Add build tool specific environment variables
	if result.Metadata != nil {
		if buildTool, ok := result.Metadata["packageManager"].(string); ok {
			if buildTool == "maven" {
				env["MAVEN_OPTS"] = "-DskipTests"
			} else if buildTool == "gradle" {
				env["GRADLE_OPTS"] = "-DskipTests=true"
			}
		}
	}

	return env
}

// NeedsNativeCompilation checks if Java project needs native compilation
func (p *JavaProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Java projects usually don't need native compilation, unless using GraalVM Native Image
	// Check for frameworks that typically use native compilation
	if result.Framework == "GraalVM" || result.Framework == "Quarkus" {
		return true
	}

	// Check metadata for native compilation flags
	if result.Metadata != nil {
		if needsNative, ok := result.Metadata["needsNativeCompilation"].(bool); ok {
			return needsNative
		}
	}

	return false
}
