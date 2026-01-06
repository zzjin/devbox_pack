# Go Examples

Complete examples of Go project analysis and deployment configuration.

## Example 1: Go Microservice with Environment Configuration

**User request**: "Pack this Go API"

**Detected files**:
- `go.mod` with Go 1.23
- `main.go`
- `.env.example`

**Output**:
```json
[
  {
    "language": "go",
    "version": "1.23",
    "apt": [],
    "dev": {
      "environment": {
        "CGO_ENABLED": "0",
        "GOARCH": "amd64",
        "GOOS": "linux",
        "PORT": "8080",
        "HOST": "0.0.0.0"
      },
      "commands": ["go run ."]
    },
    "prod": {
      "environment": {
        "CGO_ENABLED": "0",
        "GOARCH": "amd64",
        "GOOS": "linux",
        "GO_ENV": "production",
        "PORT": "8080",
        "HOST": "0.0.0.0"
      },
      "setup": ["go build -o app ."],
      "commands": ["./app"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export CGO_ENABLED=0\n    export GOARCH=amd64\n    export GOOS=linux\n    export GO_ENV=production\n    export PORT=8080\n    export HOST=0.0.0.0\n    go build -o app .\n    ./app\nelse\n    export CGO_ENABLED=0\n    export GOARCH=amd64\n    export GOOS=linux\n    export PORT=8080\n    export HOST=0.0.0.0\n    go run .\nfi",
    "evidence": {
      "files": ["go.mod", "main.go"],
      "reason": "Go 1.23 project detected with go.mod module definition"
    }
  }
]
```

**Key Points**:
- Static binaries: `CGO_ENABLED=0`
- No APT dependencies for pure Go projects
- Dev: `go run .` (no build needed)
- Prod: Build to binary, then execute
- Environment variables control configuration

---

## Example 2: Go API in Subdirectory (Monorepo)

**User request**: "Analyze monorepo/backend"

**Detected files**:
- `backend/go.mod`
- `backend/main.go`
- `backend/go.sum`

**Output**:
```json
[
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
      "commands": ["cd backend && go run ."]
    },
    "prod": {
      "environment": {
        "CGO_ENABLED": "0",
        "GO_ENV": "production",
        "PORT": "8080",
        "HOST": "0.0.0.0"
      },
      "setup": ["cd backend && go build -o app ."],
      "commands": ["cd backend && ./app"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    cd backend\n    go build -o app .\n    CGO_ENABLED=0 GO_ENV=production PORT=8080 HOST=0.0.0.0 ./app\nelse\n    cd backend\n    CGO_ENABLED=0 PORT=8080 HOST=0.0.0.0 go run .\nfi",
    "evidence": {
      "files": ["backend/go.mod", "backend/main.go"],
      "reason": "Go API server detected in backend directory"
    }
  }
]
```

**Key Points**:
- Subdirectory handling: `cd backend &&` prepended
- Working directory preserved in entrypoint
- Independent module in monorepo structure

---

## Example 3: Go with Gin Framework

**User request**: "Pack Gin web service"

**Detected files**:
- `go.mod` with `github.com/gin-gonic/gin`
- `main.go`
- `server/router.go`

**Output**:
```json
[
  {
    "language": "go",
    "version": "1.22",
    "apt": [],
    "dev": {
      "environment": {
        "GIN_MODE": "debug",
        "PORT": "8080"
      },
      "commands": ["go run ."]
    },
    "prod": {
      "environment": {
        "GIN_MODE": "release",
        "CGO_ENABLED": "0",
        "PORT": "8080"
      },
      "setup": ["go build -ldflags='-s -w' -o app ."],
      "commands": ["./app"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export GIN_MODE=release\n    export CGO_ENABLED=0\n    export PORT=8080\n    go build -ldflags='-s -w' -o app .\n    ./app\nelse\n    export GIN_MODE=debug\n    export PORT=8080\n    go run .\nfi",
    "evidence": {
      "files": ["go.mod", "main.go", "server/router.go"],
      "reason": "Go project with Gin framework detected. Using optimized build flags for production."
    }
  }
]
```

**Key Points**:
- Gin mode: `debug` (dev) vs `release` (prod)
- Optimized build flags: `-ldflags='-s -w'` (strip symbols)
- Framework-specific environment variables

---

## Go Detection Patterns

### Required Files
- `go.mod` (REQUIRED for Go detection)
- `go.sum` (optional, dependency lock)
- `*.go` files

### Common Frameworks
| Framework | Detection | Notes |
|-----------|-----------|-------|
| Gin | `github.com/gin-gonic/gin` in go.mod | Popular web framework |
| Echo | `github.com/labstack/echo` in go.mod | High-performance framework |
| Fiber | `github.com/gofiber/fiber` in go.mod | Express-inspired framework |
| Chi | `github.com/go-chi/chi` in go.mod | Lightweight router |
| Standard lib | `net/http` in imports, no framework | Built-in HTTP server |

### Host Binding Patterns

**In Code**:
```go
// Using environment variable
port := os.Getenv("PORT")
addr := fmt.Sprintf(":%s", port)  // Binds to 0.0.0.0 by default

// Explicit binding
addr := "0.0.0.0:8080"
```

**Via Environment**:
```bash
export PORT=8080
export HOST=0.0.0.0
```

### Version Detection
1. `go.mod` → `go` directive line
   ```
   go 1.23
   ```
2. `.go-version` file (if using goenv)
3. Omit if uncertain

### Build Optimization

**Development**:
```bash
go run .
```

**Production**:
```bash
# Standard build
go build -o app .

# Optimized (smaller binary)
go build -ldflags='-s -w' -o app .

# With trimpath (reproducible builds)
go build -ldflags='-s -w' -trimpath -o app .
```

### Common Environment Variables
- `CGO_ENABLED`: `0` for static binaries, `1` for CGO
- `GOOS`: Target OS (linux, darwin, windows)
- `GOARCH`: Target arch (amd64, arm64)
- `GO_ENV`: Custom environment flag
- `PORT`: Application port
- `HOST`: Bind address
- `GIN_MODE`: Gin framework mode (debug/release)

### APT Dependencies
**Default: `[]` (empty array)**

Reasons:
- Go compiles to static binaries with `CGO_ENABLED=0`
- No runtime dependencies for pure Go code
- Standard library is self-contained

**Include APT packages only if**:
- Using CGO with system libraries
- Project explicitly requires CLI tools
- Database client utilities needed

### Subdirectory Handling

When Go module is in subdirectory:
- Prepend `cd dirname &&` to all commands
- Maintain working directory in entrypoint
- Example: `cd backend && go run .`

### Default Port
**8080** (common for Go microservices and APIs)

### Module Proxy
Go modules are downloaded automatically via:
- `GOPROXY` (default: https://proxy.golang.org)
- No explicit dependency installation needed
- Dependencies resolved from `go.mod` and `go.sum`

### Production Best Practices
1. **Static binaries**: Use `CGO_ENABLED=0`
2. **Size optimization**: Use `-ldflags='-s -w'`
3. **Security**: Keep dependencies updated
4. **Performance**: Multi-stage Docker builds (if applicable)
