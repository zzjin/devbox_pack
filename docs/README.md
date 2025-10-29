# DevBox Pack - Go Implementation Documentation

DevBox Pack is a high-performance Go-based static analysis tool that implements a sophisticated provider system for automatic project detection and execution plan generation. Built with a priority-based detection engine and confidence scoring algorithms, it accurately identifies programming languages, frameworks, and generates optimized containerized deployment configurations.

## 🏗️ Technical Architecture

### Core Components (`pkg/`)

```
pkg/
├── cli/           # CLI application with argument parsing
├── detector/      # Detection engine with provider coordination
├── generators/    # Execution plan generation logic
├── providers/     # Language-specific detection providers
├── service/       # Main orchestration service
├── types/         # Type definitions and interfaces
└── utils/         # Utility functions and helpers
```

### Detection Engine Implementation

**Priority-Based Provider System:**
```go
// Provider priorities (lower = higher priority)
Node.js:     50    // JavaScript/TypeScript projects
Python:      60    // Python applications  
Java:        70    // JVM-based applications
Go:          75    // Go modules and workspaces
PHP:         80    // PHP web applications
Ruby:        90    // Ruby applications
Deno:       100    // Deno runtime projects
Rust:       110    // Rust applications
Shell:      150    // Shell scripts
Static:     200    // Static file serving
```

**Confidence Scoring Algorithm:**
```go
type ConfidenceIndicator struct {
    Weight    int  `json:"weight"`    // Indicator importance (0-100)
    Satisfied bool `json:"satisfied"` // Whether indicator is present
}

// Confidence = Σ(satisfied_weight) / Σ(total_weight)
```

## 📋 Data Structures & API

### ExecutionPlan Output
```go
type ExecutionPlan struct {
    Provider string        `json:"provider"`     // Detected provider name
    Base     BaseConfig    `json:"base"`         // Container base configuration
    Runtime  RuntimeConfig `json:"runtime"`      // Runtime environment setup
    Apt      []string      `json:"apt"`          // System packages to install
    Commands Commands      `json:"commands"`     // Build/dev/start commands
    Port     int           `json:"port"`         // Application port
    Evidence Evidence      `json:"evidence"`     // Detection evidence
}
```

### Provider Detection Results
```go
type DetectResult struct {
    Matched        bool                   `json:"matched"`        // Detection success
    Provider       *string                `json:"provider"`       // Provider name
    Confidence     float64                `json:"confidence"`     // Score (0.0-1.0)
    Evidence       Evidence               `json:"evidence"`       // Detection evidence
    Language       string                 `json:"language"`       // Primary language
    Framework      string                 `json:"framework"`      // Detected framework
    Version        string                 `json:"version"`        // Language version
    PackageManager *PackageManager        `json:"packageManager"` // Package manager info
    BuildTools     []string               `json:"buildTools"`     // Build tools detected
    Metadata       map[string]interface{} `json:"metadata"`       // Provider-specific data
}
```

## 🔧 [Core Documentation](./core/)
Technical implementation details and system architecture:

- **[Technical Overview](./core/overview.md)** - Detailed architecture, data structures, and provider system
- **[Architecture Details](./core/architecture.md)** - System design patterns and component interactions
- **[API Schema](./core/api-schema.md)** - Complete type definitions and data structure specifications
- **[CLI Usage](./core/cli-usage.md)** - Command-line interface implementation and options
- **[Testing Guide](./core/testing.md)** - Test strategy, unit tests, and integration testing
- **[Examples](./core/examples.md)** - Real-world usage examples with actual output

## 🚀 [Provider Documentation](./providers/)
Language-specific detection logic and implementation details:

