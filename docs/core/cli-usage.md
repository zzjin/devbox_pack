# CLI Usage Guide

## Overview

The DevBox Pack CLI is a command-line tool that analyzes source code repositories and generates execution plans for containerized deployment. It automatically detects programming languages, frameworks, and generates appropriate build, development, and production commands.

## Installation

### Binary Installation

Download the latest release from the [GitHub releases page](https://github.com/labring/devbox-pack/releases) and extract the binary to your PATH.

### Build from Source

```bash
git clone https://github.com/labring/devbox-pack.git
cd devbox-pack
make build
```

## Basic Usage

### Command Syntax

```bash
devbox-pack <repository> [options]
```

### Arguments

- `repository` - Git repository URL or local path to analyze

### Basic Examples

```bash
# Analyze a remote repository
devbox-pack https://github.com/user/repo

# Analyze a local directory
devbox-pack . --offline

# Analyze with verbose output
devbox-pack /path/to/project --verbose

# Generate JSON output
devbox-pack https://github.com/user/repo --format json
```

## Command Options

### Repository Options

| Option | Description | Example |
|--------|-------------|---------|
| `--ref <ref>` | Git branch, tag, or commit to analyze | `--ref develop` |
| `--subdir <path>` | Subdirectory within the repository | `--subdir backend` |
| `--offline` | Analyze local directory without cloning | `--offline` |

### Detection Options

| Option | Description | Example |
|--------|-------------|---------|
| `--provider <name>` | Force use of specific provider | `--provider node` |
| `--platform <arch>` | Target platform architecture | `--platform linux/arm64` |
| `--base <name>` | Override base image selection | `--base base:node-18` |

### Output Options

| Option | Description | Example |
|--------|-------------|---------|
| `--format <format>` | Output format (pretty, json) | `--format json` |
| `--verbose` | Enable detailed logging | `--verbose` |
| `--quiet` | Suppress non-essential output | `--quiet` |

### Help Options

| Option | Description |
|--------|-------------|
| `-h, --help` | Show help information |
| `-v, --version` | Show version information |

## Supported Providers

The CLI automatically detects and supports the following providers:

| Provider | Languages/Frameworks | Priority |
|----------|----------------------|----------|
| `staticfile` | HTML, CSS, JS, Static sites | 10 (Highest) |
| `shell` | Shell scripts, Bash | 30 |
| `php` | PHP applications | 60 |
| `ruby` | Ruby, Rails applications | 65 |
| `java` | Java, Spring Boot, Maven, Gradle | 70 |
| `go` | Go applications | 75 |
| `node` | Node.js, JavaScript, TypeScript | 80 |
| `python` | Python, Django, Flask | 80 |
| `rust` | Rust applications | 85 |
| `deno` | Deno applications | 95 (Lowest) |

*Lower priority numbers indicate higher precedence in detection.*

## Output Formats

### Pretty Format (Default)

Human-readable output with colored formatting and structured information:

```bash
devbox-pack https://github.com/user/node-app
```

```
✅ Detection completed successfully

Provider: node
Language: node (v18.19.0)
Base Image: base:node-18
Port: 3000

Commands:
  Development: npm install && npm run dev -- --host 0.0.0.0 --port ${PORT}
  Build: npm run build
  Start: npm start -- --host 0.0.0.0 --port ${PORT}

Evidence:
  Files: package.json, package-lock.json
  Reason: Detected Node.js project with package.json and npm lockfile
```

### JSON Format

Machine-readable JSON output for integration with other tools:

```bash
devbox-pack https://github.com/user/node-app --format json
```

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
  "commands": {
    "dev": ["npm install", "npm run dev -- --host 0.0.0.0 --port ${PORT}"],
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

## Advanced Usage Examples

### Analyzing Specific Branches

```bash
# Analyze the develop branch
devbox-pack https://github.com/user/repo --ref develop

# Analyze a specific tag
devbox-pack https://github.com/user/repo --ref v1.2.3

# Analyze a specific commit
devbox-pack https://github.com/user/repo --ref abc123def
```

### Monorepo Support

```bash
# Analyze a specific service in a monorepo
devbox-pack https://github.com/user/monorepo --subdir services/api

# Analyze frontend in a full-stack repository
devbox-pack https://github.com/user/fullstack --subdir frontend
```

### Provider Override

```bash
# Force Node.js detection even if other providers match
devbox-pack . --provider node --offline

# Force static file serving for a React build
devbox-pack ./dist --provider staticfile --offline
```

### Platform-Specific Builds

```bash
# Generate plan for ARM64 architecture
devbox-pack https://github.com/user/repo --platform linux/arm64

# Generate plan for AMD64 architecture
devbox-pack https://github.com/user/repo --platform linux/amd64
```

### Custom Base Images

```bash
# Use a specific Node.js version
devbox-pack . --base base:node-20 --offline

# Use a specific Python version
devbox-pack . --base base:python-3.12 --offline
```

## Integration Examples

### CI/CD Pipeline

```yaml
# GitHub Actions example
- name: Generate Execution Plan
  run: |
    devbox-pack . --format json --offline > execution-plan.json
    
- name: Upload Plan Artifact
  uses: actions/upload-artifact@v3
  with:
    name: execution-plan
    path: execution-plan.json
```

### Docker Integration

```bash
# Generate plan and use in Dockerfile
devbox-pack . --format json --offline | jq -r '.base.name'
# Output: base:node-18
```

### Scripting

```bash
#!/bin/bash
# Automated deployment script

PLAN=$(devbox-pack https://github.com/user/repo --format json)
PROVIDER=$(echo "$PLAN" | jq -r '.provider')
BASE_IMAGE=$(echo "$PLAN" | jq -r '.base.name')

echo "Detected provider: $PROVIDER"
echo "Using base image: $BASE_IMAGE"

# Use the plan for deployment...
```

## Error Handling

### Common Exit Codes

- `0` - Success
- `1` - General error (invalid arguments, detection failure, etc.)

### Common Error Scenarios

#### No Repository Provided

```bash
devbox-pack
# Error: please provide repository path or URL
```

#### Invalid Repository

```bash
devbox-pack https://github.com/invalid/repo
# Error: failed to clone repository: repository not found
```

#### No Supported Language Detected

```bash
devbox-pack https://github.com/user/empty-repo
# Error: no supported language or framework detected
```

#### Invalid Provider

```bash
devbox-pack . --provider invalid --offline
# Error: unknown Provider: invalid, available Providers: [node python go java php ruby deno rust staticfile shell]
```

## Troubleshooting

### Enable Verbose Output

Use the `--verbose` flag to get detailed information about the detection process:

```bash
devbox-pack . --verbose --offline
```

This will show:
- File scanning progress
- Provider detection attempts
- Confidence scores
- Detection reasoning

### Check Available Providers

List all supported providers:

```bash
devbox-pack --help
```

### Validate Local Analysis

For local directories, always use the `--offline` flag:

```bash
devbox-pack /path/to/project --offline --verbose
```

### Debug Detection Issues

If detection fails or produces unexpected results:

1. Check that your project has the expected configuration files
2. Use `--verbose` to see detection reasoning
3. Try forcing a specific provider with `--provider`
4. Ensure you're analyzing the correct directory with `--subdir`

## Configuration Files

The CLI looks for these key files during detection:

### Node.js Projects
- `package.json` (required)
- `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml` (package managers)
- `.node-version`, `.nvmrc` (version specification)

### Python Projects
- `requirements.txt`, `pyproject.toml`, `setup.py` (dependencies)
- `runtime.txt` (version specification)
- `main.py`, `app.py` (entry points)

### Go Projects
- `go.mod` (required)
- `go.work` (workspace)
- `main.go` (entry point)

### Java Projects
- `pom.xml` (Maven)
- `build.gradle`, `build.gradle.kts` (Gradle)
- `src/main/java` (source structure)

### Static Files
- `index.html` (entry point)
- `dist/`, `build/`, `public/` (build directories)

## Best Practices

1. **Use Specific Branches**: Always specify `--ref` for production deployments
2. **Validate Locally**: Test with `--offline` before analyzing remote repositories
3. **Check Output**: Use `--verbose` to understand detection reasoning
4. **Version Control**: Include generated plans in version control for reproducibility
5. **Platform Consistency**: Specify `--platform` for consistent cross-platform builds
6. **Monorepo Organization**: Use `--subdir` for clear service boundaries in monorepos