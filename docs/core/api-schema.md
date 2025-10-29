# API Schema Documentation

## Overview

This document describes the data structures and API specifications used by the DevBox Pack system for generating execution plans and detecting project configurations.

## Core Data Structures

### ExecutionPlan

The `ExecutionPlan` represents the complete execution configuration for a project after detection and analysis.

```go
type ExecutionPlan struct {
    // Provider name that was matched
    Provider string `json:"provider"`
    
    // Base environment configuration
    Base BaseConfig `json:"base"`
    
    // Runtime configuration
    Runtime RuntimeConfig `json:"runtime"`
    
    // System packages to install (optional)
    Apt []string `json:"apt,omitempty"`
    
    // Commands configuration (optional)
    Commands Commands `json:"commands,omitempty"`
    
    // Port configuration
    Port int `json:"port"`
    
    // Detection evidence (optional)
    Evidence Evidence `json:"evidence,omitempty"`
}
```

**Fields:**
- `provider`: String identifier of the matched provider (e.g., "node", "python", "go")
- `base`: Base container image configuration
- `runtime`: Language and runtime environment settings
- `apt`: Array of system packages to install via APT (only included if needed)
- `commands`: Development, build, and production commands (only included if available)
- `port`: Default port number for the application
- `evidence`: Detection metadata and reasoning (only included if available)

### BaseConfig

Defines the base container image configuration.

```go
type BaseConfig struct {
    // Base image name (e.g., "base:node-18", "base:python-3.11")
    Name string `json:"name"`
    
    // Target platform/architecture (optional)
    Platform string `json:"platform,omitempty"`
}
```

**Examples:**
- `base:node-18` - Node.js version 18
- `base:python-3.11` - Python version 3.11
- `base:go-1.21` - Go version 1.21
- `base:caddy` - Caddy web server for static files

### RuntimeConfig

Specifies the runtime environment configuration.

```go
type RuntimeConfig struct {
    // Programming language
    Language string `json:"language"`
    
    // Language version (optional)
    Version *string `json:"version,omitempty"`
    
    // Required tools (optional)
    Tools []string `json:"tools,omitempty"`
    
    // Environment variables (optional)
    Env map[string]string `json:"env,omitempty"`
}
```

**Language Values:**
- `node` - Node.js/JavaScript/TypeScript
- `python` - Python
- `go` - Go/Golang
- `java` - Java
- `php` - PHP
- `ruby` - Ruby
- `deno` - Deno
- `rust` - Rust
- `staticfile` - Static files (HTML/CSS/JS)
- `shell` - Shell scripts

**Common Tools:**
- `corepack` - Node.js package manager enabler
- `pnpm` - Fast Node.js package manager
- `poetry` - Python dependency management
- `uv` - Fast Python package installer
- `caddy` - Web server for static files

### Commands

Defines the available commands for different lifecycle phases.

```go
type Commands struct {
    // Development commands
    Dev []string `json:"dev,omitempty"`
    
    // Build commands
    Build []string `json:"build,omitempty"`
    
    // Production start commands
    Start []string `json:"start,omitempty"`
}
```

**Command Conventions:**
- All commands bind to `0.0.0.0` for container compatibility
- Use `${PORT}` environment variable for port configuration
- Include framework-specific optimizations and flags

### Evidence

Contains detection metadata and reasoning.

```go
type Evidence struct {
    // Key files that influenced detection
    Files []string `json:"files,omitempty"`
    
    // Human-readable reason for the match
    Reason string `json:"reason,omitempty"`
}
```

**Example Files:**
- `package.json` - Node.js project configuration
- `requirements.txt` - Python dependencies
- `go.mod` - Go module definition
- `Cargo.toml` - Rust project configuration
- `composer.json` - PHP dependencies

## Detection Data Structures

### DetectResult

The `DetectResult` represents the outcome of project detection by a specific provider.

```go
type DetectResult struct {
    // Whether the project matches this provider
    Matched bool `json:"matched"`
    
    // Provider name if matched
    Provider *string `json:"provider"`
    
    // Confidence score (0.0-1.0)
    Confidence float64 `json:"confidence"`
    
    // Detection evidence
    Evidence Evidence `json:"evidence"`
    
    // Detected language
    Language string `json:"language"`
    
    // Detected framework (optional)
    Framework string `json:"framework"`
    
    // Detected version (optional)
    Version string `json:"version"`
    
    // Package manager information (optional)
    PackageManager *PackageManager `json:"packageManager"`
    
    // Build tools (optional)
    BuildTools []string `json:"buildTools"`
    
    // Additional metadata (optional)
    Metadata map[string]interface{} `json:"metadata"`
}
```