- **[Go Provider](./providers/golang.md)** - Go modules, workspaces, framework detection (Gin, Echo, Fiber)
- **[Node.js Provider](./providers/node.md)** - npm/yarn/pnpm, framework detection (Next.js, React, Vue)
- **[Python Provider](./providers/python.md)** - pip/poetry/pipenv, framework detection (Django, Flask, FastAPI)
- **[Java Provider](./providers/java.md)** - Maven/Gradle, framework detection (Spring Boot, Quarkus)
- **[PHP Provider](./providers/php.md)** - Composer, framework detection (Laravel, Symfony)
- **[Ruby Provider](./providers/ruby.md)** - Bundler, framework detection (Rails, Sinatra)
- **[Rust Provider](./providers/rust.md)** - Cargo, framework detection (Rocket, Actix)
- **[Deno Provider](./providers/deno.md)** - Deno runtime, TypeScript support
- **[Shell Provider](./providers/shell.md)** - Shell script detection and execution
- **[Static Files Provider](./providers/staticfile.md)** - Static HTML/CSS/JS serving

## 🚀 CLI Usage Examples

### Basic Usage
```bash
# Analyze a Git repository
devbox-pack https://github.com/user/my-app

# Analyze local project
devbox-pack . --offline

# Force specific provider
devbox-pack . --provider go --format json

# Analyze specific branch/subdirectory
devbox-pack https://github.com/user/mono-repo --ref develop --subdir backend
```

### CLI Options (`pkg/cli/cli.go`)
```bash
devbox-pack <repository> [options]

Arguments:
  repository               Git repository URL or local path

Options:
  --ref <ref>             Git branch/tag (default: main)
  --subdir <path>         Subdirectory path within repository
  --provider <name>       Force specific provider (node|python|java|go|php|ruby|deno|rust|shell|staticfile)
  --format <format>       Output format (pretty|json, default: pretty)
  --verbose               Enable detailed logging and detection info
  --offline               Skip Git operations, analyze local path only
  --platform <arch>       Target platform (e.g., linux/amd64)
  --base <name>           Override base image selection
  -h, --help              Show help information
  -v, --version           Show version information
```

### Output Examples

**Pretty Format (Default):**
```
✓ Detected Go project
  Language: go (v1.21)
  Framework: gin
  Package Manager: go modules
  
Commands:
  Dev:   go mod download && go run .
  Build: go mod download && go build -o app .
  Start: ./app
  
Port: 8080
Base: golang:1.21-alpine
```

**JSON Format:**
```json
{
  "provider": "go",
  "base": {"name": "golang:1.21-alpine"},
  "runtime": {
    "language": "go",
    "version": "1.21",
    "tools": ["go"]
  },
  "commands": {
    "dev": ["go mod download", "go run ."],
    "build": ["go mod download", "go build -o app ."],
    "start": ["./app"]
  },
  "port": 8080,
  "evidence": {
    "files": ["go.mod", "main.go"],
    "reason": "Found go.mod and Go source files"
  }
}
```

## 🎯 Core Features & Implementation

### **Intelligent Detection System**
- **Static Analysis Only**: No network requests during detection phase
- **Confidence-Based Scoring**: Weighted indicators with configurable thresholds
- **Priority-Based Providers**: Deterministic detection order with override capability
- **Multi-Project Support**: Monorepo and workspace detection

### **Execution Plan Generation**
- **Container-Optimized**: Efficient layer caching and build strategies
- **Framework-Aware**: Automatic detection of 50+ frameworks across 10 languages
- **Version Resolution**: Smart version detection from multiple sources
- **Command Generation**: Dev, build, and production command generation

### **Performance & Reliability**
- **Native Go Implementation**: High-performance static analysis
- **Minimal Dependencies**: Self-contained binary with no external requirements
- **Error Handling**: Comprehensive error reporting with actionable messages
- **Test Coverage**: Extensive unit and integration test suite

## 📖 Documentation Navigation Tips

- Use the **Core Documentation** for understanding the system architecture and general usage
- Refer to **Provider Documentation** for language-specific implementation details
- Each provider document includes detection logic, supported tools, and configuration examples
- Cross-references between documents use relative paths for easy navigation

## 🤝 Contributing

When contributing to the documentation:

1. Core project changes should update files in `./core/`
2. Language-specific changes should update files in `./providers/`
3. Maintain cross-references when moving or renaming files
4. Follow the established documentation structure and formatting

---

For the most up-to-date information, please refer to the individual documentation files in their respective directories.