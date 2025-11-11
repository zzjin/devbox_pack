/**
 * DevBox Pack Execution Plan Generator - Node.js Provider
 */

package providers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/labring/devbox-pack/pkg/git"
	"github.com/labring/devbox-pack/pkg/types"
)

// NodeProvider Node.js project detector
type NodeProvider struct {
	BaseProvider
}

// NewNodeProvider creates Node.js Provider
func NewNodeProvider() *NodeProvider {
	return &NodeProvider{
		BaseProvider: BaseProvider{
			Name:     "node",
			Language: "node",
			Priority: 80,
		},
	}
}

// Detect detects if project uses Node.js
func (np *NodeProvider) Detect(
	projectPath string,
	files []types.FileInfo,
	gitHandler interface{},
) (*types.DetectResult, error) {
	gh := gitHandler.(*git.GitHandler)
	indicators := []types.ConfidenceIndicator{
		{Weight: 40, Satisfied: np.HasFile(files, "package.json")},
		{Weight: 20, Satisfied: np.HasAnyFile(files, []string{"package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb"})},
		{Weight: 15, Satisfied: np.HasFile(files, "node_modules")},
		{Weight: 10, Satisfied: np.HasAnyFile(files, []string{".nvmrc", ".node-version"})},
		{Weight: 10, Satisfied: np.HasAnyFile(files, []string{"*.js", "*.ts", "*.mjs", "*.cjs"})},
		{Weight: 5, Satisfied: np.HasAnyFile(files, []string{"tsconfig.json", "jsconfig.json"})},
	}

	confidence := np.CalculateConfidence(indicators)
	detected := confidence > 0.3

	if !detected {
		return np.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Parse package.json
	packageJSON, _ := np.SafeReadJSON(projectPath, "package.json", gh)

	// Detect version
	version, err := np.detectNodeVersion(projectPath, gh)
	if err != nil {
		// Use default version
		version = np.CreateVersionInfo("20", "default")
	}

	// Detect framework
	framework := np.detectFramework(packageJSON)

	// Detect package manager
	lockFiles := map[string]string{
		"pnpm-lock.yaml":    "pnpm",
		"yarn.lock":         "yarn",
		"bun.lockb":         "bun",
		"package-lock.json": "npm",
	}
	packageManager := np.DetectPackageManager(files, lockFiles)
	if packageManager == "" {
		packageManager = "npm"
	}

	// Detect build tool
	buildTool := np.detectBuildTool(packageJSON)

	// Build metadata
	metadata := make(map[string]interface{})
	if packageJSON != nil {
		pkgInfo := make(map[string]interface{})
		if name, ok := packageJSON["name"].(string); ok {
			pkgInfo["name"] = name
		}
		if version, ok := packageJSON["version"].(string); ok {
			pkgInfo["version"] = version
		}
		if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
			pkgInfo["scripts"] = scripts
		}
		if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
			pkgInfo["engines"] = engines
		}
		if pkgType, ok := packageJSON["type"].(string); ok {
			pkgInfo["type"] = pkgType
		}
		metadata["packageJson"] = pkgInfo
	}

	metadata["hasTypeScript"] = np.HasAnyFile(files, []string{"tsconfig.json", "*.ts", "*.tsx"})

	hasESM := false
	if packageJSON != nil {
		if pkgType, ok := packageJSON["type"].(string); ok && pkgType == "module" {
			hasESM = true
		}
	}
	if !hasESM {
		hasESM = np.HasAnyFile(files, []string{"*.mjs"})
	}
	metadata["hasESM"] = hasESM

	hasCJS := true
	if packageJSON != nil {
		if pkgType, ok := packageJSON["type"].(string); ok && pkgType == "module" {
			hasCJS = false
		}
	}
	if !hasCJS {
		hasCJS = np.HasAnyFile(files, []string{"*.cjs"})
	}
	metadata["hasCJS"] = hasCJS

	// Add framework information to metadata
	metadata["framework"] = framework

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	if np.HasFile(files, "package.json") {
		evidenceFiles = append(evidenceFiles, "package.json")
	}
	if np.HasFile(files, "package-lock.json") {
		evidenceFiles = append(evidenceFiles, "package-lock.json")
	}
	if np.HasFile(files, "yarn.lock") {
		evidenceFiles = append(evidenceFiles, "yarn.lock")
	}
	if np.HasFile(files, "pnpm-lock.yaml") {
		evidenceFiles = append(evidenceFiles, "pnpm-lock.yaml")
	}
	if np.HasFile(files, "bun.lockb") {
		evidenceFiles = append(evidenceFiles, "bun.lockb")
	}
	if np.HasFile(files, "tsconfig.json") {
		evidenceFiles = append(evidenceFiles, "tsconfig.json")
	}
	if np.HasFile(files, ".nvmrc") {
		evidenceFiles = append(evidenceFiles, ".nvmrc")
	}
	if np.HasFile(files, ".node-version") {
		evidenceFiles = append(evidenceFiles, ".node-version")
	}
	// Add framework-specific configuration files
	if np.HasFile(files, "svelte.config.js") {
		evidenceFiles = append(evidenceFiles, "svelte.config.js")
	}
	if np.HasFile(files, "svelte.config.ts") {
		evidenceFiles = append(evidenceFiles, "svelte.config.ts")
	}
	if np.HasFile(files, "svelte.config.mjs") {
		evidenceFiles = append(evidenceFiles, "svelte.config.mjs")
	}
	if np.HasFile(files, "astro.config.js") {
		evidenceFiles = append(evidenceFiles, "astro.config.js")
	}
	if np.HasFile(files, "astro.config.ts") {
		evidenceFiles = append(evidenceFiles, "astro.config.ts")
	}
	if np.HasFile(files, "astro.config.mjs") {
		evidenceFiles = append(evidenceFiles, "astro.config.mjs")
	}

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected Node.js project based on: "
	var reasons []string
	if np.HasFile(files, "package.json") {
		reasons = append(reasons, "package.json")
	}
	if packageManager != "npm" {
		reasons = append(reasons, packageManager+" lock file")
	}
	if framework != "" {
		reasons = append(reasons, "framework: "+framework)
	}
	if len(reasons) > 0 {
		reason += reasons[0]
		for i := 1; i < len(reasons); i++ {
			reason += ", " + reasons[i]
		}
	}
	evidence.Reason = reason

	return np.CreateDetectResult(
		true,
		confidence,
		"node",
		version,
		framework,
		packageManager,
		buildTool,
		metadata,
		evidence,
	), nil
}

