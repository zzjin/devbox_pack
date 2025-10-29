/**
 * DevBox Pack Execution Plan Generator - Constants Definition
 */

package utils

import (
	"regexp"

	"github.com/labring/devbox-pack/pkg/types"
)

// BaseCatalog Base image catalog
var BaseCatalog = map[types.SupportedLanguage]map[string]string{
	types.LanguageNode: {
		"16": "node:16-alpine",
		"18": "node:18-alpine",
		"20": "node:20-alpine",
		"21": "node:21-alpine",
	},
	types.LanguagePython: {
		"3.8":  "python:3.8-slim",
		"3.9":  "python:3.9-slim",
		"3.10": "python:3.10-slim",
		"3.11": "python:3.11-slim",
		"3.12": "python:3.12-slim",
	},
	types.LanguageJava: {
		"8":  "openjdk:8-jre-alpine",
		"11": "openjdk:11-jre-alpine",
		"17": "openjdk:17-jre-alpine",
		"21": "openjdk:21-jre-alpine",
	},
	types.LanguageGo: {
		"1.19": "golang:1.19-alpine",
		"1.20": "golang:1.20-alpine",
		"1.21": "golang:1.21-alpine",
		"1.22": "golang:1.22-alpine",
	},
	types.LanguagePHP: {
		"7.4": "php:7.4-fpm-alpine",
		"8.0": "php:8.0-fpm-alpine",
		"8.1": "php:8.1-fpm-alpine",
		"8.2": "php:8.2-fpm-alpine",
		"8.3": "php:8.3-fpm-alpine",
	},
	types.LanguageRuby: {
		"2.7": "ruby:2.7-alpine",
		"3.0": "ruby:3.0-alpine",
		"3.1": "ruby:3.1-alpine",
		"3.2": "ruby:3.2-alpine",
		"3.3": "ruby:3.3-alpine",
	},

	types.LanguageDeno: {
		"1.38": "denoland/deno:1.38.0",
		"1.39": "denoland/deno:1.39.0",
		"1.40": "denoland/deno:1.40.0",
		"1.41": "denoland/deno:1.41.0",
	},
	types.LanguageRust: {
		"1.68": "rust:1.68-alpine",
		"1.69": "rust:1.69-alpine",
		"1.70": "rust:1.70-alpine",
		"1.71": "rust:1.71-alpine",
	},
	types.LanguageStaticfile: {
		"1.0": "nginx:alpine",
	},
	types.LanguageShell: {
		"1.0": "alpine:latest",
	},
}

// DefaultPorts default port configuration
var DefaultPorts = struct {
	Languages  map[types.SupportedLanguage]int
	Frameworks map[string]int
}{
	Languages: map[types.SupportedLanguage]int{
		types.LanguageNode:   3000,
		types.LanguagePython: 8000,
		types.LanguageJava:   8080,
		types.LanguageGo:     8080,
		types.LanguagePHP:    8000,
		types.LanguageRuby:   3000,

		types.LanguageDeno:       8000,
		types.LanguageRust:       8080,
		types.LanguageStaticfile: 80,
		types.LanguageShell:      8000,
	},
	Frameworks: map[string]int{
		// Node.js frameworks
		"Express":            3000,
		"Koa":                3000,
		"Fastify":            3000,
		"NestJS":             3000,
		"Next.js":            3000,
		"Nuxt.js":            3000,
		"Gatsby":             8000,
		"Svelte":             5173,
		"SvelteKit":          5173,
		"Vite":               5173,
		"Webpack Dev Server": 8080,

		// Python frameworks
		"Django":  8000,
		"Flask":   5000,
		"FastAPI": 8000,
		"Tornado": 8888,
		"Bottle":  8080,
		"Pyramid": 6543,

		// Java frameworks
		"Spring Boot": 8080,
		"Spring MVC":  8080,
		"Quarkus":     8080,
		"Micronaut":   8080,
		"Vert.x":      8080,

		// Go frameworks
		"Gin":         8080,
		"Echo":        1323,
		"Fiber":       3000,
		"Gorilla Mux": 8080,
		"Chi":         8080,

		// PHP frameworks
		"Laravel":     8000,
		"Symfony":     8000,
		"CodeIgniter": 8080,
		"CakePHP":     8765,
		"Yii":         8080,

		// Ruby frameworks
		"Ruby on Rails": 3000,
		"Sinatra":       4567,
		"Grape":         9292,
		"Hanami":        2300,

		// Deno frameworks
		"Fresh":    8000,
		"Oak":      8000,
		"Hono":     8000,
		"Aleph.js": 8080,
		"Ultra":    8000,

		// Rust frameworks
		"Actix Web": 8080,
		"Axum":      3000,
		"Warp":      3030,
		"Rocket":    8000,
		"Tide":      8080,
		"Hyper":     3000,

		// Static site generators
		"Jekyll":     4000,
		"Hugo":       1313,
		"Hexo":       4000,
		"VuePress":   8080,
		"Docusaurus": 3000,
		"GitBook":    4000,
	},
}

