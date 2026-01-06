# Deno Examples

Complete examples of Deno project analysis and deployment configuration.

## Example 1: Deno Web Application

**User request**: "Analyze Deno project"

**Detected files**:
- `deno.json` or `deno.jsonc`
- `main.ts` or `mod.ts`
- TypeScript/JavaScript files

**Output**:
```json
[
  {
    "language": "deno",
    "version": "1.40",
    "apt": [],
    "dev": {
      "environment": {},
      "commands": ["deno run --allow-net --allow-read --allow-env --watch main.ts"]
    },
    "prod": {
      "environment": {
        "DENO_ENV": "production"
      },
      "commands": ["deno run --allow-net --allow-read --allow-env main.ts"]
    },
    "port": 8000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export DENO_ENV=production\n    deno run --allow-net --allow-read --allow-env main.ts\nelse\n    deno run --allow-net --allow-read --allow-env --watch main.ts\nfi",
    "evidence": {
      "files": ["deno.json", "main.ts"],
      "reason": "Deno application detected. Using --watch flag for development hot reload."
    }
  }
]
```

**Key Points**:
- Deno requires explicit permissions (--allow-*)
- Dev: `--watch` flag for hot reload
- No build step needed (TypeScript runs directly)
- Default port: 8000

---

## Example 2: Deno Fresh Framework

**User request**: "Pack Fresh app"

**Detected files**:
- `deno.json` with `fresh` in imports
- `fresh.gen.ts`
- `routes/` directory
- `islands/` directory

**Output**:
```json
[
  {
    "language": "deno",
    "version": "1.40",
    "apt": [],
    "dev": {
      "environment": {},
      "commands": ["deno task start"]
    },
    "prod": {
      "environment": {
        "DENO_ENV": "production"
      },
      "commands": ["deno task start"]
    },
    "port": 8000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export DENO_ENV=production\n    deno task start\nelse\n    deno task start\nfi",
    "evidence": {
      "files": ["deno.json", "fresh.gen.ts", "routes/"],
      "reason": "Deno Fresh framework detected. Using task runner from deno.json."
    }
  }
]
```

**Key Points**:
- Fresh is Deno's full-stack web framework
- Uses `deno task` for task running
- No build step in dev (JIT rendering)
- Production uses same command (Fresh handles optimization)

---

## Example 3: Deno Oak API

**User request**: "Analyze Oak API server"

**Detected files**:
- `deno.json`
- Import of Oak in `mod.ts` or `main.ts`
- `deps.ts` (dependency management)

**Output**:
```json
[
  {
    "language": "deno",
    "version": "1.39",
    "apt": [],
    "dev": {
      "environment": {
        "PORT": "8000"
      },
      "commands": ["deno run --allow-net --allow-env --watch mod.ts"]
    },
    "prod": {
      "environment": {
        "PORT": "8000",
        "DENO_ENV": "production"
      },
      "commands": ["deno run --allow-net --allow-env mod.ts"]
    },
    "port": 8000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export PORT=8000\n    export DENO_ENV=production\n    deno run --allow-net --allow-env mod.ts\nelse\n    export PORT=8000\n    deno run --allow-net --allow-env --watch mod.ts\nfi",
    "evidence": {
      "files": ["deno.json", "mod.ts", "deps.ts"],
      "reason": "Deno Oak framework detected for REST API server."
    }
  }
]
```

**Key Points**:
- Oak is Deno's middleware framework (like Express/Koa)
- `deps.ts` pattern for centralized dependency management
- Minimal permissions: only `--allow-net` and `--allow-env`

---

## Deno Detection Patterns

### Required Files
- `deno.json` or `deno.jsonc` (REQUIRED - Deno configuration)
- TypeScript (`.ts`) or JavaScript (`.js`) files
- Common entry points: `main.ts`, `mod.ts`, `server.ts`

### Framework Detection
| Framework | Detection | Use Case | Default Port |
|-----------|-----------|----------|--------------|
| Fresh | `fresh` in deno.json imports, `fresh.gen.ts` | Full-stack web framework | 8000 |
| Oak | Oak import in code, `deps.ts` | Middleware/REST API | 8000 |
| Aleph | `aleph` in deno.json | React SSR framework | 8000 |
| Ultra | `ultra` in deno.json | React framework | 8000 |
| Hono | `hono` import | Ultrafast web framework | 8000 |

