# Node.js Examples

Complete examples of Node.js project analysis and deployment configuration.

## Example 1: React + Vite Frontend

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
      "setup": ["pnpm install"],
      "commands": ["pnpm run dev --host 0.0.0.0 --port 5173"]
    },
    "prod": {
      "environment": {
        "NODE_ENV": "production"
      },
      "setup": ["pnpm install", "pnpm run build"],
      "commands": ["npm run preview --host 0.0.0.0 --port 4173"]
    },
    "port": 4173,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export NODE_ENV=production\n    pnpm install\n    pnpm run build\n    pnpm run preview --host 0.0.0.0 --port 4173\nelse\n    pnpm install\n    pnpm run dev --host 0.0.0.0 --port 5173\nfi",
    "evidence": {
      "files": ["package.json", "vite.config.ts", "package-lock.json"],
      "reason": "Node.js project with Vite + React detected. Using pnpm as package manager. Vite preview server runs on port 4173 by default."
    }
  }
]
```

**Key Points**:
- Vite detected from `vite.config.ts` and dependencies
- Detect scripts in `package.json` for dev/prod commands
- Dev port: 5173, Prod port: 4173 (Vite defaults)
- APT packages: Empty array (Vite uses pure JavaScript tooling)
- Must pass `--host 0.0.0.0` to both dev and preview commands

---

## Example 2: Next.js Application

**User request**: "Analyze Next.js app"

**Detected files**:
- `package.json` with `next`, `react`, `react-dom`
- `next.config.js`
- `pages/` or `app/` directory
- `pnpm-lock.yaml`

**Output**:
```json
[
  {
    "language": "node",
    "version": "20.10.0",
    "apt": [],
    "dev": {
      "environment": {},
      "setup": ["pnpm install"],
      "commands": ["pnpm dev -H 0.0.0.0 -p 3000"]
    },
    "prod": {
      "environment": {
        "NODE_ENV": "production"
      },
      "setup": ["pnpm install", "pnpm build"],
      "commands": ["pnpm start -H 0.0.0.0 -p 3000"]
    },
    "port": 3000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export NODE_ENV=production\n    pnpm install\n    pnpm build\n    pnpm start -H 0.0.0.0 -p 3000\nelse\n    pnpm install\n    pnpm dev -H 0.0.0.0 -p 3000\nfi",
    "evidence": {
      "files": ["package.json", "next.config.js", "pnpm-lock.yaml"],
      "reason": "Next.js project detected. Using pnpm as package manager (detected from pnpm-lock.yaml). Next.js default port is 3000."
    }
  }
]
```

**Key Points**:
- Package manager: pnpm (detected from `pnpm-lock.yaml`)
- Next.js specific flags: `-H` for hostname
- Build required before production start
- Framework default port: 3000

---

## Example 3: Express.js API Server

**User request**: "Pack Express API"

**Detected files**:
- `package.json` with `express`
- `server.js` or `index.js`
- `yarn.lock`

**Output**:
```json
[
  {
    "language": "node",
    "version": "18.17.0",
    "apt": [],
    "dev": {
      "environment": {
        "PORT": "8080",
        "HOST": "0.0.0.0",
        "NODE_ENV": "development"
      },
      "setup": ["yarn install"],
      "commands": ["yarn dev"]
    },
    "prod": {
      "environment": {
        "PORT": "8080",
        "HOST": "0.0.0.0",
        "NODE_ENV": "production"
      },
      "setup": ["yarn install --production"],
      "commands": ["node server.js"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export PORT=8080\n    export HOST=0.0.0.0\n    export NODE_ENV=production\n    yarn install --production\n    node server.js\nelse\n    export PORT=8080\n    export HOST=0.0.0.0\n    export NODE_ENV=development\n    yarn install\n    yarn dev\nfi",
    "evidence": {
      "files": ["package.json", "server.js", "yarn.lock"],
      "reason": "Express.js API server detected. Using yarn as package manager. Port configured via environment variables."
    }
  }
]
```

**Key Points**:
- Package manager: yarn (detected from `yarn.lock`)
- Environment variables control host/port binding
- Production uses `--production` flag to skip dev dependencies
- Express apps typically use environment variables for configuration

---

## Node.js Detection Patterns

### Package Manager Detection
- `package-lock.json` → npm
- `yarn.lock` → yarn
- `pnpm-lock.yaml` → pnpm
- `bun.lockb` → bun

### Package Manager Command Syntax Differences

**CRITICAL: Use the correct syntax for each package manager.**

**Key Points**:
- **npm**: add `--` separator before arguments if needed
- **pnpm/yarn/bun**: NO `--` separator needed

### Framework Detection
| Framework | Detection | Default Port |
|-----------|-----------|--------------|
| Vite | `vite` in deps, `vite.config.*` | 5173 (dev), 4173 (prod) |
| Next.js | `next` in deps, `next.config.*` | 3000 |
| Nuxt | `nuxt` in deps, `nuxt.config.*` | 3000 |
| Express | `express` in deps | 3000/8080 |
| NestJS | `@nestjs/core` in deps | 3000 |
| Remix | `@remix-run/node` in deps | 3000 |

### Host Binding Syntax

**All servers must bind to `0.0.0.0` for external access.**

- **Vite**: Use `--host 0.0.0.0` flag
- **Next.js**: Use `-H 0.0.0.0` or `--hostname 0.0.0.0` flag
- **Express/Custom**: Use environment variables `HOST=0.0.0.0`

### Version Detection
1. `.node-version` file
2. `.nvmrc` file
3. `package.json` → `engines.node` field
4. Omit if uncertain

### Common Environment Variables
- `NODE_ENV`: `development` or `production`
- `PORT`: Application port number
- `HOST`: Bind address (always `0.0.0.0`)
- Framework-specific: `NEXT_PUBLIC_*`, `VITE_*`, `REACT_APP_*`

### APT Dependencies
**Default: `[]` (empty array)**

Reasons:
- Node.js base images include common native module dependencies
- `canvas`: Prebuilt binaries + runtime libs in base image
- `sharp`: Bundled libvips
- Most native modules use prebuilt binaries for glibc

**Only include APT packages if**:
- Project documentation explicitly requires them
- Dockerfile shows system package installation
