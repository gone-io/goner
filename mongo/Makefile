.PHONY: test build clean install lint fmt vet

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=gofmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Build the project
build:
	$(GOBUILD) -v ./...

# Test the project
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f coverage.out coverage.html

# Install dependencies
install:
	$(GOGET) -v ./...

# Format code
fmt:
	$(GOFMT) -s -w .

# Vet code
vet:
	$(GOVET) ./...

# Lint code (requires golangci-lint)
lint:
	$(GOLINT) run

# Run all checks
check: fmt vet lint test

# Update dependencies
update:
	$(GOCMD) get -u ./...
	$(GOCMD) mod tidy

# Tidy go modules
tidy:
	$(GOCMD) mod tidy

# Download dependencies
download:
	$(GOCMD) mod download

# Verify dependencies
verify:
	$(GOCMD) mod verify

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the project"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  install      - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Lint code"
	@echo "  check        - Run all checks (fmt, vet, lint, test)"
	@echo "  update       - Update dependencies"
	@echo "  tidy         - Tidy go modules"
	@echo "  download     - Download dependencies"
	@echo "  verify       - Verify dependencies"
	@echo "  help         - Show this help"