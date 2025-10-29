# DevBox Pack

**Intelligent execution plan generator for containerized development environments**

DevBox Pack is a sophisticated static analysis tool that automatically detects programming languages, frameworks, and generates optimized execution plans for containerized development workflows. This repository contains both Go and TypeScript/Node.js implementations.

## 🏗️ Architecture Overview

- High-performance static analysis engine with advanced provider detection

## ✨ Key Features

- 🔍 **Intelligent Multi-Language Detection**: Advanced confidence-based detection for 11+ programming languages
- 🎯 **Framework-Aware Analysis**: Detects specific frameworks (Next.js, Django, Spring Boot, etc.)
- 📋 **Execution Plan Generation**: Complete containerization configuration with optimized build/dev/start commands
- 🚀 **Git Repository Support**: Direct analysis of remote repositories with branch/tag support
- 📊 **Multiple Output Formats**: JSON and human-readable pretty formats
- ⚡ **High Performance**: Native Go implementation with sub-second analysis times
- 🔧 **Extensible Provider System**: Priority-based detection with confidence scoring algorithms

## 🌐 Supported Languages & Frameworks

| Language | Package Managers | Frameworks Detected | Priority |
|----------|------------------|-------------------|----------|
| **Node.js** | npm, yarn, pnpm, bun | Next.js, Express, Nuxt, Vite | 100 |
| **Python** | pip, poetry, pipenv, uv | Django, Flask, FastAPI | 90 |
| **Java** | Maven, Gradle | Spring Boot, Quarkus | 80 |
| **Go** | go.mod | Gin, Echo, Fiber | 70 |
| **PHP** | Composer | Laravel, Symfony | 60 |
| **Ruby** | Bundler | Rails, Sinatra | 50 |
| **Deno** | deno.json | Fresh, Oak | 45 |
| **Rust** | Cargo | Actix, Rocket | 40 |
| **Static** | - | HTML/CSS/JS | 20 |
| **Shell** | - | Bash scripts | 10 |

## 📦 Installation

```bash
# Clone repository
git clone https://github.com/labring/devbox-pack.git
cd devbox-pack

# Build from source
make build

# Install globally
make install

# Or run directly
./bin/devbox-pack --help
```

## 🚀 Quick Start

### Basic Usage

```bash
# Analyze any Git repository
devbox-pack https://github.com/vercel/next.js

# Analyze local project
devbox-pack . --offline

# Get JSON output for programmatic use
devbox-pack /path/to/project --format json

# Force specific provider detection
devbox-pack . --provider node --verbose
```

### CLI Options

```
Usage: devbox-pack <repository> [options]

Arguments:
  repository               Git repository URL or local directory path

Options:
  -h, --help              Show help information
  -v, --version           Show version information
  --ref <ref>             Git branch, tag, or commit (default: main)
  --subdir <path>         Analyze subdirectory within repository
  --provider <name>       Force specific provider (node|python|java|go|php|ruby|deno|rust|staticfile|shell)
  --format <format>       Output format: pretty (default) | json
  --verbose               Enable detailed detection information
  --offline               Skip git operations, analyze local files only
  --platform <arch>       Target platform architecture (e.g., linux/amd64)
  --base <name>           Override base image selection
```

### Real-World Examples

```bash
# Next.js application with specific subdirectory
devbox-pack https://github.com/vercel/next.js --subdir examples/hello-world

# Python Django project with JSON output
devbox-pack https://github.com/django/django --format json --verbose

# Local Go microservice
devbox-pack . --offline --provider go --platform linux/arm64

# Analyze specific branch
devbox-pack https://github.com/user/repo --ref develop --format json
```

## 📚 Documentation

- **[Core Architecture Overview](docs/core/overview.md)** - Technical deep-dive into detection engine and provider system
- **[API Schema & Data Structures](docs/core/api-schema.md)** - Complete API documentation and data models
- **[Provider Documentation](docs/providers/)** - Language-specific detection logic and configuration
- **[Development Guide](docs/README.md)** - Comprehensive development and contribution guide

## 🏗️ Project Structure

