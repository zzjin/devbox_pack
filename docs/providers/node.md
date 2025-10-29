# Node Provider

- Detection: Uses confidence-based detection with weighted indicators:
  - `package.json` (weight: 40)
  - Lock files: `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml`, `bun.lockb` (weight: 20)
  - `node_modules` directory (weight: 15)
  - Version files: `.nvmrc`, `.node-version` (weight: 10)
  - JavaScript/TypeScript files: `*.js`, `*.ts`, `*.mjs`, `*.cjs` (weight: 10)
  - Config files: `tsconfig.json`, `jsconfig.json` (weight: 5)
  - Detection threshold: 0.3 (30% confidence)

- Version Detection: Priority order for Node.js version resolution:
  1. `.nvmrc` file (supports `v` prefix)
  2. `.node-version` file
  3. `package.json` engines.node field
  4. Default version: `20`

- Framework Detection: Automatically detects popular Node.js frameworks by analyzing `package.json` dependencies:
  - **Frontend Frameworks**: Next.js, Nuxt.js, React, Vue.js, Angular, Svelte, Gatsby
  - **Backend Frameworks**: Express, Koa, Fastify, NestJS
  - **Build Tools**: Vite, Webpack, Parcel, Rollup
  - **Mobile**: Electron, React Native, Expo

- Package Manager Detection: Automatically detects package manager based on lock files:
  - `pnpm-lock.yaml` → pnpm
  - `yarn.lock` → yarn
  - `bun.lockb` → bun
  - `package-lock.json` → npm (default)

- Build Tool Detection: Identifies build tools from dependencies and build scripts:
  - **Build Tools**: vite, webpack, rollup, parcel, esbuild, turbo, nx
  - **Script Analysis**: Analyzes `package.json` build scripts for tool detection

- Commands:
  - **Development**: 
    - `<package-manager> install`
    - `<package-manager> run dev` (if dev script exists)
    - `<package-manager> run start` (if start script exists, fallback)
    - `node index.js` (final fallback)
  - **Build**: 
    - `<package-manager> install`
    - `<package-manager> run build` (if build script exists)
  - **Start**: 
    - `<package-manager> run start` (if start script exists)
    - `npx serve -s build` (for React projects after build)
    - `node index.js` (fallback)

- Metadata: Provides comprehensive metadata including:
  - `packageJson`: Package information (name, version, scripts, engines, type)
  - `hasTypeScript`: Presence of TypeScript configuration or files
  - `hasESM`: ES Module support detection
  - `hasCJS`: CommonJS support detection
  - `framework`: Detected framework name