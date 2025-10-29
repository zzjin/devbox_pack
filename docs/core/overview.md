# DevBox Pack - Project Overview

DevBox Pack is an intelligent CLI tool that automatically detects project languages and frameworks to generate optimized execution plans for containerized deployments. This document provides a comprehensive overview of the system architecture, detection mechanisms, and provider ecosystem.

## System Overview

DevBox Pack analyzes codebases to automatically determine:
- Programming languages and runtime versions
- Framework dependencies and configurations  
- Build processes and optimization strategies
- Deployment and startup commands

The tool generates comprehensive execution plans that can be used for CI/CD pipelines, Docker containerization, and cloud deployment automation.

## Provider System Architecture

The core detection system is built around a provider-based architecture located in `/railpack/core/providers`. Each provider handles language-specific detection, dependency management, and execution plan generation.

### Detection System

- **Detection Priority** (first match wins): PHP → Go → Java → Rust → Ruby → Elixir → Python → Deno → Node → Staticfile → Shell (see `provider.go:GetLanguageProviders`)
- **Common Workflow**:
  - **Detect**: Determine if matched based on characteristic files/configurations
  - **Initialize**: Read necessary configurations (such as `package.json`, workspace, etc.)
  - **Plan**:
    - Install language versions and tools (via `mise` specifying version sources: environment variables, version files, configurations, etc.)
    - Install dependencies (various ecosystem package managers, optimized with lock files)
    - Optional Prune (remove dev dependencies, etc., controlled by environment variables)
    - Build: Execute framework/tool-specific build commands and set cache directories
    - StartCmd: Generate startup commands based on scripts/entry files/framework rules

### Language-Specific Providers

- **[Node.js](../providers/node.md)** - JavaScript/TypeScript projects, frameworks, and package managers
- **[Python](../providers/python.md)** - Python projects, frameworks, and dependency management
- **[Java](../providers/java.md)** - Java applications, Maven, Gradle, and Spring Boot
- **[Go](../providers/golang.md)** - Go modules, workspaces, and build configurations
- **[PHP](../providers/php.md)** - PHP applications, Composer, and Laravel framework
- **[Ruby](../providers/ruby.md)** - Ruby applications, Bundler, and Rails framework
- **[Deno](../providers/deno.md)** - Deno runtime and TypeScript applications
- **[Rust](../providers/rust.md)** - Rust applications, Cargo, and web frameworks
- **[Static Files](../providers/staticfile.md)** - Static HTML/CSS/JS sites and SPAs
- **[Shell](../providers/shell.md)** - Shell scripts and custom deployment configurations

## Key Features

- **🔍 Intelligent Detection**: Automatically identifies programming languages, frameworks, and build tools
- **📋 Execution Plans**: Generates comprehensive plans with install, build, and startup commands
- **🐳 Container Optimized**: Designed for containerized environments with efficient caching strategies
- **🛠️ Multi-Language Support**: Covers 10+ programming languages and their ecosystems
- **⚡ Performance Focused**: Minimal overhead with intelligent dependency management

## Integration and Usage

DevBox Pack integrates seamlessly with:
- **CI/CD Pipelines**: Generate build and deployment configurations
- **Container Platforms**: Optimize Docker builds and runtime configurations  
- **Cloud Deployment**: Automate application deployment across cloud providers
- **Development Workflows**: Standardize local development environments

For detailed usage instructions and examples, refer to the [CLI Usage Guide](cli-usage.md) and [Examples](examples.md).

## Contributing

When adding new providers or modifying existing ones:

1. Follow the provider interface defined in [API Schema](api-schema.md)
2. Add comprehensive tests as outlined in [Testing Guide](testing.md)
3. Update relevant documentation and examples
4. Ensure detection priority and confidence scoring are appropriate

For more information, see the main project [CONTRIBUTING.md](../CONTRIBUTING.md).