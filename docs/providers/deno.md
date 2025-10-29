# Deno Provider

- Detection: Uses confidence-based detection with weighted indicators:
  - Deno configuration: `deno.json`, `deno.jsonc` (weight: 40)
  - Source files: `*.ts`, `*.js` (weight: 25)
  - Lock file: `deno.lock` (weight: 15)
  - Module files: `deps.ts`, `mod.ts` (weight: 10)
  - Import map: `import_map.json` (weight: 5)
  - Entry points: `main.ts`, `app.ts`, `server.ts` (weight: 5)
  - Detection threshold: 0.2 (20% confidence)
  - **Exclusions**: Skips detection if `Staticfile` or `go.work` exists

- Version Detection: Priority order for Deno version resolution:
  1. `deno.json` version field
  2. `deno.jsonc` version field
  3. Default version: `1.40`

- Framework Detection: Automatically detects popular Deno frameworks by analyzing configuration and dependency files:
  - **Web Frameworks**: Fresh, Oak, Hono, Aleph.js, Ultra
  - **Analysis**: Checks `deno.json`, `deno.jsonc` imports section and `deps.ts` file for framework imports

- Package Manager: Uses Deno's built-in dependency management (no external package manager)

- Commands:
  - **Development**: 
    - `deno run --allow-all <main_file>`
  - **Build**: 
    - No build command (runtime interpreted)
  - **Start**: 
    - `deno run --allow-all <main_file>`
  - **Main File Priority**: `mod.ts` → `index.ts` → `app.ts` → `main.ts` (default)

- Native Compilation Detection: Deno projects don't require native compilation
  - Returns `false` as Deno is runtime interpreted

- Metadata: Provides comprehensive metadata including:
  - `hasDenoJson`: Presence of `deno.json`
  - `hasDenoJsonc`: Presence of `deno.jsonc`
  - `hasDenoLock`: Presence of `deno.lock`
  - `hasImportMap`: Presence of `import_map.json`
  - `hasDepsTs`: Presence of `deps.ts`
  - `hasModTs`: Presence of `mod.ts`
  - `hasTypeScript`: Presence of TypeScript files