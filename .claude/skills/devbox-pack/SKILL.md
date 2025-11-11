---
name: devbox-pack
description: Analyzes local or remote projects to generate deployment configuration JSON. Detects programming language, dependencies, and generates dev/prod runtime commands. Use when user asks to "analyze project", "generate deployment config".
---

# Project Packing and Analysis

Analyzes a project (local directory or Git repository) and generates a standardized JSON configuration for deployment. Detects the programming language, dependencies, execution commands, and port configuration for both development and production environments.

> **📁 Documentation Structure**
> - **SKILL.md** (this file) - Core workflow and detection rules
> - **[examples.md](examples.md)** - Complete working examples with explanations

## When to Use This Skill

- User asks to "pack" a project
- User mentions "analyze project" or "project analysis"
- User needs to "generate deployment config" or "runtime configuration"
- User wants to know how to run a project in dev/prod
- User provides a Git repository URL for analysis

## Core Principles

### 1. Runtime Environment: Debian Linux + Prebuilt Binaries

**Target system**: Debian-based Linux (debian:bullseye, ubuntu:22.04) with glibc

**Key implication**: Most native modules use prebuilt binaries. System dependencies are rarely needed.
- Node.js base images include graphics libs (canvas), SSL, compression
- Python base images include database client libs (psycopg2), image processing (Pillow)
- Go/Rust binaries are typically static

**Default to `"apt": []`** unless project documentation or Dockerfile explicitly requires system packages.

### 2. Network Binding: Always 0.0.0.0

**CRITICAL**: All dev and prod commands MUST bind to `0.0.0.0`, never `localhost` or `127.0.0.1`

**Syntax by language**:
- Node.js: `--host 0.0.0.0` or `HOST=0.0.0.0`
- Python: `0.0.0.0:PORT` or `--bind 0.0.0.0:PORT`
- Go: `:PORT` or `0.0.0.0:PORT`
- Java: `--server.address=0.0.0.0`
- Ruby: `-b 0.0.0.0`
- PHP: `-S 0.0.0.0:PORT`

### 3. Language Detection Priority

**CRITICAL**: Always follow this exact order:

```
PHP > Golang > Java > Rust > Ruby > Python > Deno > Node > StaticFile > Shell
```

Return up to 2 execution plans (one per detected language/subproject). Only traverse up to 2 levels of subdirectories.

**CRITICAL**: You can ONLY access files within the given project directory. You are **PROHIBITED** from:
  - Reading any files in parent directories (e.g., `../`, `../../`)
  - Accessing files outside the project root path
  - Using absolute paths that point outside the project directory
  - Following symlinks that lead outside the project directory

## Quick Detection Reference

| Priority | Language | Signature Files | Default Port |
|----------|----------|-----------------|--------------|
| 1 | PHP | `composer.json`, `*.php` | 8080 |
| 2 | Golang | `go.mod`, `*.go` | 8080 |
| 3 | Java | `pom.xml`, `build.gradle` | 8080 |
| 4 | Rust | `Cargo.toml`, `*.rs` | 8080 |
| 5 | Ruby | `Gemfile`, `*.rb` | 3000 |
| 6 | Python | `requirements.txt`, `*.py` | 8000 |
| 7 | Deno | `deno.json` + `.ts/.js` | 8000 |
| 8 | Node | `package.json` (no deno.json) | 3000 |
| 9 | Static | `index.html` (no backend) | 8080 |
| 10 | Shell | `*.sh`, `run.sh` | 8080 |

**Version detection priority**:
1. Version files (`.node-version`, `.python-version`, `.ruby-version`)
2. Package manifests (`package.json` engines, `go.mod` go directive)
3. Omit if uncertain

**Port detection priority**:
1. Explicitly configured in code
2. Environment variable in scripts
3. Framework default

## Analysis Workflow

### Step 1: Identify Source

```bash
# If URL contains .git or starts with http/https/git@
if [[ "$input" =~ \.git|^https?://|^git@ ]]; then
  temp_dir=$(mktemp -d)
  git clone --depth 1 --filter=tree:0 "$input" "$temp_dir"
  project_dir="$temp_dir"
else
  project_dir="$input"
fi
```

### Step 2: Detect Languages

Scan project directory (max 2 levels) for signature files in priority order. Stop after finding 2 languages.

**PHP**: `composer.json`, `index.php`, `*.php`
**Golang**: `go.mod` (required), `go.sum`, `main.go`
**Java**: `pom.xml` or `build.gradle`, `*.java`
**Rust**: `Cargo.toml` (required), `Cargo.lock`, `*.rs`
**Ruby**: `Gemfile`, `config.ru`, `*.rb`
**Python**: `requirements.txt`, `setup.py`, `pyproject.toml`, `*.py`
**Deno**: `deno.json` or `deno.jsonc` + `.ts/.js`
**Node.js**: `package.json` (and no deno.json)
**Static**: `index.html` + assets (no backend framework)
**Shell**: `*.sh`, `run.sh`, `start.sh`

### Step 3: Analyze Each Language

For each detected language:

1. **Extract version**: From version files or package manifests
2. **Determine dependencies**:
   - **APT packages**: Default to `[]` (empty)
   - Only include if project docs explicitly require system packages
3. **Generate commands**:
   - **Setup**: Install dependencies, build (if needed)
   - **Dev**: Development server with hot reload, bound to 0.0.0.0
   - **Prod**: Optimized production server, bound to 0.0.0.0
