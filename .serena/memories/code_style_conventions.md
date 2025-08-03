# Code Style and Conventions

## Language Standards
- **Go Version**: 1.24+ (specified in go.mod)
- **Module Path**: github.com/skuralll/df-permission

## Naming Conventions
- **Packages**: lowercase, single word or abbreviated (e.g., `permission`, `domain`)
- **Public Interface**: PascalCase (e.g., `PermissionManager`)
- **Private Structs**: camelCase (e.g., `permissionManager`)
- **Methods**: PascalCase for exported, camelCase for unexported
- **Constants/Variables**: Standard Go conventions

## Documentation Style
- **Interface Comments**: Japanese comments describing purpose (e.g., "権限管理操作のための公開インターフェース")
- **Method Comments**: Brief Japanese explanations of functionality
- **Package Documentation**: Mix of Japanese and English
- **README**: Japanese language for user-facing documentation

## Code Organization Patterns
- **Options Pattern**: Used for configuration (`WithStorage`, `WithCache`, etc.)
- **Interface Segregation**: Clean public API wrapping internal implementation
- **Error Wrapping**: Structured error system with contextual information
- **Factory Pattern**: `NewManager()` constructor with options
- **Repository Pattern**: Storage abstraction with interface

## Testing Conventions
- **Test Files**: `*_test.go` alongside source files
- **Test Functions**: `Test[PackageName]_[Functionality]` pattern
- **Table-Driven Tests**: Used extensively for comprehensive coverage
- **Test Helpers**: Helper functions like `createTestPermissionManager`
- **Cleanup**: Proper teardown with temporary files and resource cleanup

## Import Organization
```go
import (
    // Standard library
    "time"
    "context"
    
    // Third-party packages
    "github.com/google/uuid"
    
    // Local packages
    "github.com/skuralll/df-permission/internal/application"
)
```

## Error Handling
- **Error Variables**: Predefined error variables (e.g., `ErrPlayerNotFound`)
- **Error Wrapping**: Context-aware error construction
- **Error Interface**: Clean public error interface hiding internal details
- **Error Conversion**: Translation between internal and public errors