package utils

import (
	"regexp"
	"testing"

	"github.com/labring/devbox-pack/pkg/types"
)

// Test color functions
func TestColorFunctions(t *testing.T) {
	testCases := []struct {
		name     string
		function func(string) string
		input    string
		expected string
	}{
		{"Red", Red, "test", ColorRed + "test" + ColorReset},
		{"Green", Green, "test", ColorGreen + "test" + ColorReset},
		{"Yellow", Yellow, "test", ColorYellow + "test" + ColorReset},
		{"Blue", Blue, "test", ColorBlue + "test" + ColorReset},
		{"Magenta", Magenta, "test", ColorMagenta + "test" + ColorReset},
		{"Cyan", Cyan, "test", ColorCyan + "test" + ColorReset},
		{"Gray", Gray, "test", ColorGray + "test" + ColorReset},
		{"Bold", Bold, "test", ColorBold + "test" + ColorReset},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.function(tc.input)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestColorConstants(t *testing.T) {
	expectedColors := map[string]string{
		"ColorReset":   "\033[0m",
		"ColorRed":     "\033[31m",
		"ColorGreen":   "\033[32m",
		"ColorYellow":  "\033[33m",
		"ColorBlue":    "\033[34m",
		"ColorMagenta": "\033[35m",
		"ColorCyan":    "\033[36m",
		"ColorGray":    "\033[37m",
		"ColorBold":    "\033[1m",
	}

	actualColors := map[string]string{
		"ColorReset":   ColorReset,
		"ColorRed":     ColorRed,
		"ColorGreen":   ColorGreen,
		"ColorYellow":  ColorYellow,
		"ColorBlue":    ColorBlue,
		"ColorMagenta": ColorMagenta,
		"ColorCyan":    ColorCyan,
		"ColorGray":    ColorGray,
		"ColorBold":    ColorBold,
	}

	for name, expected := range expectedColors {
		if actual := actualColors[name]; actual != expected {
			t.Errorf("constant %s: expected %q, got %q", name, expected, actual)
		}
	}
}

// Test BaseCatalog
func TestBaseCatalog(t *testing.T) {
	// Test that all supported languages have entries
	expectedLanguages := []types.SupportedLanguage{
		types.LanguageNode,
		types.LanguagePython,
		types.LanguageJava,
		types.LanguageGo,
		types.LanguagePHP,
		types.LanguageRuby,
		types.LanguageDeno,
		types.LanguageRust,
		types.LanguageStaticfile,
		types.LanguageShell,
	}

	for _, lang := range expectedLanguages {
		if _, exists := BaseCatalog[lang]; !exists {
			t.Errorf("BaseCatalog missing entry for language %s", lang)
		}
	}

	// Test specific entries
	if nodeImages, exists := BaseCatalog[types.LanguageNode]; exists {
		if nodeImages["20"] != "node:20-alpine" {
			t.Errorf("expected Node.js 20 image to be 'node:20-alpine', got %s", nodeImages["20"])
		}
	}

	if pythonImages, exists := BaseCatalog[types.LanguagePython]; exists {
		if pythonImages["3.11"] != "python:3.11-slim" {
			t.Errorf("expected Python 3.11 image to be 'python:3.11-slim', got %s", pythonImages["3.11"])
		}
	}
}

// Test DefaultPorts
func TestDefaultPorts(t *testing.T) {
	// Test language default ports
	expectedLanguagePorts := map[types.SupportedLanguage]int{
		types.LanguageNode:   3000,
		types.LanguagePython: 8000,
		types.LanguageJava:   8080,
		types.LanguageGo:     8080,
		types.LanguagePHP:    8000,
		types.LanguageRuby:   3000,
		types.LanguageDeno:   8000,
		types.LanguageRust:   8080,
	}

	for lang, expectedPort := range expectedLanguagePorts {
		if actualPort := DefaultPorts.Languages[lang]; actualPort != expectedPort {
			t.Errorf("language %s: expected port %d, got %d", lang, expectedPort, actualPort)
		}
	}

	// Test framework default ports
	expectedFrameworkPorts := map[string]int{
		"Next.js":       3000,
		"Django":        8000,
		"Spring Boot":   8080,
		"Gin":           8080,
		"Laravel":       8000,
		"Ruby on Rails": 3000,
	}

	for framework, expectedPort := range expectedFrameworkPorts {
		if actualPort := DefaultPorts.Frameworks[framework]; actualPort != expectedPort {
			t.Errorf("framework %s: expected port %d, got %d", framework, expectedPort, actualPort)
		}
	}
}

// Test PackageManagers
func TestPackageManagers(t *testing.T) {
	// Test that all languages have package manager entries
	for lang := range BaseCatalog {
		if _, exists := PackageManagers[lang]; !exists {
			t.Errorf("PackageManagers missing entry for language %s", lang)
		}
	}

	// Test specific package managers
	nodeManagers := PackageManagers[types.LanguageNode]
	expectedNodeManagers := []string{"npm", "yarn", "pnpm", "bun"}
	if len(nodeManagers) != len(expectedNodeManagers) {
		t.Errorf("expected %d Node.js package managers, got %d", len(expectedNodeManagers), len(nodeManagers))
	}

	pythonManagers := PackageManagers[types.LanguagePython]
	expectedPythonManagers := []string{"pip", "poetry", "pipenv", "conda"}
	if len(pythonManagers) != len(expectedPythonManagers) {
		t.Errorf("expected %d Python package managers, got %d", len(expectedPythonManagers), len(pythonManagers))
	}
}

// Test BuildTools
func TestBuildTools(t *testing.T) {
	// Test that all languages have build tool entries
	for lang := range BaseCatalog {
		if _, exists := BuildTools[lang]; !exists {
			t.Errorf("BuildTools missing entry for language %s", lang)
		}
	}

	// Test specific build tools
	nodeTools := BuildTools[types.LanguageNode]
	expectedNodeTools := []string{"webpack", "vite", "rollup", "parcel", "esbuild", "turbo", "nx"}
	if len(nodeTools) != len(expectedNodeTools) {
		t.Errorf("expected %d Node.js build tools, got %d", len(expectedNodeTools), len(nodeTools))
	}

	javaTools := BuildTools[types.LanguageJava]
	expectedJavaTools := []string{"maven", "gradle", "sbt", "ant"}
	if len(javaTools) != len(expectedJavaTools) {
		t.Errorf("expected %d Java build tools, got %d", len(expectedJavaTools), len(javaTools))
	}
}

// Test ScanConfig
func TestScanConfig(t *testing.T) {
	if ScanConfig.DefaultDepth != 3 {
		t.Errorf("expected default depth 3, got %d", ScanConfig.DefaultDepth)
	}

	if ScanConfig.MaxFiles != 1000 {
		t.Errorf("expected max files 1000, got %d", ScanConfig.MaxFiles)
	}

	// Test that important directories are ignored
	expectedIgnoreDirs := []string{"node_modules", ".git", "vendor", "target", "build", "dist"}
	for _, dir := range expectedIgnoreDirs {
		found := false
		for _, ignoreDir := range ScanConfig.IgnoreDirs {
			if ignoreDir == dir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected ignore directory %s not found", dir)
		}
	}

	// Test that important files are included
	expectedImportantFiles := []string{"package.json", "requirements.txt", "pom.xml", "go.mod", "Cargo.toml"}
	for _, file := range expectedImportantFiles {
		found := false
		for _, importantFile := range ScanConfig.ImportantFiles {
			if importantFile == file {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected important file %s not found", file)
		}
	}
}

// Test CLIDefaults
func TestCLIDefaults(t *testing.T) {
	if CLIDefaults.Depth != 3 {
		t.Errorf("expected default depth 3, got %d", CLIDefaults.Depth)
	}

	if CLIDefaults.MaxFiles != 1000 {
		t.Errorf("expected default max files 1000, got %d", CLIDefaults.MaxFiles)
	}

	if CLIDefaults.OutputFormat != types.OutputFormatJSON {
		t.Errorf("expected default output format JSON, got %s", CLIDefaults.OutputFormat)
	}

	if CLIDefaults.ConfidenceThreshold != 20 {
		t.Errorf("expected default confidence threshold 20, got %d", CLIDefaults.ConfidenceThreshold)
	}

	if CLIDefaults.Timeout != 30000 {
		t.Errorf("expected default timeout 30000, got %d", CLIDefaults.Timeout)
	}

	if CLIDefaults.TempDirPrefix != "devbox-pack-" {
		t.Errorf("expected default temp dir prefix 'devbox-pack-', got %s", CLIDefaults.TempDirPrefix)
	}

	if CLIDefaults.LogLevel != "info" {
		t.Errorf("expected default log level 'info', got %s", CLIDefaults.LogLevel)
	}
}

// Test DefaultVersions
func TestDefaultVersions(t *testing.T) {
	expectedVersions := map[types.SupportedLanguage]string{
		types.LanguageNode:   "20",
		types.LanguagePython: "3.11",
		types.LanguageJava:   "17",
		types.LanguageGo:     "1.21",
		types.LanguagePHP:    "8.2",
		types.LanguageRuby:   "3.2",
		types.LanguageDeno:   "1.40",
		types.LanguageRust:   "1.70",
	}

	for lang, expectedVersion := range expectedVersions {
		if actualVersion := DefaultVersions[lang]; actualVersion != expectedVersion {
			t.Errorf("language %s: expected default version %s, got %s", lang, expectedVersion, actualVersion)
		}
	}
}

// Test VersionPatterns
func TestVersionPatterns(t *testing.T) {
	// Test Semver pattern
	semverTests := []struct {
		version string
		valid   bool
	}{
		{"1.2.3", true},
		{"1.2.3-alpha", true},
		{"1.2.3+build", true},
		{"1.2.3-alpha+build", true},
		{"1.2", false},
		{"1", false},
		{"invalid", false},
	}

	for _, test := range semverTests {
		matches := VersionPatterns.Semver.MatchString(test.version)
		if matches != test.valid {
			t.Errorf("semver pattern for %s: expected %t, got %t", test.version, test.valid, matches)
		}
	}

	// Test MajorMinor pattern
	majorMinorTests := []struct {
		version string
		valid   bool
	}{
		{"1.2", true},
		{"18.17", true},
		{"1.2.3", false},
		{"1", false},
		{"invalid", false},
	}

	for _, test := range majorMinorTests {
		matches := VersionPatterns.MajorMinor.MatchString(test.version)
		if matches != test.valid {
			t.Errorf("major.minor pattern for %s: expected %t, got %t", test.version, test.valid, matches)
		}
	}

	// Test MajorOnly pattern
	majorOnlyTests := []struct {
		version string
		valid   bool
	}{
		{"1", true},
		{"18", true},
		{"1.2", false},
		{"invalid", false},
	}

	for _, test := range majorOnlyTests {
		matches := VersionPatterns.MajorOnly.MatchString(test.version)
		if matches != test.valid {
			t.Errorf("major only pattern for %s: expected %t, got %t", test.version, test.valid, matches)
		}
	}

	// Test NodeVersion pattern
	nodeVersionTests := []struct {
		version string
		valid   bool
	}{
		{"18.17.0", true},
		{"v18.17.0", true},
		{"20.0.0", true},
		{"v20.0.0", true},
		{"18.17", false},
		{"v18", false},
		{"invalid", false},
	}

	for _, test := range nodeVersionTests {
		matches := VersionPatterns.NodeVersion.MatchString(test.version)
		if matches != test.valid {
			t.Errorf("node version pattern for %s: expected %t, got %t", test.version, test.valid, matches)
		}
	}
}

// Test that all patterns are valid regex
func TestVersionPatternsValidRegex(t *testing.T) {
	patterns := map[string]*regexp.Regexp{
		"Semver":        VersionPatterns.Semver,
		"MajorMinor":    VersionPatterns.MajorMinor,
		"MajorOnly":     VersionPatterns.MajorOnly,
		"NodeVersion":   VersionPatterns.NodeVersion,
		"PythonVersion": VersionPatterns.PythonVersion,
		"JavaVersion":   VersionPatterns.JavaVersion,
		"GoVersion":     VersionPatterns.GoVersion,
	}

	for name, pattern := range patterns {
		if pattern == nil {
			t.Errorf("pattern %s is nil", name)
		}
	}
}

// Test data consistency
func TestDataConsistency(t *testing.T) {
	// Test that all languages in BaseCatalog have corresponding entries in other maps
	for lang := range BaseCatalog {
		if _, exists := DefaultPorts.Languages[lang]; !exists {
			t.Errorf("language %s missing from DefaultPorts.Languages", lang)
		}

		if _, exists := PackageManagers[lang]; !exists {
			t.Errorf("language %s missing from PackageManagers", lang)
		}

		if _, exists := BuildTools[lang]; !exists {
			t.Errorf("language %s missing from BuildTools", lang)
		}

		if _, exists := DefaultVersions[lang]; !exists {
			t.Errorf("language %s missing from DefaultVersions", lang)
		}
	}

	// Test that default versions exist in BaseCatalog
	for lang, version := range DefaultVersions {
		if images, exists := BaseCatalog[lang]; exists {
			if _, versionExists := images[version]; !versionExists {
				t.Errorf("default version %s for language %s not found in BaseCatalog", version, lang)
			}
		}
	}
}
