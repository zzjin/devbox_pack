# Rust Examples

Complete examples of Rust project analysis and deployment configuration.

## Example 1: Rust Web Service with Actix-Web

**User request**: "Pack Rust API"

**Detected files**:
- `Cargo.toml` with `actix-web`
- `src/main.rs`
- `Cargo.lock`

**Output**:
```json
[
  {
    "language": "rust",
    "version": "1.75",
    "apt": [],
    "dev": {
      "environment": {
        "RUST_LOG": "debug"
      },
      "commands": ["cargo run"]
    },
    "prod": {
      "environment": {
        "RUST_LOG": "info"
      },
      "setup": ["cargo build --release"],
      "commands": ["./target/release/app"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export RUST_LOG=info\n    cargo build --release\n    ./target/release/app\nelse\n    export RUST_LOG=debug\n    cargo run\nfi",
    "evidence": {
      "files": ["Cargo.toml", "src/main.rs", "Cargo.lock"],
      "reason": "Rust web service with actix-web detected. Building optimized release binary for production."
    }
  }
]
```

**Key Points**:
- Actix-web is a popular Rust web framework
- Dev: `cargo run` for quick compilation and execution
- Prod: `cargo build --release` for optimized binary
- Release binaries are **significantly faster** than debug builds
- `RUST_LOG` controls logging verbosity

---

## Example 2: Rust API with Axum Framework

**User request**: "Analyze Axum project"

**Detected files**:
- `Cargo.toml` with `axum` and `tokio`
- `src/main.rs`
- `.cargo/config.toml`

**Output**:
```json
[
  {
    "language": "rust",
    "version": "1.76",
    "apt": [],
    "dev": {
      "environment": {
        "RUST_LOG": "debug,axum=trace",
        "PORT": "8080"
      },
      "commands": ["cargo watch -x run"]
    },
    "prod": {
      "environment": {
        "RUST_LOG": "info",
        "PORT": "8080"
      },
      "setup": ["cargo build --release"],
      "commands": ["./target/release/app"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export RUST_LOG=info\n    export PORT=8080\n    cargo build --release\n    ./target/release/app\nelse\n    export RUST_LOG=\"debug,axum=trace\"\n    export PORT=8080\n    cargo watch -x run\nfi",
    "evidence": {
      "files": ["Cargo.toml", "src/main.rs"],
      "reason": "Axum web framework detected with tokio async runtime. Using cargo-watch for hot reload in development."
    }
  }
]
```

**Key Points**:
- Axum is a modern ergonomic web framework
- Dev: `cargo watch -x run` for hot reload (requires cargo-watch)
- Granular logging: `RUST_LOG="debug,axum=trace"`
- Async runtime: tokio (included in dependencies)

---

## Example 3: Rust Microservice with Rocket

**User request**: "Pack Rocket web app"

**Detected files**:
- `Cargo.toml` with `rocket`
- `Rocket.toml` (config file)
- `src/main.rs`

**Output**:
```json
[
  {
    "language": "rust",
    "version": "1.75",
    "apt": [],
    "dev": {
      "environment": {
        "ROCKET_ENV": "development",
        "ROCKET_ADDRESS": "0.0.0.0",
        "ROCKET_PORT": "8080"
      },
      "commands": ["cargo run"]
    },
    "prod": {
      "environment": {
        "ROCKET_ENV": "production",
        "ROCKET_ADDRESS": "0.0.0.0",
        "ROCKET_PORT": "8080"
      },
      "setup": ["cargo build --release"],
      "commands": ["./target/release/app"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export ROCKET_ENV=production\n    export ROCKET_ADDRESS=0.0.0.0\n    export ROCKET_PORT=8080\n    cargo build --release\n    ./target/release/app\nelse\n    export ROCKET_ENV=development\n    export ROCKET_ADDRESS=0.0.0.0\n    export ROCKET_PORT=8080\n    cargo run\nfi",
    "evidence": {
      "files": ["Cargo.toml", "Rocket.toml", "src/main.rs"],
      "reason": "Rocket web framework detected. Using Rocket environment configuration."
    }
  }
]
```

**Key Points**:
- Rocket uses `ROCKET_*` environment variables
- `ROCKET_ENV`: development | production
- Configuration can be in `Rocket.toml` or environment variables
- Always set `ROCKET_ADDRESS=0.0.0.0` for container deployments

---

## Rust Detection Patterns

