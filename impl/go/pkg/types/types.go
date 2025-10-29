package types

import (
	"fmt"
	"time"
)

// ExecutionPlan represents the complete execution configuration for a project
type ExecutionPlan struct {
	// Provider name that was matched
	Provider string `json:"provider"`

	// Base environment configuration
	Base BaseConfig `json:"base"`

	// Runtime configuration
	Runtime RuntimeConfig `json:"runtime"`

	// System packages to install
	Apt []string `json:"apt,omitempty"`

	// Commands configuration
	Commands Commands `json:"commands,omitempty"`

	// Port configuration
	Port int `json:"port"`

	// Detection evidence
	Evidence Evidence `json:"evidence,omitempty"`
}

// Base environment configuration
type BaseConfig struct {
	// Base image name
	Name string `json:"name"`
	// Platform specification
	Platform string `json:"platform,omitempty"`
}

// RuntimeConfig represents the runtime configuration
type RuntimeConfig struct {
	// Language type
	Language string `json:"language"`
	// Language version
	Version *string `json:"version,omitempty"`
	// Required tools, e.g., corepack, pnpm, poetry, uv, caddy
	Tools []string `json:"tools,omitempty"`
	// Environment variables for runtime
	Environment map[string]string `json:"environment,omitempty"`
}

// Evidence represents detection evidence
type Evidence struct {
	// Key files, e.g., package.json, lockfiles, etc.
	Files []string `json:"files,omitempty"`
	// Reason for match
	Reason string `json:"reason,omitempty"`
}

// DetectResult represents the result of project detection
type DetectResult struct {
	// Whether the project matches this provider
	Matched bool `json:"matched"`
	// Provider name if matched
	Provider *string `json:"provider"`
	// Confidence score (0-1)
	Confidence float64 `json:"confidence"`
	// Detection evidence
	Evidence Evidence `json:"evidence"`
	// Detected language
	Language string `json:"language"`
	// Detected framework
	Framework string `json:"framework"`
	// Detected version
	Version string `json:"version"`
	// Package manager information
	PackageManager *PackageManager `json:"packageManager"`
	// Build tools
	BuildTools []string `json:"buildTools"`
	// Additional metadata
	Metadata map[string]interface{} `json:"metadata"`
}

// Provider interface defines the contract for all language providers
type Provider interface {
	// Get provider name
	GetName() string
	// Get provider priority (lower number = higher priority)
	GetPriority() int
	// Detect if project uses this provider
	Detect(projectPath string, files []FileInfo, gitHandler interface{}) (*DetectResult, error)
}

// CLIOptions represents command line interface options
type CLIOptions struct {
	Repository string  `json:"repository"`
	Ref        *string `json:"ref,omitempty"`
	Subdir     *string `json:"subdir,omitempty"`
	Provider   *string `json:"provider,omitempty"`
	Format     string  `json:"format"`
	Verbose    bool    `json:"verbose"`
	Offline    bool    `json:"offline"`
	Platform   *string `json:"platform,omitempty"`
	Base       *string `json:"base,omitempty"`
	Help       bool    `json:"help,omitempty"`
	Version    bool    `json:"version,omitempty"`
	Quiet      bool    `json:"quiet,omitempty"`
	Pretty     bool    `json:"pretty,omitempty"`
}

// GitRepository represents a Git repository
type GitRepository struct {
	// Repository URL
	URL string `json:"url"`
	// Git reference (branch, tag, commit)
	Ref *string `json:"ref,omitempty"`
	// Subdirectory within the repository
	Subdir *string `json:"subdir,omitempty"`
	// Whether it's a local repository
	IsLocal bool `json:"isLocal,omitempty"`
}

// FileInfo represents file information
type FileInfo struct {
	// File path
	Path string `json:"path"`
	// File name
	Name string `json:"name"`
	// Whether it's a directory
	IsDirectory bool `json:"isDirectory"`
	// File size
	Size *int64 `json:"size,omitempty"`
	// File extension
	Extension *string `json:"extension,omitempty"`
	// Modification time
	MTime *time.Time `json:"mtime,omitempty"`
}