### Host Binding Patterns

Deno frameworks typically bind via environment or code:

```typescript
// Oak
const app = new Application();
await app.listen({ hostname: "0.0.0.0", port: 8000 });

// Fresh (configured in fresh.config.ts)
export default defineConfig({
  server: {
    hostname: "0.0.0.0",
    port: 8000,
  },
});

// Environment variable
const port = parseInt(Deno.env.get("PORT") || "8000");
```

### Version Detection
1. `deno.json` → no version field (Deno is globally installed)
2. `.dvmrc` file (Deno Version Manager)
3. Docker image version
4. Omit version field (Deno auto-updates)

### Deno Permission Flags

**Common permissions**:
```bash
--allow-net              # Network access
--allow-read             # File system read
--allow-write            # File system write
--allow-env              # Environment variables
--allow-run              # Subprocess execution
--allow-all              # All permissions (not recommended)
```

**Granular permissions**:
```bash
--allow-net=0.0.0.0:8000           # Specific host/port
--allow-read=/app,/tmp             # Specific directories
--allow-env=PORT,DATABASE_URL      # Specific env vars
```

### Deno Commands

**Development**:
```bash
deno run --watch main.ts              # Hot reload
deno task dev                         # Run dev task from deno.json
deno run --inspect main.ts            # With debugger
```

**Production**:
```bash
deno run main.ts                      # Direct execution
deno task start                       # Run start task
deno compile --output app main.ts    # Compile to binary (optional)
```

**Tasks in deno.json**:
```json
{
  "tasks": {
    "dev": "deno run --watch --allow-all main.ts",
    "start": "deno run --allow-net --allow-env main.ts"
  }
}
```

### Common Environment Variables
- `DENO_ENV`: `development` | `production`
- `PORT`: Application port (typically 8000)
- Custom app variables (accessed via `Deno.env.get()`)

### Import Maps and Dependencies

**deno.json imports**:
```json
{
  "imports": {
    "oak": "https://deno.land/x/oak@v12.6.1/mod.ts",
    "std/": "https://deno.land/std@0.208.0/"
  }
}
```

**deps.ts pattern**:
```typescript
// deps.ts - centralized dependency management
export { Application, Router } from "https://deno.land/x/oak/mod.ts";
export { load } from "https://deno.land/std/dotenv/mod.ts";
```

### APT Dependencies
**Default: `[]` (empty array)**

Reasons:
- Deno is a single binary with no external dependencies
- TypeScript runtime included
- No need for Node.js or npm
- Standard library is comprehensive

**Only include APT packages if**:
- Using FFI (Foreign Function Interface) with native libraries
- Explicitly required system tools

### TypeScript Configuration

Deno uses built-in TypeScript compiler. Optional `deno.json` config:
```json
{
  "compilerOptions": {
    "lib": ["deno.window"],
    "strict": true
  }
}
```

### Default Port
**8000** (Deno convention, different from Node.js 3000)

### Production Considerations
1. **No build step**: TypeScript runs directly (JIT or cached)
2. **Permissions**: Use minimal required permissions
3. **Caching**: Dependencies cached in `DENO_DIR`
4. **Binary compilation**: Optional with `deno compile` for standalone executable
5. **Environment**: Use `DENO_ENV=production` for app logic

### Fresh Framework Specifics

**File structure**:
```
project/
├── deno.json
├── fresh.gen.ts       # Auto-generated
├── routes/            # File-based routing
│   ├── index.tsx
│   └── api/
├── islands/           # Client-side interactive components
└── static/            # Static assets
```

**No build needed**: Fresh uses JIT rendering in both dev and prod

### Comparison with Node.js
| Feature | Deno | Node.js |
|---------|------|---------|
| TypeScript | Native | Requires transpilation |
| Permissions | Explicit flags | Full access |
| Package manager | URL imports | npm/yarn/pnpm |
| Standard library | Comprehensive | Minimal |
| Default port | 8000 | 3000 |

### Common Patterns
- **URL imports**: Dependencies from URLs (deno.land, esm.sh)
- **Top-level await**: Supported natively
- **No package.json**: Uses `deno.json` instead
- **Web standard APIs**: Fetch, Web Workers, etc.
- **Secure by default**: All permissions denied unless explicitly allowed