**Confidence Scoring:**
- `1.0` - Perfect match with strong indicators
- `0.8-0.9` - High confidence with multiple indicators
- `0.6-0.7` - Medium confidence with some indicators
- `0.4-0.5` - Low confidence with weak indicators
- `0.0-0.3` - Very low confidence or no match

### PackageManager

Describes the detected package manager.

```go
type PackageManager struct {
    // Package manager name
    Name string `json:"name"`
    
    // Version (optional)
    Version string `json:"version,omitempty"`
    
    // Configuration files (optional)
    ConfigFiles []string `json:"configFiles,omitempty"`
}
```

**Common Package Managers:**
- `npm` - Node.js default package manager
- `yarn` - Alternative Node.js package manager
- `pnpm` - Fast Node.js package manager
- `pip` - Python package installer
- `poetry` - Python dependency management
- `composer` - PHP dependency manager
- `bundler` - Ruby gem manager

## CLI Options

### CLIOptions

Configuration options for the CLI interface.

```go
type CLIOptions struct {
    // Git repository URL or local path
    Repo *string `json:"repo,omitempty"`
    
    // Git reference (branch, tag, commit)
    Ref *string `json:"ref,omitempty"`
    
    // Subdirectory within the repository
    SubDir *string `json:"subdir,omitempty"`
    
    // Force specific provider
    Provider *string `json:"provider,omitempty"`
    
    // Output format (json, yaml, etc.)
    Format string `json:"format"`
    
    // Enable verbose output
    Verbose bool `json:"verbose"`
    
    // Offline mode (no network requests)
    Offline bool `json:"offline"`
    
    // Target platform
    Platform *string `json:"platform,omitempty"`
    
    // Force specific base image
    Base *string `json:"base,omitempty"`
}
```

## API Response Formats

### Success Response

```json
{
  "provider": "node",
  "base": {
    "name": "base:node-18",
    "platform": "linux/amd64"
  },
  "runtime": {
    "language": "node",
    "version": "18.19.0",
    "tools": ["corepack"],
    "env": {
      "NODE_ENV": "production"
    }
  },
  "apt": ["build-essential"],
  "commands": {
    "dev": ["npm run dev -- --host 0.0.0.0 --port ${PORT}"],
    "build": ["npm run build"],
    "start": ["npm start -- --host 0.0.0.0 --port ${PORT}"]
  },
  "port": 3000,
  "evidence": {
    "files": ["package.json", "package-lock.json"],
    "reason": "Detected Node.js project with package.json and npm lockfile"
  }
}
```

### Error Response

```json
{
  "error": "no supported language or framework detected",
  "details": "Scanned 45 files but found no recognizable project structure",
  "suggestions": [
    "Ensure the repository contains valid project files",
    "Check if the subdirectory path is correct",
    "Try specifying a provider manually with --provider"
  ]
}
```

## Version Information

### VersionInfo

Represents version detection results.

```go
type VersionInfo struct {
    // Detected version string
    Version string `json:"version"`
    
    // Version source (e.g., ".node-version", "package.json")
    Source string `json:"source"`
    
    // Whether this is a constraint or exact version
    IsConstraint bool `json:"isConstraint"`
}
```

## File Information

### FileInfo

Represents scanned file metadata.

```go
type FileInfo struct {
    // Relative path from project root
    Path string `json:"path"`
    
    // File size in bytes
    Size int64 `json:"size"`
    
    // Whether this is a directory
    IsDir bool `json:"isDir"`
    
    // File modification time
    ModTime time.Time `json:"modTime"`
}
```

## Scan Options

### ScanOptions

Configuration for project file scanning.

```go
type ScanOptions struct {
    // Maximum directory depth to scan
    MaxDepth int `json:"maxDepth"`
    
    // Maximum number of files to scan
    MaxFiles int `json:"maxFiles"`
    
    // File patterns to exclude
    ExcludePatterns []string `json:"excludePatterns,omitempty"`
    
    // File patterns to include
    IncludePatterns []string `json:"includePatterns,omitempty"`
}
```

**Default Values:**
- `MaxDepth`: 3 levels deep
- `MaxFiles`: 1000 files maximum
- Common exclusions: `node_modules`, `.git`, `vendor`, `target`

## Provider Interface

### Provider

Interface that all language providers must implement.

```go
type Provider interface {
    // Get provider name
    GetName() string
    
    // Get provider priority (lower = higher priority)
    GetPriority() int
    
    // Detect if this provider matches the project
    Detect(projectPath string, files []FileInfo) (*DetectResult, error)
    
    // Generate commands for the detected project
    GenerateCommands(result *DetectResult, options CLIOptions) Commands
    
    // Check if native compilation is needed
    NeedsNativeCompilation(result *DetectResult) bool
}
```

This interface ensures consistent behavior across all language and framework providers while allowing for provider-specific detection logic and command generation.