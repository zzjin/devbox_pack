---
name: devbox-pack
description: Analyzes local or remote projects to generate deployment configuration JSON. Detects programming language, dependencies, and generates dev/prod runtime commands. Use when user asks to "analyze project", "generate deployment config".
---

# Project Packing and Analysis

Analyzes a project (local directory or Git repository) and generates standardized JSON configuration for deployment. Detects programming language, dependencies, and runtime commands for dev/prod environments.

> **📁 Progressive Documentation**
> - **SKILL.md** (this file) - Core workflow and quick reference
> - **[examples/](examples/)** - Language-specific examples:
>   - [nodejs.md](examples/nodejs.md) - Node.js, Next.js, Vite, Express
>   - [python.md](examples/python.md) - Django, FastAPI, Flask
>   - [golang.md](examples/golang.md) - Go microservices, Gin
>   - [ruby.md](examples/ruby.md) - Rails, Sinatra, Puma
>   - [php.md](examples/php.md) - Laravel, Symfony
>   - [java.md](examples/java.md) - Spring Boot, Maven, Gradle
>   - [rust.md](examples/rust.md) - Actix-web, Axum, Rocket
>   - [deno.md](examples/deno.md) - Fresh, Oak, Deno runtime
>   - [static.md](examples/static.md) - HTML/CSS/JS, SPAs

## Core Workflow

### Step 1: Detect Language

Scan project directory (max 2 levels) for signature files in priority order:

```
PHP > Golang > Java > Rust > Ruby > Python > Deno > Node > Static > Shell
```

| Priority | Language | Signature Files | Default Port |
|----------|----------|-----------------|--------------|
| 1 | PHP | `composer.json`, `*.php` | 8080 |
| 2 | Golang | `go.mod` (required) | 8080 |
| 3 | Java | `pom.xml`, `build.gradle` | 8080 |
| 4 | Rust | `Cargo.toml` (required) | 8080 |
| 5 | Ruby | `Gemfile` | 3000 |
| 6 | Python | `requirements.txt`, `*.py` | 8000 |
| 7 | Deno | `deno.json` + `.ts/.js` | 8000 |
| 8 | Node | `package.json` (no deno.json) | 3000 |
| 9 | Static | `index.html` (no backend) | 8080 |
| 10 | Shell | `*.sh`, `run.sh` | 8080 |

Return up to 2 execution plans (one per detected language/subproject).

### Step 2: Analyze Project

For each detected language:

1. **Determine framework** - Check dependencies/config files
2. **Extract version** - From version files (`.node-version`, `.python-version`) or manifests
3. **Generate commands**:
   - Setup: Install dependencies, build if needed
   - Dev: Development server with hot reload, **bound to 0.0.0.0**
   - Prod: Optimized production server, **bound to 0.0.0.0**
4. **Detect port** - From code/config or use framework default
5. **Build evidence** - List detected files and reasoning

**For language-specific patterns, see → [examples/](examples/)**

### Step 3: Generate JSON Output

Output JSON array with 1-2 execution plans matching this structure:

```json
[
  {
    "language": "string",
    "version": "string (optional)",
    "apt": [],
    "dev": {
      "environment": {"KEY": "value"},
      "setup": ["command"],
      "commands": ["command with --host 0.0.0.0"]
    },
    "prod": {
      "environment": {"KEY": "value"},
      "setup": ["command"],
      "commands": ["command with --host 0.0.0.0"]
    },
    "port": 3000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    # prod setup and commands\nelse\n    # dev setup and commands\nfi",
    "evidence": {
      "files": ["detected files"],
      "reason": "explanation of detection"
    }
  }
]
```

## Critical Rules

### 1. Network Binding: Always 0.0.0.0

**All dev and prod commands MUST bind to `0.0.0.0`**, never `localhost` or `127.0.0.1`.

Syntax by language:
- Node: `--host 0.0.0.0`
- Python: `0.0.0.0:PORT` or `--bind 0.0.0.0:PORT`
- Go: `:PORT` or `HOST=0.0.0.0`
- Ruby: `-b 0.0.0.0`
- PHP: `--host=0.0.0.0`
- Java: `SERVER_ADDRESS=0.0.0.0`