// GenerateCommands generates commands for Node.js project
func (np *NodeProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	// Determine package manager
	packageManager := "npm"
	if np.HasFileInEvidence(result.Evidence.Files, "yarn.lock") {
		packageManager = "yarn"
	} else if np.HasFileInEvidence(result.Evidence.Files, "pnpm-lock.yaml") {
		packageManager = "pnpm"
	} else if np.HasFileInEvidence(result.Evidence.Files, "bun.lockb") {
		packageManager = "bun"
	}

	// Setup commands - install dependencies
	commands.Setup = []string{
		fmt.Sprintf("%s install", packageManager),
	}

	// Development commands
	if np.hasScript(result, "dev") {
		commands.Dev = []string{fmt.Sprintf("%s run dev", packageManager)}
	} else if np.hasScript(result, "start") {
		commands.Dev = []string{fmt.Sprintf("%s run start", packageManager)}
	} else {
		commands.Dev = []string{"node index.js"}
	}

	// Build commands
	if np.hasScript(result, "build") {
		buildCommands := []string{fmt.Sprintf("%s run build", packageManager)}

		// Add TypeScript-specific build steps for SvelteKit
		if np.isSvelteKitProject(result) && result.Metadata != nil {
			if hasTypeScript, ok := result.Metadata["hasTypeScript"].(bool); ok && hasTypeScript {
				// Prepend type checking before build
				buildCommands = append([]string{
					fmt.Sprintf("%s run check", packageManager), // Run SvelteKit type checking
				}, buildCommands...)
			}
		}

		commands.Build = buildCommands
	}

	// Run commands - prioritize framework-specific logic
	if np.isSvelteKitProject(result) {
		// SvelteKit projects use "preview" script for production
		commands.Run = []string{fmt.Sprintf("%s run preview", packageManager)}
	} else if np.isAstroProject(result) {
		// Astro projects use "preview" script for production
		commands.Run = []string{fmt.Sprintf("%s run preview", packageManager)}
	} else if np.hasScript(result, "start") {
		commands.Run = []string{fmt.Sprintf("%s run start", packageManager)}
	} else if len(commands.Build) > 0 {
		// If there are build steps, assume it creates dist or build directory
		if np.isNextJSProject(result) {
			commands.Run = []string{fmt.Sprintf("%s run start", packageManager)}
		} else if np.isReactProject(result) {
			commands.Run = []string{"npx serve -s build"}
		} else {
			commands.Run = []string{"node index.js"}
		}
	} else {
		commands.Run = []string{"node index.js"}
	}

	return commands
}

