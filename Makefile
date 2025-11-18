# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install

# Binary name
BINARY_NAME=app-hexagonal
BINARY_UNIX=$(BINARY_NAME)_unix

# Application
APP_PATH=./cmd/main.go

# Default target
all: build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(APP_PATH)

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(APP_PATH)

# Run the application
run:
	$(GOCMD) run $(APP_PATH)

# Run with Air for development (hot reload)
dev:
	air -c .air.linux.conf

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Install dependencies
deps:
	$(GOMOD) download

# Tidy go.mod and go.sum
tidy:
	$(GOMOD) tidy

# Update dependencies
update:
	$(GOGET) -u ./...

# Generate Swagger documentation
swag:
	swag init -g $(APP_PATH)

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code for potential issues
vet:
	$(GOCMD) vet ./...

# Install tools needed for development
install-tools:
	$(GOINSTALL) github.com/swaggo/swag/cmd/swag@latest
	$(GOINSTALL) github.com/cosmtrek/air@latest
	$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
lint:
	golangci-lint run

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Build the application (default)"
	@echo "  build        - Build the application"
	@echo "  build-linux  - Build for Linux"
	@echo "  run          - Run the application"
	@echo "  dev          - Run with Air for development (hot reload)"
	@echo "  test         - Run tests"
	@echo "  test-cover   - Run tests with coverage"
	@echo "  clean        - Clean build files"
	@echo "  deps         - Install dependencies"
	@echo "  tidy         - Tidy go.mod and go.sum"
	@echo "  update       - Update dependencies"
	@echo "  swag         - Generate Swagger documentation"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code for potential issues"
	@echo "  install-tools - Install tools needed for development"
	@echo "  lint         - Run linter"
	@echo "  help         - Show this help message"