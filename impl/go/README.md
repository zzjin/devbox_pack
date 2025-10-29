# DevBox Pack - Go Implementation

Go language implementation of DevBox Pack execution plan generator.

## Features

- 🔍 Intelligent detection of 11 programming languages and frameworks
- 📋 Generate complete execution plan configuration
- 🚀 Support for Git repositories and local project analysis
- 📊 Multiple output formats (JSON, Pretty)
- ⚡ High performance with native Go implementation

## Supported Languages

- Node.js (npm, yarn, pnpm, bun)
- Python (pip, poetry, pipenv, uv)
- Java (Maven, Gradle)
- Go (go.mod)
- PHP (Composer)
- Ruby (Bundler)
- Deno
- Rust (Cargo)
- Static files
- Shell scripts

## Installation

### Build from Source

```bash
# Clone repository
git clone https://github.com/labring/devbox-pack.git
cd devbox-pack/impl/go

# Build
make build

# Install to GOPATH
make install
```

## Usage

### Basic Usage

```bash
# Analyze Git repository
./bin/devbox-pack https://github.com/user/repo

# Analyze local project
./bin/devbox-pack . --offline

# Specify output format
./bin/devbox-pack /path/to/project --format json

# Verbose output
./bin/devbox-pack https://github.com/user/repo --verbose
```

### Command Line Options

```
Usage:
  devbox-pack <repository> [options]

Arguments:
  repository               Git repository URL or local path

Options:
  -h, --help              Show help information
  -v, --version           Show version information
  --ref <ref>             Git branch or tag (default: main)
  --subdir <path>         Subdirectory path
  --provider <name>       Force use specified Provider
  --format <format>      Output format (pretty|json, default: pretty)
  --verbose               Show verbose information
  --offline               Offline mode, don't clone repository
  --platform <arch>       Target platform (e.g.: linux/amd64)
  --base <name>           Specify base image
```

### Examples

```bash
# Analyze Next.js project
devbox-pack https://github.com/vercel/next.js --subdir examples/hello-world

# Analyze Python Django project
devbox-pack https://github.com/django/django --format json

# Analyze local Go project
devbox-pack . --offline --verbose --provider go
```

## Development

### Project Structure

```
impl/go/
├── cmd/devbox-pack/     # Main entry point
├── pkg/
│   ├── cli/            # CLI interface
│   ├── types/          # Data type definitions
│   ├── git/            # Git operations
│   ├── providers/      # Language detectors
│   ├── generator/      # Execution plan generator
│   └── formatters/     # Output formatters
├── go.mod              # Go module definition
├── Makefile           # Build scripts
└── README.md          # Documentation
```

### Development Commands

```bash
# Run tests
make test

# Format code
make fmt

# Code linting
make lint

# Clean build artifacts
make clean

# Show help
make help
```

## Output Formats

### JSON Format

```json
{
  "provider": "node",
  "base": {
    "name": "base:node-20",
    "platform": "linux/amd64"
  },
  "runtime": {
    "language": "node",
    "version": "20.10.0",
    "tools": ["pnpm"],
    "env": {
      "NODE_ENV": "production"
    }
  },
  "apt": [],
  "dev": {
    "cmd": "pnpm dev",
    "notes": ["Development server running at http://localhost:3000"]
  },
  "build": {
    "cmd": "pnpm build",
    "caches": ["node_modules", ".next"],
    "notes": ["Build artifacts in .next directory"]
  },
  "start": {
    "cmd": "pnpm start",
    "portEnv": "PORT",
    "notes": ["Production server startup"]
  },
  "evidence": {
    "files": ["package.json", "pnpm-lock.yaml"],
    "reason": "Detected package.json and pnpm-lock.yaml"
  }
}
```

### Pretty Format

```
🔍 DevBox Pack Execution Plan

📋 Project Information:
  Provider: node
  Language: Node.js 20.10.0
  Tools: pnpm

🏗️ Base Environment:
  Image: base:node-20
  Platform: linux/amd64

⚙️ Development Environment:
  Command: pnpm dev
  Description: Development server running at http://localhost:3000

🔨 Build Configuration:
  Command: pnpm build
  Cache: node_modules, .next
  Description: Build artifacts in .next directory

🚀 Production Startup:
  Command: pnpm start
  Port: PORT
  Description: Production server startup

📁 Detection Evidence:
  Files: package.json, pnpm-lock.yaml
  Reason: Detected package.json and pnpm-lock.yaml
```

## License

MIT License