```
devbox-pack/
├── cmd/                    # Go application entry point
│   └── main.go            # CLI application main
├── pkg/                   # Go implementation core packages
│   ├── cli/              # Command-line interface logic
│   ├── detector/         # Detection engine and provider coordination
│   ├── providers/        # Language-specific detection providers
│   ├── generators/       # Execution plan generation logic
│   ├── formatters/       # Output formatting (JSON, Pretty)
│   ├── git/             # Git repository operations
│   ├── types/           # Core data structures and interfaces
│   └── utils/           # Shared utilities
├── docs/                # Comprehensive documentation
│   ├── core/           # Architecture and API documentation
│   └── providers/      # Provider-specific documentation
├── tests/              # Integration and validation tests
├── go.mod             # Go module definition
├── Makefile          # Build automation and development tasks
└── README.md         # This file
```

## 🛠️ Development

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linting
make lint

# Build for development
make build

# Clean build artifacts
make clean

# Show all available commands
make help
```

## 📊 Output Formats & Examples

### JSON Format (Programmatic Integration)

Perfect for CI/CD pipelines and automated tooling:

```json
{
  "provider": "node",
  "confidence": 0.95,
  "base": {
    "name": "base:node-20",
    "platform": "linux/amd64"
  },
  "runtime": {
    "language": "node",
    "version": "20.10.0",
    "framework": "nextjs",
    "tools": ["pnpm"],
    "env": {
      "NODE_ENV": "production",
      "NEXT_TELEMETRY_DISABLED": "1"
    }
  },
  "apt": [],
  "dev": {
    "cmd": "pnpm dev",
    "port": 3000,
    "notes": ["Development server with hot reload"]
  },
  "build": {
    "cmd": "pnpm build",
    "caches": ["node_modules", ".next"],
    "notes": ["Optimized production build"]
  },
  "start": {
    "cmd": "pnpm start",
    "portEnv": "PORT",
    "notes": ["Production server startup"]
  },
  "evidence": {
    "files": ["package.json", "pnpm-lock.yaml", "next.config.js"],
    "reason": "Next.js project detected with pnpm package manager",
    "confidence_breakdown": {
      "package.json": 0.4,
      "pnpm-lock.yaml": 0.3,
      "next.config.js": 0.25
    }
  }
}
```

### Pretty Format (Human-Readable)

Ideal for development and debugging:

```
🔍 DevBox Pack Execution Plan

📋 Project Analysis:
  Provider: node (confidence: 95%)
  Language: Node.js 20.10.0
  Framework: Next.js
  Package Manager: pnpm

🏗️ Base Environment:
  Image: base:node-20
  Platform: linux/amd64

⚙️ Development Environment:
  Command: pnpm dev
  Port: 3000
  Notes: Development server with hot reload

🔨 Build Configuration:
  Command: pnpm build
  Cache Directories: node_modules, .next
  Notes: Optimized production build

🚀 Production Startup:
  Command: pnpm start
  Port Environment: PORT
  Notes: Production server startup

📁 Detection Evidence:
  Files Found: package.json, pnpm-lock.yaml, next.config.js
  Detection Reason: Next.js project detected with pnpm package manager
  Confidence Breakdown:
    • package.json: 40%
    • pnpm-lock.yaml: 30%
    • next.config.js: 25%
```

## 🔧 Advanced Usage

### Integration with CI/CD

```bash
# Generate execution plan for deployment
devbox-pack https://github.com/user/repo --format json > execution-plan.json

# Use in GitHub Actions
- name: Generate DevBox Plan
  run: |
    devbox-pack . --format json --offline > plan.json
    echo "DEVBOX_PLAN=$(cat plan.json)" >> $GITHUB_ENV
```

### Custom Provider Detection

```bash
# Force specific provider when auto-detection fails
devbox-pack . --provider python --verbose

# Override base image selection
devbox-pack . --base "custom:python-3.11" --format json

# Target specific platform
devbox-pack . --platform linux/arm64 --provider go
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. **Fork and Clone**: Fork the repository and clone your fork
2. **Create Branch**: Create a feature branch for your changes
3. **Develop**: Make your changes with appropriate tests
4. **Test**: Run the full test suite (`make test`)
5. **Submit**: Create a pull request with a clear description

### Adding New Providers

See the [Provider Development Guide](docs/providers/README.md) for detailed instructions on adding support for new languages and frameworks.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Built with ❤️ by the [Labring](https://github.com/labring) team and contributors.