# Weather TUI Makefile

# Variables
BINARY_NAME=weather-tui
MAIN_PACKAGE=.
BUILD_DIR=./bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test deps run dev install help

# Default target
all: deps build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "Multi-platform build complete"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated"

# Run the application directly
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PACKAGE)

# Development mode with live reload (requires entr)
dev:
	@echo "Starting development mode..."
	@echo "Note: Install 'entr' for auto-reload: brew install entr (macOS) or apt install entr (Ubuntu)"
	@find . -name "*.go" | entr -r $(GOCMD) run $(MAIN_PACKAGE)

# Install to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(MAIN_PACKAGE)
	@echo "Installation complete"

# Create release build
release: clean deps build-all
	@echo "Creating release archive..."
	@cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-release.tar.gz *
	@echo "Release created: $(BUILD_DIR)/$(BINARY_NAME)-release.tar.gz"

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from template"; fi
	@echo "Please edit .env file and add your OpenWeatherMap API key"
	@echo "Get your API key at: https://openweathermap.org/api"
	$(MAKE) deps

# Check for required tools
check-tools:
	@echo "Checking required tools..."
	@command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Visit: https://golang.org/"; exit 1; }
	@echo "✓ Go is installed: $$(go version)"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "Code formatting complete"

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found. Install: https://golangci-lint.run/"; exit 1; }
	golangci-lint run
	@echo "Linting complete"

# Update dependencies
update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "Dependencies updated"

# Show help
help:
	@echo "Weather TUI Makefile Commands:"
	@echo ""
	@echo "  make build      - Build the application"
	@echo "  make build-all  - Build for multiple platforms"
	@echo "  make run        - Run the application directly"
	@echo "  make dev        - Run in development mode with auto-reload"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Download and tidy dependencies"
	@echo "  make install    - Install to GOPATH/bin"
	@echo "  make release    - Create release build for all platforms"
	@echo "  make setup      - Setup development environment"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Lint code (requires golangci-lint)"
	@echo "  make update     - Update dependencies"
	@echo "  make check-tools- Check for required tools"
	@echo "  make help       - Show this help message"
	@echo ""
	@echo "Environment Variables:"
	@echo "  OPENWEATHER_API_KEY - Your OpenWeatherMap API key"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make setup       # Setup development environment"
	@echo "  2. Edit .env file   # Add your API key"
	@echo "  3. make run         # Run the application"
