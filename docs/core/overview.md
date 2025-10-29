# DevBox Pack - Technical Architecture Overview

DevBox Pack is a Go-based static analysis tool that implements a priority-based provider system for automatic project detection and execution plan generation. The system uses confidence scoring algorithms and weighted indicators to accurately identify programming languages, frameworks, and generate containerized deployment configurations.

## Core Architecture Components

### 1. Detection Engine (`pkg/detector/engine.go`)

The `DetectionEngine` coordinates all language providers through a priority-based detection system:

```go
type DetectionEngine struct {
    providers map[string]Provider
}
```

**Provider Priority Order** (lower number = higher priority):
- Static Files: Priority 200
- Node.js: Priority 50  
- Python: Priority 60
- Java: Priority 70
- Go: Priority 75
- PHP: Priority 80
- Ruby: Priority 90
- Deno: Priority 100
- Rust: Priority 110
- Shell: Priority 150

### 2. Provider Interface (`pkg/types/types.go`)

All language providers implement the standardized `Provider` interface:

```go
type Provider interface {
    GetName() string
    GetPriority() int
    Detect(projectPath string, files []FileInfo, gitHandler interface{}) (*DetectResult, error)
}
```

### 3. Confidence Scoring System

Each provider uses weighted indicators to calculate detection confidence:

```go
type ConfidenceIndicator struct {
    Weight    int  `json:"weight"`
    Satisfied bool `json:"satisfied"`
}
```

**Example: Go Provider Detection Logic**
```go
indicators := []types.ConfidenceIndicator{
    {Weight: 40, Satisfied: p.HasFile(files, "go.mod")},
    {Weight: 30, Satisfied: p.HasFile(files, "go.work")},
    {Weight: 25, Satisfied: p.HasAnyFile(files, []string{"*.go"})},
    {Weight: 15, Satisfied: p.HasFile(files, "go.sum")},
    {Weight: 10, Satisfied: p.HasAnyFile(files, []string{"main.go", "cmd/"})},
    {Weight: 5, Satisfied: p.HasFile(files, "vendor/")},
    {Weight: 5, Satisfied: p.HasAnyFile(files, []string{"Makefile", "makefile"})},
}
```

**Detection Thresholds:**
- Go: 0.2 (20% confidence)
- Python: 0.3 (30% confidence)  
- Java: 0.2 (20% confidence)
- PHP: 0.2 (20% confidence)
- Deno: 0.2 (20% confidence)

## Data Structures

### ExecutionPlan (`pkg/types/types.go`)

The core output structure containing complete deployment configuration:

```go
type ExecutionPlan struct {
    Provider string        `json:"provider"`
    Base     BaseConfig    `json:"base"`
    Runtime  RuntimeConfig `json:"runtime"`
    Apt      []string      `json:"apt,omitempty"`
    Commands Commands      `json:"commands,omitempty"`
    Port     int           `json:"port"`
    Evidence Evidence      `json:"evidence,omitempty"`
}

type Commands struct {
    Dev   []string `json:"dev,omitempty"`
    Build []string `json:"build,omitempty"`
    Start []string `json:"start,omitempty"`
}
```

### DetectResult (`pkg/types/types.go`)

Provider detection output with confidence metrics:

```go
type DetectResult struct {
    Matched        bool                   `json:"matched"`
    Provider       *string                `json:"provider"`
    Confidence     float64                `json:"confidence"`
    Evidence       Evidence               `json:"evidence"`
    Language       string                 `json:"language"`
    Framework      string                 `json:"framework"`
    Version        string                 `json:"version"`
    PackageManager *PackageManager        `json:"packageManager"`
    BuildTools     []string               `json:"buildTools"`
    Metadata       map[string]interface{} `json:"metadata"`
}
```

## CLI Implementation (`pkg/cli/cli.go`)

### Command Structure

```bash
devbox-pack <repository> [options]
```

### Available Options

```go
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
}
```

### Supported Output Formats
- `pretty` - Human-readable format (default)
- `json` - Machine-readable JSON format

## Language-Specific Provider Details

### **[Go Provider](../providers/golang.md)** (`pkg/providers/golang.go`)
- **Detection Files**: `go.mod` (40%), `go.work` (30%), `*.go` (25%), `go.sum` (15%)
- **Version Sources**: `go.work` → `go.mod` → `.go-version` → default (1.21)
- **Framework Detection**: Gin, Echo, Fiber, Gorilla Mux, Beego, Revel, Cobra CLI

### **[Python Provider](../providers/python.md)** (`pkg/providers/python.go`)
- **Detection Files**: `requirements.txt|pyproject.toml|setup.py|Pipfile` (30%), `*.py` (25%)
- **Package Managers**: Poetry, Pipenv, pip (with lock file optimization)
- **Framework Detection**: Django, Flask, FastAPI, Streamlit

### **[Java Provider](../providers/java.md)** (`pkg/providers/java.go`)
- **Detection Files**: `pom.xml|build.gradle|build.gradle.kts` (30%), `*.java|*.kt|*.scala` (25%)
- **Build Tools**: Maven, Gradle (with wrapper detection)
- **Framework Detection**: Spring Boot, Quarkus, Micronaut

### **[Node.js Provider](../providers/node.md)** (`pkg/providers/node.go`)
- **Detection Files**: `package.json` (40%), `*.js|*.ts` (25%), lock files (15%)
- **Package Managers**: npm, yarn, pnpm (with corepack support)
- **Framework Detection**: Next.js, React, Vue, Express, Nuxt

## Service Layer (`pkg/service/devbox.go`)

The main service orchestrates the detection and plan generation process:

```go
type DevBoxPack struct {
    gitHandler      *git.GitHandler
    detectionEngine *detector.DetectionEngine
    planGenerator   *generators.ExecutionPlanGenerator
    outputUtils     *formatters.OutputUtils
}
```

**Core Workflow:**
1. Repository cloning/local path validation
2. File system scanning with configurable depth
3. Provider detection with confidence scoring
4. Execution plan generation
5. Output formatting (JSON/Pretty)

## Contributing

When adding new providers or modifying existing ones:

1. Follow the provider interface defined in [API Schema](api-schema.md)
2. Add comprehensive tests as outlined in [Testing Guide](testing.md)
3. Update relevant documentation and examples
4. Ensure detection priority and confidence scoring are appropriate

For more information, see the main project [CONTRIBUTING.md](../CONTRIBUTING.md).