// PackageManagers package manager configuration
var PackageManagers = map[types.SupportedLanguage][]string{
	types.LanguageNode:   {"npm", "yarn", "pnpm", "bun"},
	types.LanguagePython: {"pip", "poetry", "pipenv", "conda"},
	types.LanguageJava:   {"maven", "gradle", "sbt"},
	types.LanguageGo:     {"go"},
	types.LanguagePHP:    {"composer"},
	types.LanguageRuby:   {"bundler", "gem"},

	types.LanguageDeno:       {"deno"},
	types.LanguageRust:       {"cargo"},
	types.LanguageStaticfile: {},
	types.LanguageShell:      {},
}

// BuildTools build tool configuration
var BuildTools = map[types.SupportedLanguage][]string{
	types.LanguageNode:   {"webpack", "vite", "rollup", "parcel", "esbuild", "turbo", "nx"},
	types.LanguagePython: {"poetry", "setuptools", "flit", "hatch", "pdm"},
	types.LanguageJava:   {"maven", "gradle", "sbt", "ant"},
	types.LanguageGo:     {"make", "task", "mage", "go"},
	types.LanguagePHP:    {"composer"},
	types.LanguageRuby:   {"bundler", "rake"},

	types.LanguageDeno:       {"deno"},
	types.LanguageRust:       {"cargo"},
	types.LanguageStaticfile: {"webpack", "vite", "gulp", "grunt"},
	types.LanguageShell:      {"make"},
}

// ScanConfig scan configuration
var ScanConfig = struct {
	DefaultDepth   int
	MaxFiles       int
	IgnoreDirs     []string
	IgnorePatterns []string
	ImportantFiles []string
}{
	DefaultDepth: 3,
	MaxFiles:     1000,
	IgnoreDirs: []string{
		"node_modules",
		".git",
		".svn",
		".hg",
		"vendor",
		"target",
		"build",
		"dist",
		"out",
		".next",
		".nuxt",
		".vuepress",
		"__pycache__",
		".pytest_cache",
		".coverage",
		"coverage",
		".nyc_output",
		"logs",
		"log",
		"tmp",
		"temp",
		".tmp",
		".temp",
		".cache",
		".DS_Store",
		"Thumbs.db",
	},
	IgnorePatterns: []string{
		"*.log",
		"*.tmp",
		"*.temp",
		"*.cache",
		"*.pid",
		"*.lock",
		"*.swp",
		"*.swo",
		"*~",
		".DS_Store",
		"Thumbs.db",
		"*.min.js",
		"*.min.css",
		"*.map",
	},
	ImportantFiles: []string{
		"package.json",
		"requirements.txt",
		"Pipfile",
		"pyproject.toml",
		"pom.xml",
		"build.gradle",
		"go.mod",
		"composer.json",
		"Gemfile",

		"deno.json",
		"deno.jsonc",
		"Cargo.toml",
		"Dockerfile",
		"docker-compose.yml",
		"docker-compose.yaml",
		".nvmrc",
		".node-version",
		".python-version",
		".ruby-version",
		".go-version",
		"runtime.txt",
		"Procfile",
		"app.json",
		"now.json",
		"vercel.json",
		"netlify.toml",
		"_config.yml",
		"gatsby-config.js",
		"next.config.js",
		"nuxt.config.js",
		"vue.config.js",
		"angular.json",
		"ember-cli-build.js",
		"svelte.config.js",
		"vite.config.js",
		"webpack.config.js",
		"rollup.config.js",
		"tsconfig.json",
		"jsconfig.json",
		"babel.config.js",
		".babelrc",
		"eslint.config.js",
		".eslintrc.js",
		".eslintrc.json",
		"prettier.config.js",
		".prettierrc",
		"tailwind.config.js",
		"postcss.config.js",
	},
}