// VersionInfo represents version information
type VersionInfo struct {
	// Version string
	Version string `json:"version"`
	// Source of version information
	Source string `json:"source"` // 'file' | 'env' | 'default'
	// Source detail
	SourceDetail *string `json:"sourceDetail,omitempty"`
}

// PackageManager represents a package manager
type PackageManager struct {
	// Package manager name
	Name string `json:"name"`
	// Lock file name
	LockFile *string `json:"lockFile,omitempty"`
	// Configuration file name
	ConfigFile *string `json:"configFile,omitempty"`
	// Whether to use corepack
	UseCorepack bool `json:"useCorepack,omitempty"`
}

// BaseCatalog represents the base image catalog
type BaseCatalog map[string]map[string]string

// DevBoxPackError represents a custom error for DevBox Pack
type DevBoxPackError struct {
	Message string
	Code    string
	Details interface{}
}

// ConfidenceIndicator confidence indicator
type ConfidenceIndicator struct {
	Weight    int  `json:"weight"`
	Satisfied bool `json:"satisfied"`
}

// Commands represents the command configuration
type Commands struct {
	Dev   []string `json:"dev,omitempty"`
	Build []string `json:"build,omitempty"`
	Start []string `json:"start,omitempty"`
}

// ScanOptions represents scanning configuration
type ScanOptions struct {
	Depth    int `json:"depth"`
	MaxDepth int `json:"maxDepth"`
	MaxFiles int `json:"maxFiles"`
}

// Error codes
const (
	ErrorCodeGitError          = "GIT_ERROR"
	ErrorCodeLocalAccessError  = "LOCAL_ACCESS_ERROR"
	ErrorCodeInvalidPath       = "INVALID_PATH"
	ErrorCodeCloneError        = "CLONE_ERROR"
	ErrorCodeGitCheckoutError  = "GIT_CHECKOUT_ERROR"
	ErrorCodeSubdirAccessError = "SUBDIR_ACCESS_ERROR"
	ErrorCodeSubdirNotFound    = "SUBDIR_NOT_FOUND"
	ErrorCodeTempDirError      = "TEMP_DIR_ERROR"
	ErrorCodeFileReadError     = "FILE_READ_ERROR"
	ErrorCodeJSONParseError    = "JSON_PARSE_ERROR"
	ErrorCodeInvalidFormat     = "INVALID_FORMAT"
	ErrorCodeInvalidPlatform   = "INVALID_PLATFORM"
	ErrorCodeInvalidGitURL     = "INVALID_GIT_URL"
	ErrorCodeInvalidInput      = "INVALID_INPUT"
	ErrorCodeScanError         = "SCAN_ERROR"
	ErrorCodeInvalidProvider   = "INVALID_PROVIDER"
	ErrorCodeInvalidArgument   = "INVALID_ARGUMENT"
)

func (e *DevBoxPackError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDevBoxPackError creates a new DevBoxPackError
func NewDevBoxPackError(message, code string, details interface{}) *DevBoxPackError {
	return &DevBoxPackError{
		Message: message,
		Code:    code,
		Details: details,
	}
}

// SupportedLanguage represents supported programming languages
type SupportedLanguage string

const (
	LanguageNode   SupportedLanguage = "node"
	LanguagePython SupportedLanguage = "python"
	LanguageJava   SupportedLanguage = "java"
	LanguageGo     SupportedLanguage = "go"
	LanguagePHP    SupportedLanguage = "php"
	LanguageRuby   SupportedLanguage = "ruby"

	LanguageDeno       SupportedLanguage = "deno"
	LanguageRust       SupportedLanguage = "rust"
	LanguageStaticfile SupportedLanguage = "staticfile"
	LanguageShell      SupportedLanguage = "shell"
)

// OutputFormat represents output format types
type OutputFormat string

const (
	OutputFormatJSON   OutputFormat = "json"
	OutputFormatPretty OutputFormat = "pretty"
)

// Platform represents supported platforms
type Platform string

const (
	PlatformLinuxAMD64  Platform = "linux/amd64"
	PlatformLinuxARM64  Platform = "linux/arm64"
	PlatformDarwinAMD64 Platform = "darwin/amd64"
	PlatformDarwinARM64 Platform = "darwin/arm64"
)