// GenerateEnvironment generates environment variables for Node.js project
func (np *NodeProvider) GenerateEnvironment(result *types.DetectResult) map[string]string {
	env := make(map[string]string)

	// Set Node.js specific environment variables
	env["NODE_ENV"] = "development"
	env["NPM_CONFIG_FUND"] = "false"
	env["NPM_CONFIG_AUDIT"] = "false"

	// Add framework-specific environment variables
	if np.isSvelteKitProject(result) {
		// SvelteKit-specific variables
		env["SVELTEKIT_ADAPTER_NAME"] = "auto"
		env["SVELTEKIT_PRERENDER"] = "auto"

		// Set correct port for SvelteKit (default 5173)
		env["PORT"] = "5173"
	} else if np.isAstroProject(result) {
		// Astro-specific variables
		env["ASTRO_ADAPTER"] = "auto"

		// Set correct port for Astro (default 4321)
		env["PORT"] = "4321"
	} else {
		// Set port if not specified
		env["PORT"] = "3000"
	}

	// Add TypeScript specific variables if TypeScript is detected
	if result.Metadata != nil {
		if hasTypeScript, ok := result.Metadata["hasTypeScript"].(bool); ok && hasTypeScript {
			env["TS_NODE_PROJECT"] = "tsconfig.json"
			// For SvelteKit with TypeScript
			if np.isSvelteKitProject(result) {
				env["SVELTEKIT TypeScript"] = "true"
			}
		}
	}

	// Add Node.js version if available
	if result.Version != "" {
		env["NODE_VERSION"] = result.Version
	}

	return env
}

// NeedsNativeCompilation checks if Node.js project needs native compilation
func (np *NodeProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Check metadata for native modules flag
	if result.Metadata != nil {
		if hasNativeModules, ok := result.Metadata["hasNativeModules"].(bool); ok {
			return hasNativeModules
		}
	}

	// Check native modules
	nativeModules := []string{
		"node-gyp", "node-sass", "sharp", "sqlite3",
		"bcrypt", "canvas", "puppeteer", "fsevents",
		"fibers", "grpc", "leveldown", "node-expat",
	}

	evidenceStr := strings.Join(result.Evidence.Files, " ")
	for _, module := range nativeModules {
		if strings.Contains(evidenceStr, module) {
			return true
		}
	}

	// Check dependencies in package.json (if available)
	if metadata, ok := result.Metadata["packageJson"].(map[string]interface{}); ok {
		// Can further check dependencies and devDependencies here
		_ = metadata // Not implementing detailed check for now
	}

	return false
}

// Helper methods

// hasScript checks if package.json has specific script
func (np *NodeProvider) hasScript(result *types.DetectResult, script string) bool {
	// Check if we have package.json in evidence files
	if !np.HasFileInEvidence(result.Evidence.Files, "package.json") {
		return false
	}

	// For SvelteKit projects, we know they have specific scripts
	if np.isSvelteKitProject(result) {
		switch script {
		case "dev", "build", "preview", "check":
			return true // SvelteKit projects typically have these scripts
		}
	}

	// For Astro projects, we know they have specific scripts
	if np.isAstroProject(result) {
		switch script {
		case "dev", "build", "preview":
			return true // Astro projects typically have these scripts
		}
	}

	// For other projects, check common scripts
	switch script {
	case "dev", "build", "start":
		return true // Most modern projects have these scripts
	}

	return false
}

// isNextJSProject checks if it's a Next.js project
func (np *NodeProvider) isNextJSProject(result *types.DetectResult) bool {
	return np.HasFileInEvidence(result.Evidence.Files, "next.config.js") ||
		np.HasFileInEvidence(result.Evidence.Files, "next.config.ts") ||
		np.HasFileInEvidence(result.Evidence.Files, "next.config.mjs")
}

// isReactProject checks if it's a React project
func (np *NodeProvider) isReactProject(result *types.DetectResult) bool {
	// This needs to check React dependencies in package.json
	return np.HasFileInEvidence(result.Evidence.Files, "package.json")
}

// isSvelteKitProject checks if it's a SvelteKit project
func (np *NodeProvider) isSvelteKitProject(result *types.DetectResult) bool {
	return np.HasFileInEvidence(result.Evidence.Files, "svelte.config.js") ||
		np.HasFileInEvidence(result.Evidence.Files, "svelte.config.ts") ||
		np.HasFileInEvidence(result.Evidence.Files, "svelte.config.mjs")
}

