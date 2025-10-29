# Devbox Pack Documentation

Welcome to the Devbox Pack documentation! This directory contains comprehensive documentation for the project, organized into two main sections for better navigation and maintainability.

## 📁 Documentation Structure

### 🔧 [Core Documentation](./core/)
Essential project documentation covering the fundamental aspects of Devbox Pack:

- **[Overview & Architecture](./core/overview.md)** - Project overview, architecture, and core concepts
- **[Architecture Details](./core/architecture.md)** - Detailed system architecture and design patterns
- **[API Schema](./core/api-schema.md)** - Complete API schema and data structure definitions
- **[CLI Usage](./core/cli-usage.md)** - Command-line interface documentation and usage examples
- **[Testing Guide](./core/testing.md)** - Testing strategy, structure, and best practices
- **[Examples](./core/examples.md)** - Comprehensive examples and use cases

### 🚀 [Provider Documentation](./providers/)
Language and framework-specific provider documentation:

- **[Deno](./providers/deno.md)** - Deno runtime detection and configuration
- **[Go](./providers/golang.md)** - Go language support and toolchain detection
- **[Java](./providers/java.md)** - Java ecosystem support (Maven, Gradle)
- **[Node.js](./providers/node.md)** - Node.js and npm/yarn/pnpm support
- **[PHP](./providers/php.md)** - PHP and Composer support
- **[Python](./providers/python.md)** - Python and pip/poetry/pipenv support
- **[Ruby](./providers/ruby.md)** - Ruby and gem/bundler support
- **[Rust](./providers/rust.md)** - Rust and Cargo support
- **[Shell](./providers/shell.md)** - Shell script detection and execution
- **[Static Files](./providers/staticfile.md)** - Static file serving configuration

## 🚀 Quick Start

1. **New to Devbox Pack?** Start with the [Core Documentation](./core/overview.md)
2. **Looking for specific language support?** Check the [Provider Documentation](./providers/)
3. **Want to see examples?** Browse the [Examples](./core/examples.md)
4. **Need CLI help?** Refer to [CLI Usage](./core/cli-usage.md)
5. **Contributing?** See the [Testing Guide](./core/testing.md)

## 🔍 What is Devbox Pack?

Devbox Pack is an intelligent project detection and containerization tool that automatically analyzes your codebase, detects the programming languages and frameworks used, and generates optimized execution plans for deployment and development environments.

### Key Features

- **🔍 Automatic Detection** - Intelligently detects programming languages, frameworks, and build tools
- **📋 Execution Plans** - Generates comprehensive execution plans with build, dev, and start commands
- **🐳 Container Ready** - Optimized for containerized environments and cloud deployment
- **🛠️ Multi-Language** - Supports 10+ programming languages and their ecosystems
- **⚡ Fast & Efficient** - Minimal overhead with intelligent caching and optimization

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