### Required Files
- `Cargo.toml` (REQUIRED - Rust's manifest file)
- `Cargo.lock` (dependency lock file)
- `src/` directory with Rust source files (`*.rs`)

### Framework Detection
| Framework | Detection | Use Case | Default Port |
|-----------|-----------|----------|--------------|
| Actix-web | `actix-web` in dependencies | High-performance web | 8080 |
| Axum | `axum` in dependencies | Modern ergonomic web | 8080 |
| Rocket | `rocket` in dependencies | Feature-rich web | 8080 |
| Warp | `warp` in dependencies | Filter-based web | 8080 |
| Tide | `tide` in dependencies | Async web | 8080 |
| Hyper | `hyper` in dependencies | Low-level HTTP | 8080 |

### Host Binding Patterns

**Most frameworks** bind to `0.0.0.0` by default or via environment:

```rust
// Actix-web
HttpServer::new(|| App::new())
    .bind("0.0.0.0:8080")?
    .run()
    .await

// Axum
axum::Server::bind(&"0.0.0.0:8080".parse()?)
    .serve(app.into_make_service())
    .await

// Rocket
// Set via ROCKET_ADDRESS environment variable
```

### Version Detection
1. `Cargo.toml` → `package.rust-version` field
2. `rust-toolchain.toml` file
3. `.rust-version` file
4. Omit if uncertain

### Cargo Commands

**Development**:
```bash
cargo run                # Compile and run
cargo watch -x run       # Auto-reload on file changes (requires cargo-watch)
cargo run --bin <name>   # Run specific binary
```

**Production**:
```bash
cargo build --release    # Optimized release build
./target/release/app     # Run binary
```

**Optimization Levels**:
- Debug build: Fast compilation, slow runtime
- Release build: Slow compilation, **very fast** runtime (10-100x faster)

### Common Environment Variables
- `RUST_LOG`: Logging level
  - `error` | `warn` | `info` | `debug` | `trace`
  - Module-specific: `RUST_LOG="debug,actix_web=trace"`
- `RUST_BACKTRACE`: `1` for backtrace on panic
- Framework-specific:
  - `ROCKET_ADDRESS`, `ROCKET_PORT`, `ROCKET_ENV`
  - `ACTIX_*` (less common, usually configured in code)
  - Custom env vars defined in application

### Build Output Locations
- **Debug**: `target/debug/` (default with `cargo run` or `cargo build`)
- **Release**: `target/release/` (with `--release` flag)
- Binary name matches package name in `Cargo.toml`

### APT Dependencies
**Default: `[]` (empty array)**

Reasons:
- Rust compiles to **static binaries** by default
- No runtime dependencies for pure Rust code
- All dependencies linked at compile time
- musl-based builds are fully static

**Only include APT packages if**:
- Using system libraries (rare): `libssl-dev`, `pkg-config`
- Native bindings to C libraries
- Database drivers requiring system libs

**Note**: Most Rust crates use pure Rust implementations or include static libs

### Async Runtimes
| Runtime | Detection | Notes |
|---------|-----------|-------|
| Tokio | `tokio` in dependencies | Most popular, full-featured |
| async-std | `async-std` in dependencies | Alternative to tokio |
| smol | `smol` in dependencies | Lightweight runtime |

### Development Tools

**cargo-watch** (hot reload):
```bash
cargo install cargo-watch
cargo watch -x run              # Watch and run
cargo watch -x test             # Watch and test
```

**cargo-edit** (dependency management):
```bash
cargo install cargo-edit
cargo add axum                  # Add dependency
cargo upgrade                   # Upgrade dependencies
```

### Default Port
**8080** (common for Rust web services)

### Production Best Practices
1. **Always use `--release`**: Release builds are dramatically faster
2. **Strip binaries**: Reduce size with `strip` or `cargo-strip`
3. **Logging**: Use `RUST_LOG=info` in production
4. **Static linking**: Default for most Rust apps (no runtime deps)
5. **Security**: Update dependencies regularly with `cargo audit`
6. **Binary size**: Consider `cargo-bloat` for size analysis

### Binary Optimization (Optional)

Add to `Cargo.toml` for smaller binaries:
```toml
[profile.release]
opt-level = "z"        # Optimize for size
lto = true             # Link-time optimization
codegen-units = 1      # Single codegen unit
strip = true           # Strip symbols
```

### Cross-Compilation
Rust excels at cross-compilation:
```bash
# For Linux x86_64 (most common)
cargo build --release --target x86_64-unknown-linux-gnu

# For musl (fully static)
cargo build --release --target x86_64-unknown-linux-musl
```

### Common Patterns
- **Zero-copy parsing**: Rust frameworks often use zero-copy deserialization
- **Type-safe routing**: Compile-time route verification
- **Async by default**: Most modern frameworks use async/await
- **Safety guarantees**: Memory safety without garbage collection