### 2. APT Dependencies: Default to Empty

**Default: `"apt": []`** for all languages.

Reasons:
- Debian base images include common libraries
- Prebuilt binaries handle native dependencies (canvas, sharp, psycopg2-binary, Pillow)
- Go/Rust produce static binaries

**Only include APT packages if**:
- Project documentation explicitly requires them
- Dockerfile shows system package installation
- CLI tools needed: `postgresql-client`, `ffmpeg`

### 3. Entrypoint Script

**REQUIRED**: The `entrypoint` field must be a unified shell script:

```bash
#!/bin/bash
app_env=${1:-development}
if [ "$app_env" = "production" ] || [ "$app_env" = "prod" ] ; then
    # Export prod environment vars
    # Run prod setup commands
    # Run prod commands
else
    # Export dev environment vars
    # Run dev setup commands
    # Run dev commands
fi
```

**Rules**:
- Accept optional first argument (defaults to `development`)
- Check for `production` or `prod`
- Export all environment variables using `export KEY=value`
- Run setup before main commands
- Preserve `cd` directory changes
- For identical dev/prod (static), omit if/else

### 4. Security Constraints

**You can ONLY access files within the given project directory.**

**PROHIBITED**:
- Reading parent directories (`../`, `../../`)
- Accessing files outside project root
- Using absolute paths outside project
- Following symlinks outside project

### 5. Environment Variables

**Remove localhost URLs from production environment**:

❌ Remove: `NEXT_PUBLIC_APP_URL=http://localhost:3000`, `VITE_API_URL=http://localhost:8080`

✅ Keep internal services: `DATABASE_URL=postgresql://localhost:5432/db`

Common patterns to check:
- `NEXT_PUBLIC_*`, `VITE_*`, `REACT_APP_*`
- `API_BASE_URL`, `APP_URL`, `BACKEND_URL`, `PUBLIC_URL`

## Validation Checklist

Before outputting JSON, verify:

- [ ] All commands bind to `0.0.0.0` (not localhost)
- [ ] Port number included
- [ ] Evidence lists detected files and clear reasoning
- [ ] Language priority order respected
- [ ] Dev and prod environments differentiated
- [ ] APT packages default to `[]`
- [ ] JSON valid and matches schema
- [ ] Entrypoint script generated and properly formatted
- [ ] Entrypoint exports environment variables
- [ ] Entrypoint runs setup before main commands

## Error Handling

- **No language detected**: Output empty array `[]`
- **Git clone fails**: Report error, do not proceed
- **Multiple subprojects > 2**: Select top 2 by priority

## Language-Specific Guides

**When you need detailed guidance for a specific language**:

- **Node.js projects** → See [examples/nodejs.md](examples/nodejs.md)
  - Package manager detection (npm/yarn/pnpm)
  - Framework patterns (Vite, Next.js, Express)
  - Host binding syntax

- **Python projects** → See [examples/python.md](examples/python.md)
  - Framework detection (Django, FastAPI, Flask)
  - WSGI/ASGI server selection
  - Django-specific commands

- **Go projects** → See [examples/golang.md](examples/golang.md)
  - Static binary compilation
  - Subdirectory handling
  - Gin/Echo framework patterns

- **Ruby projects** → See [examples/ruby.md](examples/ruby.md)
  - Rails vs Sinatra detection
  - Bundle install strategies
  - Puma server configuration

- **PHP projects** → See [examples/php.md](examples/php.md)
  - Laravel/Symfony detection
  - Composer strategies
  - Artisan commands

- **Java projects** → See [examples/java.md](examples/java.md)
  - Spring Boot with Maven/Gradle
  - JVM options and profiles
  - Build tool detection

- **Rust projects** → See [examples/rust.md](examples/rust.md)
  - Actix-web, Axum, Rocket frameworks
  - Release build optimization
  - Cargo commands

- **Deno projects** → See [examples/deno.md](examples/deno.md)
  - Fresh, Oak frameworks
  - Permission flags
  - TypeScript runtime

- **Static sites** → See [examples/static.md](examples/static.md)
  - Pure HTML/CSS/JS
  - Built SPAs
  - Static file serving
