# Suggested Commands for df-permission Development

## Essential Go Commands

### Build and Dependencies
```bash
# Build the entire project
go build ./...

# Clean and update dependencies
go mod tidy

# Format all Go code
go fmt ./...

# Generate code (if needed)
go generate ./...
```

### Testing Commands
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific packages
go test ./permission
go test ./internal/domain
go test ./internal/application
go test ./internal/repository

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestPermissionManager_HasPermission ./permission
```

### Development Tools
```bash
# Lint (if golangci-lint is available)
golangci-lint run

# Check for potential issues
go vet ./...

# List all dependencies
go list -m all

# Show module info
go mod graph
```

### Release Commands
```bash
# Tag for release (triggers GitHub Actions)
git tag v1.0.0
git push origin v1.0.0

# Test GoReleaser locally
goreleaser release --snapshot --clean
```

## Linux System Commands
```bash
# File operations
ls -la                    # List files with details
find . -name "*.go"       # Find Go files
grep -r "pattern" .       # Search in files
rg "pattern"              # ripgrep (faster grep)

# Git operations
git status
git log --oneline -10
git diff
git add .
git commit -m "message"
git push

# Directory navigation
pwd                       # Current directory
cd /path/to/project      # Change directory
```

## Project-Specific Testing
```bash
# Test public API
go test ./permission -v

# Test core domain logic
go test ./internal/domain -v

# Test storage layer
go test ./internal/repository -v

# Test application layer
go test ./internal/application -v

# Run all tests with race detection
go test -race ./...
```