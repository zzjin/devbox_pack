# Rust Provider

- Detection: Uses confidence-based detection with weighted indicators:
  - Cargo manifest: `Cargo.toml` (weight: 40)
  - Lock file: `Cargo.lock` (weight: 20)
  - Rust source files: `*.rs` (weight: 20)
  - Entry points: `src/main.rs`, `src/lib.rs` (weight: 10)
  - Build artifacts: `target/` (weight: 5)
  - Toolchain config: `rust-toolchain` (weight: 5)
  - Detection threshold: 0.3 (30% confidence)

- Version Detection: Priority order for Rust version resolution:
  1. `rust-toolchain.toml` channel field
  2. `rust-toolchain` file content
  3. `Cargo.toml` rust-version field (MSRV)
  4. Default version: `1.70`

- Framework Detection: Automatically detects popular Rust frameworks and libraries by analyzing Cargo.toml:
  - **Web Frameworks**: Actix Web, Axum, Warp, Rocket, Tide, Hyper
  - **Async Runtime**: Tokio, async-std
  - **Serialization**: Serde
  - **Database**: Diesel, SQLx, SeaORM
  - **Desktop/Mobile**: Tauri
  - **Frontend**: Yew, Leptos, Dioxus
  - **Analysis**: Checks `Cargo.toml` dependencies section

- Package Manager: Always uses Cargo for dependency management

- Commands:
  - **Development**: 
    - `cargo run`
  - **Build**: 
    - `cargo build --release`
  - **Start**: 
    - `./target/release/app`

- Native Compilation Detection: Rust projects require compilation
  - Returns `true` as Rust is a compiled language requiring build step

- Metadata: Provides comprehensive metadata including:
  - `hasCargoToml`: Presence of `Cargo.toml`
  - `hasCargoLock`: Presence of `Cargo.lock`
  - `hasRustSrc`: Presence of Rust source files
  - `hasLib`: Presence of `src/lib.rs`
  - `hasMain`: Presence of `src/main.rs`
  - `hasToolchain`: Presence of toolchain configuration files
  - `framework`: Detected framework name