// DetectionWeights detection weight configuration
var DetectionWeights = struct {
	FileExists   map[string]int
	ContentMatch map[string]int
	Framework    map[string]int
}{
	FileExists: map[string]int{
		"HIGH":   30, // Critical configuration files
		"MEDIUM": 20, // Important files
		"LOW":    10, // General files
	},
	ContentMatch: map[string]int{
		"HIGH":   25, // Strong match (e.g., dependency declarations)
		"MEDIUM": 15, // Medium match (e.g., import statements)
		"LOW":    5,  // Weak match (e.g., comments)
	},
	Framework: map[string]int{
		"EXPLICIT": 40, // Explicit framework declarations
		"IMPLICIT": 20, // Implicit framework features
		"WEAK":     10, // Weak framework features
	},
}

// ConfidenceThresholds confidence thresholds
var ConfidenceThresholds = struct {
	High    int
	Medium  int
	Low     int
	Minimum int
}{
	High:    80, // High confidence
	Medium:  60, // Medium confidence
	Low:     40, // Low confidence
	Minimum: 20, // Minimum confidence
}

// ErrorCodes error codes
const (
	// General errors
	ErrorUnknown      = "UNKNOWN_ERROR"
	ErrorInvalidInput = "INVALID_INPUT"

	// Git related errors
	ErrorGitCloneFailed = "GIT_CLONE_FAILED"
	ErrorGitInvalidURL  = "GIT_INVALID_URL"

	// File system errors
	ErrorFileNotFound      = "FILE_NOT_FOUND"
	ErrorDirectoryNotFound = "DIRECTORY_NOT_FOUND"
	ErrorPermissionDenied  = "PERMISSION_DENIED"

	// Detection related errors
	ErrorNoLanguageDetected  = "NO_LANGUAGE_DETECTED"
	ErrorNoValidDetection    = "NO_VALID_DETECTION"
	ErrorUnsupportedLanguage = "UNSUPPORTED_LANGUAGE"

	// Plan generation errors
	ErrorBaseImageNotFound    = "BASE_IMAGE_NOT_FOUND"
	ErrorPlanGenerationFailed = "PLAN_GENERATION_FAILED"

	// Output related errors
	ErrorOutputFormatInvalid = "OUTPUT_FORMAT_INVALID"
	ErrorOutputWriteFailed   = "OUTPUT_WRITE_FAILED"
)

// CLIDefaults CLI default configuration
var CLIDefaults = struct {
	Depth               int
	MaxFiles            int
	OutputFormat        types.OutputFormat
	ConfidenceThreshold int
	Timeout             int // milliseconds
	TempDirPrefix       string
	LogLevel            string
}{
	Depth:               3,
	MaxFiles:            1000,
	OutputFormat:        types.OutputFormatJSON,
	ConfidenceThreshold: 20,
	Timeout:             30000, // 30 seconds
	TempDirPrefix:       "devbox-pack-",
	LogLevel:            "info",
}

// DefaultVersions default version configuration
var DefaultVersions = map[types.SupportedLanguage]string{
	types.LanguageNode:   "20",
	types.LanguagePython: "3.11",
	types.LanguageJava:   "17",
	types.LanguageGo:     "1.21",
	types.LanguagePHP:    "8.2",
	types.LanguageRuby:   "3.2",

	types.LanguageDeno:       "1.40",
	types.LanguageRust:       "1.70",
	types.LanguageStaticfile: "1.0",
	types.LanguageShell:      "1.0",
}

// VersionPatterns version matching patterns
var VersionPatterns = struct {
	Semver        *regexp.Regexp
	MajorMinor    *regexp.Regexp
	MajorOnly     *regexp.Regexp
	NodeVersion   *regexp.Regexp
	PythonVersion *regexp.Regexp
	JavaVersion   *regexp.Regexp
	GoVersion     *regexp.Regexp
}{
	Semver:        regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*))?(?:\+([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*))?$`),
	MajorMinor:    regexp.MustCompile(`^(\d+)\.(\d+)$`),
	MajorOnly:     regexp.MustCompile(`^(\d+)$`),
	NodeVersion:   regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)$`),
	PythonVersion: regexp.MustCompile(`^(\d+)\.(\d+)(?:\.(\d+))?$`),
	JavaVersion:   regexp.MustCompile(`^(\d+)(?:\.(\d+))?(?:\.(\d+))?$`),
	GoVersion:     regexp.MustCompile(`^(\d+)\.(\d+)(?:\.(\d+))?$`),
}
