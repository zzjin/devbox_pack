/**
 * DevBox Pack Execution Plan Generator - Rust Provider
 */

package providers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/labring/devbox-pack/pkg/types"
)

// RustWorkspaceInfo represents Cargo workspace information
type RustWorkspaceInfo struct {
	IsWorkspace bool     `json:"isWorkspace"`
	Members    []string `json:"members"`
	Excludes   []string `json:"excludes,omitempty"`
}

// RustProvider Rust project detector
type RustProvider struct {
	BaseProvider
}

// NewRustProvider creates Rust Provider
func NewRustProvider() *RustProvider {
	return &RustProvider{
		BaseProvider: BaseProvider{
			Name:     "rust",
			Language: "rust",
			Priority: 40,
		},
	}
}

// GetName gets Provider name
func (p *RustProvider) GetName() string {
	return p.Name
}

// GetPriority gets Provider priority
func (p *RustProvider) GetPriority() int {
	return p.Priority
}

// Detect detects Rust project
func (p *RustProvider) Detect(projectPath string, files []types.FileInfo, gitHandler interface{}) (*types.DetectResult, error) {
	// Check for workspace first
	isWorkspace := p.HasFile(files, "Cargo.toml") // Will check for [workspace] section later

	indicators := []types.ConfidenceIndicator{
		{Weight: 40, Satisfied: p.HasFile(files, "Cargo.toml")},
		{Weight: 20, Satisfied: p.HasFile(files, "Cargo.lock")},
		{Weight: 20, Satisfied: p.HasAnyFile(files, []string{"*.rs"})},
		{Weight: 15, Satisfied: p.HasFile(files, "Cargo.lock")}, // Higher weight for locked dependencies
		{Weight: 10, Satisfied: p.HasAnyFile(files, []string{"src/main.rs", "src/lib.rs"})},
		{Weight: 5, Satisfied: p.HasAnyFile(files, []string{"target/"})},
		{Weight: 5, Satisfied: p.HasFile(files, "rust-toolchain")},
	}

	confidence := p.CalculateConfidence(indicators)
	detected := confidence > 0.3

	if !detected {
		return p.CreateDetectResult(false, confidence, "", nil, "", "", "", nil, types.Evidence{}), nil
	}

	// Detect version
	version, err := p.detectRustVersion(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	// Detect framework
	framework, err := p.detectFramework(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	// Detect workspace and binary targets
	workspaceInfo, err := p.detectWorkspaceInfo(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	// Detect binary targets
	binaryTargets, err := p.detectBinaryTargets(projectPath, gitHandler)
	if err != nil {
		return nil, err
	}

	isWorkspace = workspaceInfo != nil && workspaceInfo.IsWorkspace

	metadata := map[string]interface{}{
		"hasCargoToml":  p.HasFile(files, "Cargo.toml"),
		"hasCargoLock":  p.HasFile(files, "Cargo.lock"),
		"hasRustSrc":    p.HasAnyFile(files, []string{"*.rs"}),
		"hasLib":        p.HasFile(files, "src/lib.rs"),
		"hasMain":       p.HasFile(files, "src/main.rs"),
		"hasToolchain":  p.HasAnyFile(files, []string{"rust-toolchain", "rust-toolchain.toml"}),
		"isWorkspace":   isWorkspace,
		"workspaceInfo": workspaceInfo,
		"binaryTargets": binaryTargets,
		"framework":     framework,
	}

	// Build Evidence
	evidence := types.Evidence{}

	// Collect key files
	var evidenceFiles []string
	if p.HasFile(files, "Cargo.toml") {
		evidenceFiles = append(evidenceFiles, "Cargo.toml")
	}
	if p.HasFile(files, "Cargo.lock") {
		evidenceFiles = append(evidenceFiles, "Cargo.lock")
	}
	if p.HasFile(files, "rust-toolchain") {
		evidenceFiles = append(evidenceFiles, "rust-toolchain")
	}
	if p.HasFile(files, "rust-toolchain.toml") {
		evidenceFiles = append(evidenceFiles, "rust-toolchain.toml")
	}
	if p.HasFile(files, "src/main.rs") {
		evidenceFiles = append(evidenceFiles, "src/main.rs")
	}
	if p.HasFile(files, "src/lib.rs") {
		evidenceFiles = append(evidenceFiles, "src/lib.rs")
	}

	evidence.Files = evidenceFiles

	// Build detection reason
	reason := "Detected Rust project based on: "
	var reasons []string
	if p.HasAnyFile(files, []string{"*.rs"}) {
		reasons = append(reasons, "Rust source files")
	}
	if p.HasFile(files, "Cargo.toml") {
		if isWorkspace {
			reasons = append(reasons, "Cargo workspace")
			if workspaceInfo != nil && len(workspaceInfo.Members) > 0 {
				reasons = append(reasons, fmt.Sprintf("%d members", len(workspaceInfo.Members)))
			}
		} else {
			reasons = append(reasons, "Cargo manifest (Cargo.toml)")
		}
	}
	if p.HasFile(files, "Cargo.lock") {
		reasons = append(reasons, "Cargo lock file")
	}
	if len(binaryTargets) > 0 {
		reasons = append(reasons, fmt.Sprintf("binary targets: %s", strings.Join(binaryTargets, ", ")))
	}
	if p.HasAnyFile(files, []string{"rust-toolchain", "rust-toolchain.toml"}) {
		reasons = append(reasons, "Rust toolchain configuration")
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

	return p.CreateDetectResult(
		true,
		confidence,
		"rust",
		version,
		framework,
		"cargo",
		"cargo",
		metadata,
		evidence,
	), nil
}

// detectRustVersion detects Rust version
func (p *RustProvider) detectRustVersion(projectPath string, gitHandler interface{}) (*types.VersionInfo, error) {
	// Read from rust-toolchain.toml
	toolchainToml, err := p.SafeReadText(projectPath, "rust-toolchain.toml", gitHandler)
	if err == nil && toolchainToml != "" {
		re := regexp.MustCompile(`channel\s*=\s*"([^"]+)"`)
		matches := re.FindStringSubmatch(toolchainToml)
		if len(matches) > 1 {
			return p.CreateVersionInfo(matches[1], "rust-toolchain.toml"), nil
		}
	}

	// Read from rust-toolchain
	toolchain, err := p.SafeReadText(projectPath, "rust-toolchain", gitHandler)
	if err == nil && toolchain != "" {
		version := strings.TrimSpace(toolchain)
		return p.CreateVersionInfo(version, "rust-toolchain"), nil
	}

	// Read MSRV (Minimum Supported Rust Version) from Cargo.toml
	cargoToml, err := p.SafeReadText(projectPath, "Cargo.toml", gitHandler)
	if err == nil && cargoToml != "" {
		re := regexp.MustCompile(`rust-version\s*=\s*"([^"]+)"`)
		matches := re.FindStringSubmatch(cargoToml)
		if len(matches) > 1 {
			return p.CreateVersionInfo(matches[1], "Cargo.toml"), nil
		}
	}

	// Default version
	return p.CreateVersionInfo("1.70", "default"), nil
}

// detectWorkspaceInfo detects Cargo workspace information
func (p *RustProvider) detectWorkspaceInfo(projectPath string, gitHandler interface{}) (*RustWorkspaceInfo, error) {
	cargoToml, err := p.SafeReadText(projectPath, "Cargo.toml", gitHandler)
	if err != nil || cargoToml == "" {
		return nil, nil
	}

	// Check for [workspace] section
	re := regexp.MustCompile(`\[workspace\]`)
	if !re.MatchString(cargoToml) {
		return nil, nil
	}

	workspaceInfo := &RustWorkspaceInfo{
		IsWorkspace: true,
		Members:    []string{},
	}

	// Extract workspace members
	membersRe := regexp.MustCompile(`members\s*=\s*\[([^\]]+)\]`)
	matches := membersRe.FindStringSubmatch(cargoToml)
	if len(matches) > 1 {
		membersStr := strings.ReplaceAll(matches[1], "\"", "")
		membersStr = strings.ReplaceAll(membersStr, " ", "")
		members := strings.Split(membersStr, ",")
		workspaceInfo.Members = members
	}

	// Extract workspace exclude patterns
	excludeRe := regexp.MustCompile(`exclude\s*=\s*\[([^\]]+)\]`)
	excludeMatches := excludeRe.FindStringSubmatch(cargoToml)
	if len(excludeMatches) > 1 {
		excludeStr := strings.ReplaceAll(excludeMatches[1], "\"", "")
		excludeStr = strings.ReplaceAll(excludeStr, " ", "")
		excludes := strings.Split(excludeStr, ",")
		workspaceInfo.Excludes = excludes
	}

	return workspaceInfo, nil
}

// detectBinaryTargets detects binary targets in Cargo.toml
func (p *RustProvider) detectBinaryTargets(projectPath string, gitHandler interface{}) ([]string, error) {
	cargoToml, err := p.SafeReadText(projectPath, "Cargo.toml", gitHandler)
	if err != nil || cargoToml == "" {
		return []string{}, nil
	}

	var binaries []string

	// Check for [[bin]] sections
	binRe := regexp.MustCompile(`\[\[bin\]\][\s\S]*?name\s*=\s*"([^"]+)"`)
	binMatches := binRe.FindAllStringSubmatch(cargoToml, -1)
	for _, match := range binMatches {
		if len(match) > 1 {
			binaries = append(binaries, match[1])
		}
	}

	// If no explicit binaries and has src/main.rs, use package name as binary
	if len(binaries) == 0 {
		packageRe := regexp.MustCompile(`name\s*=\s*"([^"]+)"`)
		packageMatches := packageRe.FindStringSubmatch(cargoToml)
		if len(packageMatches) > 1 {
			binaries = append(binaries, packageMatches[1])
		}
	}

	return binaries, nil
}

// detectFramework detects framework
func (p *RustProvider) detectFramework(projectPath string, gitHandler interface{}) (string, error) {
	cargoToml, err := p.SafeReadText(projectPath, "Cargo.toml", gitHandler)
	if err != nil || cargoToml == "" {
		return "", nil
	}

	// Priority-based framework detection - web frameworks first
	frameworkPriorities := []struct {
		dependency string
		framework  string
	}{
		{"actix-web", "Actix Web"},
		{"axum", "Axum"},
		{"warp", "Warp"},
		{"rocket", "Rocket"},
		{"tide", "Tide"},
		{"hyper", "Hyper"},
		{"tauri", "Tauri"},
		{"yew", "Yew"},
		{"leptos", "Leptos"},
		{"dioxus", "Dioxus"},
		{"tokio", "Tokio"},
		{"async-std", "async-std"},
		{"diesel", "Diesel"},
		{"sqlx", "SQLx"},
		{"sea-orm", "SeaORM"},
	}

	// Only consider actual frameworks, not utility crates
	for _, fp := range frameworkPriorities {
		if strings.Contains(cargoToml, fp.dependency) {
			return fp.framework, nil
		}
	}

	return "", nil
}

// GenerateCommands generates commands for Rust project
func (p *RustProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	// Check if this is a workspace
	isWorkspace := false
	var binaryTargets []string
	if result.Metadata != nil {
		if workspace, ok := result.Metadata["isWorkspace"].(bool); ok {
			isWorkspace = workspace
		}
		if binaries, ok := result.Metadata["binaryTargets"].([]string); ok {
			binaryTargets = binaries
		}
	}

	if isWorkspace {
		// Workspace commands
		commands.Setup = []string{"cargo build --workspace"}
		commands.Build = []string{"cargo build --release --workspace"}

		if len(binaryTargets) > 0 {
			// Run the first binary target
			commands.Dev = []string{"cargo run -p " + binaryTargets[0]}
			commands.Run = []string{"./target/release/" + binaryTargets[0]}
		} else {
			commands.Dev = []string{"cargo run --workspace"}
			commands.Run = []string{"cargo run --release --workspace"}
		}
	} else {
		// Single package commands
		commands.Setup = []string{"cargo build"}
		commands.Build = []string{"cargo build --release"}

		if len(binaryTargets) > 1 {
			// Multiple binary targets - provide commands for the first one
			commands.Dev = []string{"cargo run --bin " + binaryTargets[0]}
			commands.Run = []string{"./target/release/" + binaryTargets[0]}
		} else if len(binaryTargets) == 1 {
			// Single binary target
			commands.Dev = []string{"cargo run --bin " + binaryTargets[0]}
			commands.Run = []string{"./target/release/" + binaryTargets[0]}
		} else {
			// Default behavior
			commands.Dev = []string{"cargo run"}
			commands.Run = []string{"cargo run --release"}
		}
	}

	return commands
}

// GenerateEnvironment generates environment variables for Rust project
func (p *RustProvider) GenerateEnvironment(result *types.DetectResult) map[string]string {
	env := make(map[string]string)

	// Check if this is a workspace
	isWorkspace := false
	if result.Metadata != nil {
		if workspace, ok := result.Metadata["isWorkspace"].(bool); ok {
			isWorkspace = workspace
		}
	}

	// Set Rust specific environment variables
	env["RUST_ENV"] = "production"
	env["RUST_BACKTRACE"] = "1"
	env["CARGO_NET_GIT_FETCH_WITH_CLI"] = "true"

	// Set port for web applications
	env["PORT"] = "8080"

	// Add Rust version if available
	if result.Version != "" {
		env["RUST_VERSION"] = result.Version
	}

	// Add workspace-specific environment variables
	if isWorkspace {
		env["CARGO_WORKSPACE"] = "true"
		// Add workspace member count if available
		if result.Metadata != nil {
			if workspaceInfo, ok := result.Metadata["workspaceInfo"].(*RustWorkspaceInfo); ok && workspaceInfo != nil {
				env["CARGO_WORKSPACE_MEMBERS"] = fmt.Sprintf("%d", len(workspaceInfo.Members))
			}
		}
	}

	// Add framework-specific environment variables
	if strings.Contains(result.Framework, "Actix") {
		env["ACTIX_PROFILE"] = "production"
	} else if strings.Contains(result.Framework, "Axum") {
		env["AXUM_LOG_LEVEL"] = "info"
	} else if strings.Contains(result.Framework, "Warp") {
		env["WARP_LOG_LEVEL"] = "info"
	} else if strings.Contains(result.Framework, "Rocket") {
		env["ROCKET_PROFILE"] = "release"
		env["ROCKET_PORT"] = "8080"
	} else if strings.Contains(result.Framework, "Tokio") {
		env["TOKIO_LOG_LEVEL"] = "info"
	}

	return env
}

// NeedsNativeCompilation checks if Rust project needs native compilation
func (p *RustProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Rust projects are compiled languages that need compilation
	// Here "native compilation" might refer to whether additional native dependencies are needed
	// For Rust projects, compilation is usually required, so return true
	return true
}
