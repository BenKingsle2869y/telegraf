# Makefile for telegraf
# Provides common build, test, and lint targets

PLATFORM ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "unknown")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

BINARY := telegraf
MAIN_PKG := ./cmd/telegraf

LDFLAGS := -ldflags "-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.branch=$(BRANCH) \
	-X main.buildDate=$(BUILD_DATE)"

.PHONY: all build clean test lint fmt vet deps docker help

all: deps build

## build: Compile the telegraf binary
build:
	@echo "Building $(BINARY) $(VERSION) for $(PLATFORM)/$(ARCH)..."
	go build $(LDFLAGS) -o $(BINARY) $(MAIN_PKG)

## build-static: Compile a statically linked binary
build-static:
	@echo "Building static $(BINARY) $(VERSION)..."
	CGO_ENABLED=0 go build $(LDFLAGS) -a -installsuffix cgo -o $(BINARY) $(MAIN_PKG)

## test: Run all unit tests
test:
	@echo "Running tests..."
	go test -v -race -timeout 30s ./...

## test-short: Run short tests only
test-short:
	@echo "Running short tests..."
	go test -short -timeout 30s ./...

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

## fmt: Format Go source files
fmt:
	@echo "Formatting source files..."
	gofmt -s -w .

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## deps: Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY)
	rm -rf dist/

## docker: Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t telegraf:$(VERSION) .

## install: Install telegraf to GOPATH/bin
install:
	@echo "Installing $(BINARY)..."
	go install $(LDFLAGS) $(MAIN_PKG)

## check: Run fmt, vet, and lint
check: fmt vet lint

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /' | column -t -s ':'
