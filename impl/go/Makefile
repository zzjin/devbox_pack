# Production-ready Makefile for devbox-pack Go implementation (Linux/AMD64 only)
.PHONY: all build run clean clean-bin install fmt lint test deps tidy help

# Project configuration
PROJECT_NAME := devbox-pack
MODULE_NAME := github.com/labring/devbox-pack
OUTPUT_DIR := bin
MAIN_FILE := cmd/main.go

# Target platform (Linux/AMD64 only)
GOOS := linux
GOARCH := amd64

# Version information from git
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null || echo "")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null | tr ' ' '_' || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
VERSION := $(if $(GIT_TAG),$(GIT_TAG),$(if $(GIT_COMMIT),dev-$(GIT_COMMIT),1.0.0))

# Go configuration
GO_VERSION := $(shell go version | cut -d' ' -f3)

# Optimized build flags (maximum size optimization)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT) -X main.branch=$(GIT_BRANCH) -X main.buildTime=$(BUILD_TIME) -X main.goVersion=$(GO_VERSION) -s -w -extldflags '-static'
BUILD_FLAGS := -ldflags "$(LDFLAGS)" -trimpath -a -installsuffix cgo -buildvcs=false -tags 'netgo,osusergo,release' -gcflags=all='-l -B' -mod=readonly

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
CYAN := \033[36m
RESET := \033[0m

# Default target
all: clean deps build

# Production build (maximum size optimization)
build:
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -rf $(OUTPUT_DIR)
	@echo "$(GREEN)Build artifacts cleaned$(RESET)"
	@echo "$(BLUE)Building $(PROJECT_NAME) v$(VERSION) for production (Linux/AMD64)...$(RESET)"
	@echo "$(CYAN)Commit: $(GIT_COMMIT), Branch: $(GIT_BRANCH)$(RESET)"
	@echo "$(CYAN)Target: $(GOOS)/$(GOARCH) (optimized static binary)$(RESET)"
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o $(OUTPUT_DIR)/$(PROJECT_NAME) $(MAIN_FILE) || { \
		echo "$(RED)Build failed$(RESET)"; \
		exit 1; \
	}
	@echo "$(GREEN)Production build complete: $(OUTPUT_DIR)/$(PROJECT_NAME)$(RESET)"
	@ls -lh $(OUTPUT_DIR)/$(PROJECT_NAME)
	@file $(OUTPUT_DIR)/$(PROJECT_NAME)

# Run the application
run: build
	@echo "$(BLUE)Running $(PROJECT_NAME)...$(RESET)"
	@$(OUTPUT_DIR)/$(PROJECT_NAME) $(ARGS)

# Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning all artifacts...$(RESET)"
	@rm -rf $(OUTPUT_DIR)
	@go clean -cache -modcache -testcache
	@echo "$(GREEN)All artifacts cleaned$(RESET)"

# Clean build artifacts only
clean-bin:
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -rf $(OUTPUT_DIR)
	@echo "$(GREEN)Build artifacts cleaned$(RESET)"

# Install binary to GOPATH/bin
install: build
	@echo "$(BLUE)Installing $(PROJECT_NAME) to GOPATH/bin...$(RESET)"
	@go install $(BUILD_FLAGS) $(MAIN_FILE) || { \
		echo "$(RED)Installation failed$(RESET)"; \
		exit 1; \
	}
	@echo "$(GREEN)Installation complete$(RESET)"

# Format Go code
fmt:
	@echo "$(BLUE)Formatting Go code...$(RESET)"
	@gofmt -s -w . || { \
		echo "$(RED)Code formatting failed$(RESET)"; \
		exit 1; \
	}
	@echo "$(GREEN)Code formatting complete$(RESET)"

# Lint code
lint:
	@echo "$(BLUE)Running golangci-lint...$(RESET)"
	@golangci-lint run ./... || { \
		echo "$(RED)Linting failed$(RESET)"; \
		exit 1; \
	}
	@echo "$(GREEN)Linting complete$(RESET)"

# Run tests
test:
	@echo "$(BLUE)Running tests...$(RESET)"
	@go test -short ./... || { \
		echo "$(RED)Tests failed$(RESET)"; \
		exit 1; \
	}
	@echo "$(GREEN)Tests complete$(RESET)"

# Install and verify dependencies
deps:
	@echo "$(BLUE)Installing dependencies...$(RESET)"
	@go mod download || { \
		echo "$(RED)Dependency installation failed$(RESET)"; \
		exit 1; \
	}
	@go mod verify || { \
		echo "$(RED)Dependency verification failed$(RESET)"; \
		exit 1; \
	}
	@echo "$(GREEN)Dependencies installed and verified$(RESET)"

# Tidy dependencies
tidy:
	@echo "$(BLUE)Tidying dependencies...$(RESET)"
	@go mod tidy || { \
		echo "$(RED)Dependency tidying failed$(RESET)"; \
		exit 1; \
	}
	@echo "$(GREEN)Dependencies tidied$(RESET)"

# Help information
help:
	@echo "$(BLUE)DevBox Pack - Production Makefile (Linux/AMD64 only)$(RESET)"
	@echo ""
	@echo "$(CYAN)Project:$(RESET) $(PROJECT_NAME) $(VERSION)"
	@echo "$(CYAN)Target:$(RESET) $(GOOS)/$(GOARCH)"
	@echo ""
	@echo "$(YELLOW)Production Targets:$(RESET)"
	@echo "  build        - Build optimized production binary (maximum size optimization)"
	@echo "  run          - Build and run the application (use ARGS=\"arg1 arg2\" for arguments)"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(RESET)"
	@echo "  fmt          - Format Go code with gofmt"
	@echo "  lint         - Run golangci-lint"
	@echo "  test         - Run unit tests"
	@echo ""
	@echo "$(YELLOW)Dependencies:$(RESET)"
	@echo "  deps         - Install and verify dependencies"
	@echo "  tidy         - Tidy dependencies"
	@echo ""
	@echo "$(YELLOW)Maintenance:$(RESET)"
	@echo "  clean        - Clean all artifacts and caches"
	@echo "  clean-bin    - Clean build artifacts only"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "$(YELLOW)Shortcuts:$(RESET)"
	@echo "  all          - clean + deps + build (default target)"
	@echo ""
	@echo "$(GREEN)Example: make clean && make all$(RESET)"
	@echo "$(CYAN)Target Platform: Linux/AMD64 (static binaries)$(RESET)"