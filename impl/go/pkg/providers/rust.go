/**
 * DevBox Pack Execution Plan Generator - Rust Provider
 */

package providers

import (
	"regexp"
	"strings"

	"github.com/labring/devbox-pack/pkg/types"
)

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
			Priority: 85,
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
	indicators := []types.ConfidenceIndicator{
		{Weight: 40, Satisfied: p.HasFile(files, "Cargo.toml")},
		{Weight: 20, Satisfied: p.HasFile(files, "Cargo.lock")},
		{Weight: 20, Satisfied: p.HasAnyFile(files, []string{"*.rs"})},
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

	metadata := map[string]interface{}{
		"hasCargoToml": p.HasFile(files, "Cargo.toml"),
		"hasCargoLock": p.HasFile(files, "Cargo.lock"),
		"hasRustSrc":   p.HasAnyFile(files, []string{"*.rs"}),
		"hasLib":       p.HasFile(files, "src/lib.rs"),
		"hasMain":      p.HasFile(files, "src/main.rs"),
		"hasToolchain": p.HasAnyFile(files, []string{"rust-toolchain", "rust-toolchain.toml"}),
		"framework":    framework,
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
		reasons = append(reasons, "Cargo manifest (Cargo.toml)")
	}
	if p.HasFile(files, "Cargo.lock") {
		reasons = append(reasons, "Cargo lock file")
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
			return p.CreateVersionInfo(matches[1], "Cargo.toml rust-version"), nil
		}
	}

	// Default version
	return p.CreateVersionInfo("1.70", "default"), nil
}

// detectFramework detects framework
func (p *RustProvider) detectFramework(projectPath string, gitHandler interface{}) (string, error) {
	cargoToml, err := p.SafeReadText(projectPath, "Cargo.toml", gitHandler)
	if err != nil || cargoToml == "" {
		return "", nil
	}

	frameworkMap := map[string]string{
		"actix-web": "Actix Web",
		"axum":      "Axum",
		"warp":      "Warp",
		"rocket":    "Rocket",
		"tide":      "Tide",
		"hyper":     "Hyper",
		"tokio":     "Tokio",
		"async-std": "async-std",
		"serde":     "Serde",
		"diesel":    "Diesel",
		"sqlx":      "SQLx",
		"sea-orm":   "SeaORM",
		"tauri":     "Tauri",
		"yew":       "Yew",
		"leptos":    "Leptos",
		"dioxus":    "Dioxus",
	}

	for dependency, framework := range frameworkMap {
		if strings.Contains(cargoToml, dependency) {
			return framework, nil
		}
	}

	return "", nil
}

// GenerateCommands generates commands for Rust project
func (p *RustProvider) GenerateCommands(result *types.DetectResult, options types.CLIOptions) types.Commands {
	commands := types.Commands{}

	commands.Dev = []string{
		"cargo run",
	}

	commands.Build = []string{
		"cargo build --release",
	}

	commands.Start = []string{
		"./target/release/app",
	}

	return commands
}

// NeedsNativeCompilation checks if Rust project needs native compilation
func (p *RustProvider) NeedsNativeCompilation(result *types.DetectResult) bool {
	// Rust projects are compiled languages that need compilation
	// Here "native compilation" might refer to whether additional native dependencies are needed
	// For Rust projects, compilation is usually required, so return true
	return true
}
