# Architecture Overview

## Introduction

The DevBox Pack system is a static analysis tool that generates execution plans for various programming languages and frameworks. It analyzes source code repositories to automatically detect the technology stack and generate appropriate build, development, and production commands along with environment configurations.

## Core Components

### 1. Detection Engine

The detection engine is responsible for analyzing source code and identifying the appropriate technology provider. It operates through a priority-based system where multiple providers can be evaluated, and the one with the highest confidence score is selected.

**Key Features:**
- Multi-provider detection with confidence scoring
- Priority-based provider ordering
- Static analysis only (no network requests during detection)
- Support for monorepos and composite projects

### 2. Provider System

Providers are the core components that implement language and framework-specific detection and plan generation logic. Each provider implements the `Provider` interface with the following methods:

- `GetName()` - Returns the provider name
- `GetPriority()` - Returns priority value (lower = higher priority)
- `Detect()` - Analyzes source code and returns detection results

**Current Provider Priorities:**
1. Static File (Priority: 10) - Highest priority for static websites
2. Shell (Priority: 30) - Shell scripts and basic executables
3. PHP (Priority: 60) - PHP applications
4. Ruby (Priority: 65) - Ruby applications
5. Java (Priority: 70) - Java applications
6. Go (Priority: 75) - Go applications
7. Node.js (Priority: 80) - JavaScript/TypeScript applications
8. Python (Priority: 80) - Python applications
9. Rust (Priority: 85) - Rust applications
10. Deno (Priority: 95) - Deno applications

### 3. Plan Generation

Once a provider is selected, it generates an `ExecutionPlan` containing:

- **Base Configuration**: Container base image selection
- **Runtime Environment**: Language version, tools, environment variables
- **Dependencies**: APT packages and system requirements
- **Commands**: Development, build, and production start commands
- **Evidence**: Detection metadata and confidence indicators

## Detection Workflow

### Phase 1: Source Code Acquisition
1. Clone or access the target repository
2. Navigate to the specified subdirectory (if provided)
3. Scan for relevant files and directory structure

### Phase 2: Provider Detection
1. Initialize all available providers
2. Execute detection logic for each provider
3. Collect confidence scores and evidence
4. Select the provider with the highest confidence
5. Handle manual provider override if specified

### Phase 3: Plan Generation
1. Analyze project structure and dependencies
2. Determine appropriate base image and version
3. Resolve package managers and lock files
4. Generate development, build, and production commands
5. Identify required system dependencies
6. Compile evidence and metadata

### Phase 4: Output Generation
1. Format the execution plan according to specified output format
2. Include detection evidence and confidence metrics
3. Handle error cases and fallback scenarios

## Base Image Selection

The system uses a standardized base image naming convention:

- **Node.js**: `base:node-<version>` (e.g., `base:node-18`)
- **Python**: `base:python-<version>` (e.g., `base:python-3.11`)
- **Go**: `base:go-<version>` (e.g., `base:go-1.21`)
- **Java**: `base:java-<version>` (e.g., `base:java-17`)
- **PHP**: `base:php-<version>` (e.g., `base:php-8.2`)
- **Ruby**: `base:ruby-<version>` (e.g., `base:ruby-3.2`)
- **Deno**: `base:deno-<version>` (e.g., `base:deno-1.40`)
- **Rust**: `base:rust-<version>` (e.g., `base:rust-1.75`)
- **Static**: `base:caddy` for static file serving
- **Shell**: `base:debian` for shell scripts

## Version Resolution Strategy

The system follows a consistent version resolution order across all providers:

1. **Explicit Configuration**: Version specified in language-specific files (e.g., `.node-version`, `runtime.txt`)
2. **Lock Files**: Version constraints from dependency lock files
3. **Package Manager Config**: Version from package manager configuration
4. **Default Fallback**: Provider-specific default version

## Command Generation

### Development Commands
- Bind to `0.0.0.0` for container compatibility
- Use framework-specific development servers when available
- Include hot-reload and debugging capabilities
- Set appropriate environment variables (e.g., `NODE_ENV=development`)

### Build Commands
- Use standard build tools and scripts
- Handle framework-specific build processes
- Generate optimized production assets
- Support custom build configurations

### Production Commands
- Bind to `0.0.0.0:${PORT}` for container deployment
- Use production-optimized settings
- Handle static file serving with Caddy when appropriate
- Include health check endpoints when available

## Network and Port Handling

All generated commands follow these networking conventions:

- **Host Binding**: Always use `0.0.0.0` instead of `localhost` or `127.0.0.1`
- **Port Configuration**: Use `${PORT}` environment variable with sensible defaults
- **Static Sites**: Serve through Caddy on port 80 with SPA support
- **API Services**: Bind to the specified port with proper CORS configuration

## Error Handling and Fallbacks

### Detection Failures
- If no provider matches with sufficient confidence, fall back to Shell provider
- Log detection attempts and confidence scores for debugging
- Provide clear error messages for unsupported configurations

### Version Resolution Failures
- Fall back to provider default versions
- Warn about potential compatibility issues
- Support manual version override through CLI parameters

### Command Generation Failures
- Provide basic fallback commands for common scenarios
- Include manual configuration hints in error messages
- Support custom command override through configuration

## Extensibility

The architecture supports easy extension through:

- **Custom Providers**: Implement the `Provider` interface for new languages/frameworks
- **Plugin System**: Future support for external provider plugins
- **Configuration Override**: CLI and configuration file support for customization
- **Rule Engine**: Planned support for custom detection and generation rules

## Integration Points

### CLI Interface
- Command-line tool for direct execution plan generation
- Support for various output formats (JSON, YAML, etc.)
- Integration with CI/CD pipelines and development workflows

### Library Usage
- Go package for programmatic access
- API for integration with other tools and services
- Structured output for automated processing

### Container Runtime
- Compatible with Docker and container orchestration platforms
- Standardized base image requirements
- Environment variable and port configuration support