4. **Detect port**: From code/config or use framework default
5. **Collect evidence**: List detected files and reasoning

### Step 4: Generate JSON Output

Output a JSON array with 1-2 execution plans.

## Output Format

```json
[
  {
    "language": "string",
    "version": "string (optional)",
    "apt": ["string (optional)"],
    "dev": {
      "environment": {
        "KEY": "value"
      },
      "setup": ["command"],
      "commands": ["command with --host 0.0.0.0"]
    },
    "prod": {
      "environment": {
        "KEY": "value"
      },
      "setup": ["command"],
      "commands": ["command with --host 0.0.0.0"]
    },
    "port": 3000,
    "evidence": {
      "files": ["detected files"],
      "reason": "explanation of detection"
    }
  }
]
```

## Language-Specific Patterns

### Node.js

**Package manager detection**:
- `package-lock.json` → npm
- `yarn.lock` → yarn
- `pnpm-lock.yaml` → pnpm

**Common patterns**:
```json
{
  "language": "node",
  "version": "20.10.0",
  "apt": [],
  "dev": {
    "setup": ["npm install"],
    "commands": ["npm run dev -- --host 0.0.0.0 --port 3000"]
  },
  "prod": {
    "environment": {"NODE_ENV": "production"},
    "setup": ["npm install", "npm run build"],
    "commands": ["npm start -- --host 0.0.0.0 --port 3000"]
  },
  "port": 3000
}
```

### Python

**Framework detection**:
- `manage.py` → Django
- `flask` in requirements → Flask
- `fastapi` + `uvicorn` → FastAPI

**Common patterns**:
```json
{
  "language": "python",
  "version": "3.11",
  "apt": [],
  "dev": {
    "setup": ["pip install -r requirements.txt"],
    "commands": ["python manage.py runserver 0.0.0.0:8000"]
  },
  "prod": {
    "environment": {"PYTHONENV": "production"},
    "setup": ["pip install -r requirements.txt"],
    "commands": ["gunicorn --bind 0.0.0.0:8000 app:app"]
  },
  "port": 8000
}
```

### Golang

**Common patterns**:
```json
{
  "language": "go",
  "version": "1.23",
  "apt": [],
  "dev": {
    "environment": {
      "CGO_ENABLED": "0",
      "PORT": "8080",
      "HOST": "0.0.0.0"
    },
    "commands": ["go run ."]
  },
  "prod": {
    "environment": {
      "CGO_ENABLED": "0",
      "GO_ENV": "production",
      "PORT": "8080",
      "HOST": "0.0.0.0"
    },
    "setup": ["go build -o app ."],
    "commands": ["./app"]
  },
  "port": 8080
}
```

### Ruby

**Common patterns**:
```json
{
  "language": "ruby",
  "version": "3.2",
  "apt": [],
  "dev": {
    "setup": ["bundle install"],
    "commands": ["bundle exec rails server -b 0.0.0.0 -p 3000"]
  },
  "prod": {
    "environment": {"RAILS_ENV": "production"},
    "setup": ["bundle install --without development test"],
    "commands": ["bundle exec rails server -b 0.0.0.0 -p 3000"]
  },
  "port": 3000
}
```

## Special Cases

### Multi-Language Projects

**Example structure**:
```
project/
├── package.json          # Node.js (primary)
├── api/
│   └── go.mod           # Golang (secondary)
```

Generate separate execution plans for each language (max 2), respecting priority order.

### Static Files

Only detect as static if NO backend framework is present:

```json
{
  "language": "static",
  "apt": [],
  "dev": {
    "commands": ["python3 -m http.server 8080 --bind 0.0.0.0"]
  },
  "prod": {
    "commands": ["python3 -m http.server 8080 --bind 0.0.0.0"]
  },
  "port": 8080
}
```

## Validation Checklist

Before outputting JSON, verify:

- [ ] All commands bind to `0.0.0.0` (not localhost)
- [ ] Port number is included
- [ ] Evidence includes detected files and clear reason
- [ ] Language priority order is respected
- [ ] Dev and prod environments are differentiated
- [ ] Setup commands install necessary dependencies
- [ ] Commands are executable (no placeholders)
- [ ] APT packages default to `[]`
- [ ] JSON is valid and matches schema

## Error Handling

**No language detected**: Output empty array `[]`

**Git clone fails**: Report error, do not proceed

**Multiple subprojects > 2**: Select top 2 by priority

## Complete Examples

For detailed examples with full project analysis, see:

**→ [examples.md](examples.md)**

Available examples:
1. Simple Node.js Project (React + Vite)
2. Multi-Language Monorepo (Go + Node.js)
3. Python Django with PostgreSQL
4. Static Website
5. Go Microservice
6. Ruby on Rails Application

## Common Patterns Summary

| Package/Module | Language | APT Needed? | Why |
|----------------|----------|-------------|-----|
| `canvas` | Node.js | ❌ NO | Prebuilt binary + runtime libs in base image |
| `sharp` | Node.js | ❌ NO | Bundled libvips |
| `Pillow` | Python | ❌ NO | Binary wheels for glibc |
| `psycopg2-binary` | Python | ❌ NO | Binary wheel includes libpq |
| `lxml` | Python | ❌ NO | Binary wheels with libs |
| Go/Rust binaries | All | ❌ NO | Static or self-contained |

**When APT IS needed** (rare):
- CLI tools: `postgresql-client`, `redis-tools`
- Explicit in docs: `ffmpeg`, specialized system libs
