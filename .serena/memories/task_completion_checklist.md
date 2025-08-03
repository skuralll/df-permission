# Task Completion Checklist

## Required Steps When Completing Any Development Task

### 1. Code Quality Checks
```bash
# Format code
go fmt ./...

# Vet for potential issues
go vet ./...

# Lint (if available)
golangci-lint run
```

### 2. Testing Requirements
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Test specific affected packages
go test ./permission          # If public API changed
go test ./internal/domain     # If core logic changed
go test ./internal/repository # If storage changed
```

### 3. Build Verification
```bash
# Ensure project builds
go build ./...

# Update dependencies if needed
go mod tidy
```

### 4. Documentation Updates (if applicable)
- Update CLAUDE.md for development guide changes
- Update README.md for user-facing changes
- Update code comments for new/modified public APIs

### 5. Git Operations
```bash
# Check status
git status

# Add changes
git add .

# Commit with descriptive message
git commit -m "descriptive message"

# Push if ready (usually not automatically)
git push origin branch-name
```

## Specific Testing Scenarios

### For Permission Logic Changes
```bash
go test ./internal/domain -v
go test ./permission -run TestPermissionManager_HasPermission
```

### For Storage Changes
```bash
go test ./internal/repository -v
go test ./internal/application -v  # Test integration
```

### For Public API Changes
```bash
go test ./permission -v
# Ensure backward compatibility
```

### For Configuration/Options Changes
```bash
go test ./permission -run TestWithStorage
go test ./permission -run TestMultipleOptions
```

## Pre-Commit Checklist
- [ ] Code formatted (`go fmt ./...`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] All tests pass (`go test ./...`)
- [ ] Project builds (`go build ./...`)
- [ ] Dependencies clean (`go mod tidy`)
- [ ] Documentation updated (if needed)
- [ ] Error handling appropriate
- [ ] Thread safety maintained
- [ ] No performance regressions