// isAstroProject checks if it's an Astro project
func (np *NodeProvider) isAstroProject(result *types.DetectResult) bool {
	return np.HasFileInEvidence(result.Evidence.Files, "astro.config.js") ||
		np.HasFileInEvidence(result.Evidence.Files, "astro.config.ts") ||
		np.HasFileInEvidence(result.Evidence.Files, "astro.config.mjs")
}

// detectNodeVersion detects Node.js version
func (np *NodeProvider) detectNodeVersion(projectPath string, gitHandler interface{}) (*types.VersionInfo, error) {
	// Read from .nvmrc
	nvmrcPattern := regexp.MustCompile(`^v?(.+)$`)
	if version, err := np.ParseVersionFromText(projectPath, ".nvmrc", gitHandler, nvmrcPattern); err == nil {
		return np.CreateVersionInfo(np.NormalizeVersion(version), ".nvmrc"), nil
	}

	// Read from .node-version
	if version, err := np.ParseVersionFromText(projectPath, ".node-version", gitHandler, nvmrcPattern); err == nil {
		return np.CreateVersionInfo(np.NormalizeVersion(version), ".node-version"), nil
	}

	// Read from package.json engines
	if packageJSON, err := np.SafeReadJSON(projectPath, "package.json", gitHandler); err == nil {
		if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
			if nodeVersion, ok := engines["node"].(string); ok {
				return np.CreateVersionInfo(np.NormalizeVersion(nodeVersion), "package.json engines"), nil
			}
		}
	}

	// Default version
	return np.CreateVersionInfo("20", "default"), nil
}

// detectFramework detects framework with SvelteKit priority
func (np *NodeProvider) detectFramework(packageJSON map[string]interface{}) string {
	if packageJSON == nil {
		return ""
	}

	// Collect all dependencies
	allDeps := make(map[string]interface{})
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		for k, v := range deps {
			allDeps[k] = v
		}
	}
	if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		for k, v := range devDeps {
			allDeps[k] = v
		}
	}
	if peerDeps, ok := packageJSON["peerDependencies"].(map[string]interface{}); ok {
		for k, v := range peerDeps {
			allDeps[k] = v
		}
	}

	// Priority framework detection - check for specific frameworks first
	if _, exists := allDeps["@sveltejs/kit"]; exists {
		return "sveltekit"
	}
	if _, exists := allDeps["astro"]; exists {
		return "astro"
	}

	// Then check other frameworks in order of specificity
	frameworkPriority := []string{
		"next",
		"nuxt",
		"@angular/core",
		"gatsby",
		"nestjs",
		"@nestjs/core",
		"express",
		"koa",
		"fastify",
		"svelte",
		"react",
		"vue",
		"vite",
		"webpack",
		"parcel",
		"rollup",
		"electron",
		"react-native",
		"expo",
	}

	for _, depName := range frameworkPriority {
		if _, exists := allDeps[depName]; exists {
			switch depName {
			case "next":
				return "next"
			case "nuxt":
				return "nuxt"
			case "@angular/core":
				return "angular"
			case "gatsby":
				return "gatsby"
			case "nestjs", "@nestjs/core":
				return "nestjs"
			case "express":
				return "express"
			case "koa":
				return "koa"
			case "fastify":
				return "fastify"
			case "svelte":
				return "svelte"
			case "react":
				return "react"
			case "vue":
				return "vue"
			case "vite":
				return "vite"
			case "webpack":
				return "webpack"
			case "parcel":
				return "parcel"
			case "rollup":
				return "rollup"
			case "electron":
				return "electron"
			case "react-native":
				return "react-native"
			case "expo":
				return "expo"
			}
		}
	}

	return ""
}

// detectBuildTool detects build tool
func (np *NodeProvider) detectBuildTool(packageJSON map[string]interface{}) string {
	if packageJSON == nil {
		return ""
	}

	allDeps := make(map[string]interface{})

	// Merge all dependencies
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		for k, v := range deps {
			allDeps[k] = v
		}
	}
	if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		for k, v := range devDeps {
			allDeps[k] = v
		}
	}

	// Check build tools in dependencies
	buildTools := []string{"vite", "webpack", "rollup", "parcel", "esbuild", "turbo", "nx"}
	for _, tool := range buildTools {
		if _, exists := allDeps[tool]; exists {
			return tool
		}
	}

	// Check build tools in scripts
	if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
		if buildScript, ok := scripts["build"].(string); ok {
			scriptTools := []string{"vite", "webpack", "rollup", "parcel", "esbuild", "tsc"}
			for _, tool := range scriptTools {
				if regexp.MustCompile(tool).MatchString(buildScript) {
					return tool
				}
			}
		}
	}

	return ""
}
