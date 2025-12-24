.PHONY: test build fmt lint version release-patch release-minor release-major help

# Default target
help:
	@echo "Ujeebu Go SDK - Available targets:"
	@echo ""
	@echo "  Development:"
	@echo "    test          Run all tests"
	@echo "    test-cover    Run tests with coverage"
	@echo "    build         Build the package"
	@echo "    fmt           Format code"
	@echo "    lint          Run linter (requires golangci-lint)"
	@echo ""
	@echo "  Versioning:"
	@echo "    version       Show current version"
	@echo "    release-patch Bump patch version (x.x.X)"
	@echo "    release-minor Bump minor version (x.X.0)"
	@echo "    release-major Bump major version (X.0.0)"
	@echo ""
	@echo "  Examples:"
	@echo "    examples      Build all examples"
	@echo ""

# Development
test:
	go test ./...

test-cover:
	go test -cover ./...

test-verbose:
	go test -v ./...

build:
	go build ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

# Examples
examples:
	go build ./examples/...

# Versioning
version:
	@git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"

release-patch:
	@./scripts/bump-version.sh patch

release-minor:
	@./scripts/bump-version.sh minor

release-major:
	@./scripts/bump-version.sh major

# Clean
clean:
	go clean ./...
	rm -f coverage.out

