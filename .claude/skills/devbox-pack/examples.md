# Pack Skill Examples

Complete examples of project analysis and deployment configuration generation.

## Example 1: React Project With Vite

**User request**: "Pack this project: examples/node-vite-react/"

**Detected files**:
- `package.json` with `vite`, `react`, `react-dom`
- `vite.config.ts`
- `package-lock.json`
- Scripts: `"dev": "vite"`, `"build": "tsc -b && vite build"`, `"preview": "vite preview"`

**Output**:
```json
[
  {
    "language": "node",
    "apt": [],
    "dev": {
      "setup": ["npm install"],
      "commands": ["npm run dev -- --host 0.0.0.0 --port 5173"]
    },
    "prod": {
      "environment": {
        "NODE_ENV": "production"
      },
      "setup": ["npm install", "npm run build"],
      "commands": ["npm run preview -- --host 0.0.0.0 --port 4173"]
    },
    "port": 4173,
    "evidence": {
      "files": ["package.json", "vite.config.ts", "package-lock.json"],
      "reason": "Node.js project with Vite + React detected. Using npm as package manager. Vite preview server runs on port 4173 by default."
    }
  }
]
```

**Key Points**:
- Vite detected from `vite.config.ts` and dependencies
- Detect scripts in `package.json` for dev/prod commands
- Dev port and Pord port: 5173 (Vite default)
- APT packages: Empty array (Vite uses pure JavaScript tooling)
- Must pass `--host 0.0.0.0` to both dev and preview commands

---

## Example 2: Multi-Language Monorepo

**User request**: "Analyze /path/to/monorepo"

**Structure**:
```
monorepo/
├── package.json (Node.js workspace)
├── frontend/
│   └── package.json (React/Vite)
└── backend/
    └── go.mod (Golang API)
```

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
    "evidence": {
      "files": ["backend/go.mod", "backend/main.go"],
      "reason": "Go API server detected in backend directory"
    }
  },
  {
    "language": "node",
    "version": "20.10.0",
    "apt": [],
    "dev": {
      "environment": {},
      "setup": ["cd frontend && npm install"],
      "commands": ["cd frontend && npm run dev -- --host 0.0.0.0 --port 3000"]
    },
    "prod": {
      "environment": {
        "NODE_ENV": "production"
      },
      "setup": ["cd frontend && npm install && npm run build"],
      "commands": ["cd frontend && npm start -- --host 0.0.0.0 --port 3000"]
    },
    "port": 3000,
    "evidence": {
      "files": ["frontend/package.json", "frontend/vite.config.ts"],
      "reason": "Node.js React frontend with Vite detected in frontend directory"
    }
  }
]
```

**Key Points**:
- Language priority order respected: Golang (priority 2) comes before Node.js (priority 8)
- Each execution plan is independent with its own working directory
- Maximum 3 execution plans returned
- Different ports assigned to avoid conflicts (8080 for backend, 3000 for frontend)

---

## Example 3: Python Django with PostgreSQL

**User request**: "Pack my Django project"

**Detected files**:
- `requirements.txt` with `Django==4.2` and `psycopg2`
- `manage.py`
- `.python-version` → 3.11

**Output**:
```json
[
  {
    "language": "python",
    "version": "3.11",
    "apt": [],
    "dev": {
      "environment": {
        "DEBUG": "True"
      },
      "setup": ["pip install -r requirements.txt", "python manage.py migrate"],
      "commands": ["python manage.py runserver 0.0.0.0:8000"]
    },
    "prod": {
      "environment": {
        "PYTHONENV": "production",
        "DEBUG": "False"
      },
      "setup": [
        "pip install -r requirements.txt",
        "python manage.py migrate",
        "python manage.py collectstatic --noinput"
      ],
      "commands": ["gunicorn --bind 0.0.0.0:8000 --workers 4 myproject.wsgi:application"]
    },
    "port": 8000,
    "evidence": {
      "files": ["requirements.txt", "manage.py", ".python-version"],
      "reason": "Python Django 4.2 project detected. psycopg2 uses prebuilt binaries on Debian/glibc, no system dependencies needed."
    }
  }
]
```

**Key Points**:
- APT packages: Empty array even though project uses `psycopg2`
  - Reason: `psycopg2-binary` wheels include libpq
  - On Debian/glibc, no postgresql-dev packages needed
- Django-specific setup commands included (migrate, collectstatic)
- Production uses gunicorn instead of Django dev server
- Debug mode toggled via environment variables

---

## Example 4: Static Website

**User request**: "Pack this HTML site"

**Detected files**:
- `index.html`
- `css/style.css`
- `js/app.js`
- No backend frameworks

**Output**:
```json
[
  {
    "language": "static",
    "apt": [],
    "dev": {
      "environment": {},
      "commands": ["python3 -m http.server 8080 --bind 0.0.0.0"]
    },
    "prod": {
      "environment": {},
      "commands": ["python3 -m http.server 8080 --bind 0.0.0.0"]
    },
    "port": 8080,
    "evidence": {
      "files": ["index.html", "css/style.css", "js/app.js"],
      "reason": "Static HTML/CSS/JS website detected with no backend framework"
    }
  }
]
```

**Key Points**:
- Static files have lowest priority in detection order
- Same command for dev and prod (simple HTTP server)
- No setup commands needed
- No version field (not applicable for static files)

---

## Example 5: Go Microservice with Environment-Based Configuration

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
    "evidence": {
      "files": ["go.mod", "main.go"],
      "reason": "Go 1.23 project detected with go.mod module definition"
    }
  }
]
```

