.PHONY: all build install clean test fmt help

# Application name
APP_NAME := llm-radar

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOINSTALL := $(GOCMD) install
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOFMT := gofmt
GOMOD := $(GOCMD) mod

# Build directory
BUILD_DIR := bin

# Binary name
BINARY_NAME := $(APP_NAME)
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)

# Installation path
INSTALL_PATH := $(HOME)/.local/bin

# Version info
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
all: build

## build: Build the application binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) .
	@echo "✅ Build complete: $(BINARY_PATH)"

## install: Install the application to $(INSTALL_PATH)
install: build
	@echo "Installing $(APP_NAME) to $(INSTALL_PATH)..."
	@mkdir -p $(INSTALL_PATH)
	@cp $(BINARY_PATH) $(INSTALL_PATH)/$(BINARY_NAME)
	@chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Installed to $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "Make sure $(INSTALL_PATH) is in your PATH"

## install-global: Install globally to /usr/local/bin (requires sudo)
install-global: build
	@echo "Installing $(APP_NAME) globally..."
	@sudo cp $(BINARY_PATH) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Globally installed to /usr/local/bin/$(BINARY_NAME)"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

## test: Run tests (when available)
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## fmt: Format Go source code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -w .
	@echo "✅ Code formatted"

## mod-download: Download dependencies
mod-download:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	@echo "✅ Dependencies downloaded"

## mod-tidy: Tidy go.mod and go.sum
mod-tidy:
	@echo "Tidying modules..."
	$(GOMOD) tidy
	@echo "✅ Modules tidied"

## run: Build and run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(BINARY_PATH)

## run-cached: Run with cache enabled
run-cached: build
	@echo "Running $(APP_NAME) with cache..."
	@$(BINARY_PATH) --cache

## run-fast: Run with limited models (fast test)
run-fast: build
	@echo "Running $(APP_NAME) with 2 workers..."
	@$(BINARY_PATH) -c 2

## version: Show version information
version: build
	@$(BINARY_PATH) --version

## help: Show this help message
help:
	@echo "$(APP_NAME) Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^##//p' Makefile | column -t -s ':' | sed -e 's/^/ /'
