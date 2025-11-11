package providers

import (
	"testing"

	"github.com/labring/devbox-pack/pkg/types"
)

func TestNewJavaProvider(t *testing.T) {
	provider := NewJavaProvider()
	AssertProviderBasic(t, provider, "java", 30)
}

func TestJavaProvider_Detect_NoJavaFiles(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	provider := NewJavaProvider()

	files := []types.FileInfo{
		{Path: "README.md", IsDirectory: false},
		{Path: "package.json", IsDirectory: false},
	}

	result, err := provider.Detect(helper.TempDir, files, helper.GitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	AssertDetectResult(t, result, false, "")
}

func TestJavaProvider_Detect_WithMaven(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	provider := NewJavaProvider()

	// Create pom.xml
	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
    </properties>
</project>`

	helper.WriteFile("pom.xml", pomContent)
	helper.CreateTempDir("src/main/java")

	files := []types.FileInfo{
		{Path: "pom.xml", IsDirectory: false},
		{Path: "src/main/java/", IsDirectory: true},
		{Path: "src/main/java/Main.java", IsDirectory: false},
	}

	result, err := provider.Detect(helper.TempDir, files, helper.GitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	AssertDetectResult(t, result, true, "java")

	if result.Version != "17" {
		t.Errorf("expected version '17', got %s", result.Version)
	}

	if len(result.BuildTools) == 0 || result.BuildTools[0] != "Maven" {
		t.Errorf("expected build tool 'Maven', got %v", result.BuildTools)
	}
}

func TestJavaProvider_Detect_WithGradle(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	provider := NewJavaProvider()

	// Create build.gradle
	gradleContent := `plugins {
    id 'java'
}

java {
    sourceCompatibility = JavaVersion.VERSION_17
    targetCompatibility = JavaVersion.VERSION_17
}`

	helper.WriteFile("build.gradle", gradleContent)
	helper.CreateTempDir("src/main/java")

	files := []types.FileInfo{
		{Path: "build.gradle", IsDirectory: false},
		{Path: "src/main/java/", IsDirectory: true},
		{Path: "src/main/java/Main.java", IsDirectory: false},
	}

	result, err := provider.Detect(helper.TempDir, files, helper.GitHandler)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	AssertDetectResult(t, result, true, "java")

	if len(result.BuildTools) == 0 || result.BuildTools[0] != "Gradle" {
		t.Errorf("expected build tool 'Gradle', got %v", result.BuildTools)
	}
}

func TestJavaProvider_GenerateCommands_MavenProject(t *testing.T) {
	provider := NewJavaProvider()

	result := &types.DetectResult{
		Matched:    true,
		Language:   "java",
		BuildTools: []string{"Maven"},
		Framework:  "Spring Boot",
		Evidence: types.Evidence{
			Files: []string{"pom.xml", "src/main/java/Main.java"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Setup) == 0 {
		t.Error("expected setup commands for Maven project")
	}

	AssertCommandContains(t, commands.Setup, "mvn")

	if len(commands.Build) == 0 {
		t.Error("expected build commands for Maven project")
	}

	AssertCommandContains(t, commands.Build, "mvn")

	if len(commands.Run) == 0 {
		t.Error("expected run commands for Maven project")
	}
}

func TestJavaProvider_GenerateCommands_GradleProject(t *testing.T) {
	provider := NewJavaProvider()

	result := &types.DetectResult{
		Matched:    true,
		Language:   "java",
		BuildTools: []string{"Gradle"},
		Framework:  "Spring Boot",
		Evidence: types.Evidence{
			Files: []string{"build.gradle", "src/main/java/Main.java"},
		},
	}

	options := types.CLIOptions{}
	commands := provider.GenerateCommands(result, options)

	if len(commands.Setup) == 0 {
		t.Error("expected setup commands for Gradle project")
	}

	AssertCommandContains(t, commands.Setup, "gradle")

	if len(commands.Build) == 0 {
		t.Error("expected build commands for Gradle project")
	}

	AssertCommandContains(t, commands.Build, "gradle")
}

func TestJavaProvider_GenerateEnvironment(t *testing.T) {
	provider := NewJavaProvider()

	result := &types.DetectResult{
		Matched:    true,
		Language:   "java",
		Version:    "17",
		BuildTools: []string{"Maven"},
		Framework:  "Spring Boot",
	}

	env := provider.GenerateEnvironment(result)

	if env == nil {
		t.Fatal("GenerateEnvironment returned nil")
	}

	AssertEnvironmentVar(t, env, "JAVA_VERSION", "17")
}

func TestJavaProvider_NeedsNativeCompilation(t *testing.T) {
	provider := NewJavaProvider()

	tests := []struct {
		name     string
		result   *types.DetectResult
		expected bool
	}{
		{
			name: "Spring Boot project",
			result: &types.DetectResult{
				Framework: "Spring Boot",
			},
			expected: false,
		},
		{
			name: "GraalVM project",
			result: &types.DetectResult{
				Framework: "GraalVM",
			},
			expected: true,
		},
		{
			name: "Quarkus project",
			result: &types.DetectResult{
				Framework: "Quarkus",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertNeedsNativeCompilation(t, provider, tt.result, tt.expected)
		})
	}
}
