.PHONY: test test-unit test-integration test-coverage test-race test-verbose clean mocks

# Generate mocks
mocks:
	@echo "🔨 Generating mocks..."
	@mockgen -source=internal/service/event_service.go -destination=internal/mocks/mock_service.go -package=mocks
	@mockgen -source=internal/repository/event_repository.go -destination=internal/mocks/mock_repository.go -package=mocks
	@echo "✅ Mocks generated successfully"

# Run all tests
test:
	@echo "🧪 Running all tests..."
	@go test ./... -v

# Run only unit tests
test-unit:
	@echo "🧪 Running unit tests..."
	@go test ./internal/domain ./internal/service ./internal/middleware ./geolocation -v

# Run only integration tests
test-integration:
	@echo "🧪 Running integration tests..."
	@go test ./internal/tests -v

# Run integration tests in short mode (skip slow tests)
test-short:
	@echo "🧪 Running quick tests..."
	@go test ./... -short -v

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total

# Run tests with race detection
test-race:
	@echo "🏁 Running tests with race detector..."
	@go test ./... -race -v

# Run tests with verbose output and race detection
test-verbose:
	@echo "🔍 Running verbose tests with race detection..."
	@go test ./... -v -race -count=1

# Run benchmarks
bench:
	@echo "⚡ Running benchmarks..."
	@go test ./... -bench=. -benchmem -run=^Benchmark

# Clean test artifacts
clean:
	@echo "🧹 Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@go clean -testcache

# Install test dependencies
deps:
	@echo "📦 Installing test dependencies..."
	@go get -t ./...
	@go mod tidy

# Run specific package tests
test-domain:
	@go test ./internal/domain -v

test-service:
	@go test ./internal/service -v

test-middleware:
	@go test ./internal/middleware -v

test-geolocation:
	@go test ./geolocation -v

# CI/CD test command
test-ci:
	@echo "🤖 Running CI tests..."
	@go test ./... -race -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out

# Help
help:
	@echo "Available test commands:"
	@echo "  make test               - Run all tests"
	@echo "  make test-unit          - Run unit tests only"
	@echo "  make test-integration   - Run integration tests only"
	@echo "  make test-short         - Run quick tests (skip slow)"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make test-race          - Run tests with race detection"
	@echo "  make test-verbose       - Run verbose tests"
	@echo "  make bench              - Run benchmarks"
	@echo "  make mocks              - Generate mock files"
	@echo "  make clean              - Clean test artifacts"
	@echo "  make deps               - Install test dependencies"
	@echo "  make test-ci            - Run CI tests"
