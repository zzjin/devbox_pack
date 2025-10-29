# Contributing to DevBox Pack

## Project Overview

DevBox Pack is a tool for generating execution plans for various development environments and frameworks. This Go implementation provides language detection, plan generation, and build configuration for containerized development environments.

## Pull Requests

We welcome pull requests that push the project forward in meaningful ways.
Please ensure your PRs:

- Address a specific problem or add a well-defined feature
- Include tests for new functionality
- Follow the existing code style and patterns
- Update documentation as needed

Note: We prefer focused, well-thought-out contributions over "drive-by" PRs that make superficial changes.

## Setup

### Prerequisites

- Go 1.21 or later
- Make (for using the Makefile)

### Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/labring/devbox-pack.git
   cd devbox-pack/impl/go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make build
   # or
   go build -o bin/devbox-pack ./cmd/main.go
   ```

4. Run tests:
   ```bash
   make test
   # or
   go test ./...
   ```

## Testing

### Unit Tests

Run unit tests with:
```bash
go test ./...
```

For verbose output:
```bash
go test -v ./...
```

For short tests only (excluding integration tests):
```bash
go test -short ./...
```

### Integration Tests

Integration tests use the testdata directory with real project examples:
```bash
go test ./pkg/service -v
```

### Test Coverage

Generate test coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Code Style

### Formatting

We use standard Go formatting. Before submitting a PR, run:
```bash
go fmt ./...
```

### Linting

We use golangci-lint for code quality checks:
```bash
golangci-lint run
```

### Vetting

Run Go vet to catch common issues:
```bash
go vet ./...
```

## Project Structure

```
├── cmd/                    # Main application entry point
├── pkg/
│   ├── cli/               # CLI interface and commands
│   ├── detector/          # Language and framework detection
│   ├── generators/        # Plan generation logic
│   ├── providers/         # Language-specific providers
│   ├── service/           # Core business logic
│   ├── types/             # Type definitions
│   └── utils/             # Utility functions
└── testdata/              # Test data and examples
```

## Adding New Language Support

To add support for a new programming language:

1. Create a new provider in `pkg/providers/`
2. Implement the `Provider` interface
3. Register the provider in `pkg/registry/registry.go`
4. Add detection logic in `pkg/detector/`
5. Add test cases in `testdata/`
6. Update documentation

## Useful Commands

Available Make targets:
- `make build` - Build the binary
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make lint` - Run linting
- `make fmt` - Format code

## Debugging

For debugging, you can use the built-in logging:
```bash
./bin/devbox-pack --debug <command>
```

Or set the log level:
```bash
LOG_LEVEL=debug ./bin/devbox-pack <command>
```

## Submitting Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-new-feature`
3. Make your changes and add tests
4. Ensure all tests pass: `make test`
5. Format your code: `make fmt`
6. Run linting: `make lint`
7. Commit your changes: `git commit -am 'Add some feature'`
8. Push to the branch: `git push origin feature/my-new-feature`
9. Submit a pull request

## Getting Help

If you have questions or need help:
- Check existing issues and discussions
- Create a new issue with detailed information
- Provide minimal reproduction steps when reporting bugs