**Key Points**:
- Go binaries are typically statically compiled (CGO_ENABLED=0)
- No APT dependencies needed for pure Go projects
- Build creates standalone binary in production
- Environment variables control host/port binding

---

## Example 6: Ruby on Rails Application

**User request**: "Analyze Rails app"

**Detected files**:
- `Gemfile` with `rails 7.0`
- `config/application.rb`
- `.ruby-version` → 3.2

**Output**:
```json
[
  {
    "language": "ruby",
    "version": "3.2",
    "apt": [],
    "dev": {
      "environment": {},
      "setup": ["bundle install"],
      "commands": ["bundle exec rails server -b 0.0.0.0 -p 3000"]
    },
    "prod": {
      "environment": {
        "RAILS_ENV": "production"
      },
      "setup": [
        "bundle install --without development test",
        "bundle exec rails db:migrate"
      ],
      "commands": ["bundle exec rails server -b 0.0.0.0 -p 3000"]
    },
    "port": 3000,
    "evidence": {
      "files": ["Gemfile", "config/application.rb", ".ruby-version"],
      "reason": "Ruby on Rails 7.0 application detected"
    }
  }
]
```

**Key Points**:
- Rails-specific commands (db:migrate)
- Bundle install with different groups for dev/prod
- Version from `.ruby-version` file
- RAILS_ENV environment variable for production mode

---

## Common Patterns Across Examples

### APT Dependencies
- **All examples use `"apt": []`** (empty array)
- Reason: Debian base images + prebuilt binaries handle everything
- Only exception: CLI tools explicitly mentioned in docs (postgresql-client, ffmpeg, etc.)

### Host Binding
- **All commands bind to `0.0.0.0`**, never localhost
- Critical for container/cloud deployments
- Various syntax: `--host`, `-b`, `--bind`, `--hostname`, or environment variable

### Version Detection Priority
1. Version files (`.node-version`, `.python-version`, `.ruby-version`)
2. Package manifests (`package.json` engines, `go.mod` go directive)
3. Dockerfile `FROM` statement (if needed)
4. Omit if uncertain

### Port Assignment Strategy
- Default framework ports where possible
- Avoid conflicts in multi-service projects
- Always include in JSON output
- Common ports: 3000 (Node.js/Rails), 8000 (Python/Django), 8080 (Go/Java)

### Evidence Quality
- List actual files found during detection
- Explain reasoning clearly
- Mention special cases (e.g., "psycopg2 uses prebuilt binaries")
- Reference version numbers when known
