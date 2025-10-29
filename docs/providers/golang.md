# Go Provider

- Detection: Uses confidence-based detection with weighted indicators:
  - `go.mod` (weight: 40)
  - `go.work` (weight: 30)
  - `*.go` files (weight: 25)
  - `go.sum` (weight: 15)
  - `main.go` or `cmd/` directory (weight: 10)
  - `vendor/` directory (weight: 5)
  - `Makefile` or `makefile` (weight: 5)
  - Detection threshold: 0.2 (20% confidence)

- Version Detection: Priority order for Go version resolution:
  1. `go X.XX` directive in `go.work` file
  2. `go X.XX` directive in `go.mod` file
  3. Version specified in `.go-version` file
  4. Default version: `1.21`

- Framework Detection: Automatically detects popular Go frameworks by analyzing `go.mod` dependencies:
  - **Web Frameworks**: Gin, Echo, Fiber, Gorilla Mux, Beego, Revel
  - **CLI Tools**: Cobra CLI, CLI
  - **Dependency Injection**: Fx
  - **ORM**: GORM

- Commands:
  - **Development**: `go mod download`, `go run .`
  - **Build**: `go mod download`, `go build -o app .`
  - **Start**: `./app`

- Metadata: Provides the following metadata fields:
  - `hasGoMod`: Presence of `go.mod` file
  - `hasGoSum`: Presence of `go.sum` file
  - `hasMainGo`: Presence of `main.go` file
  - `hasVendor`: Presence of `vendor/` directory
  - `hasGoWork`: Presence of `go.work` file
  - `hasMakefile`: Presence of `Makefile` or `makefile`
  - `framework`: Detected framework name (if any)