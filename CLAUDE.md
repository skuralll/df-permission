# df-permission - Claude Development Guide

## Project Overview

**df-permission** is a comprehensive permission management library for Minecraft Bedrock Edition Dragonfly servers, written in Go. It provides a robust, thread-safe, cached permission system with wildcard support, group-based permissions, and flexible storage options.

### Key Features
- **Wildcard Permission System**: Supports `*` (global) and `prefix.*` (prefix) wildcards
- **Group-Based Management**: Hierarchical permission management through groups
- **Thread-Safe Operations**: Full concurrent operation support with RWMutex
- **Performance Optimized**: TTL-based caching with automatic cleanup
- **Flexible Configuration**: Options pattern for customizable setup
- **Clean Architecture**: Layered design following domain-driven principles

## Architecture Overview

The project follows a clean, layered architecture pattern:

```
┌─────────────────────┐
│   Public API Layer  │  permission/ (Public interface)
├─────────────────────┤
│ Application Layer   │  internal/application/ (Business orchestration)
├─────────────────────┤
│   Domain Layer      │  internal/domain/ (Core business logic)
├─────────────────────┤
│ Repository Layer    │  internal/repository/ (Data access)
├─────────────────────┤
│   Shared Layer      │  internal/shared/ (Common types & errors)
└─────────────────────┘
```

### Layer Responsibilities

#### 1. Public API Layer (`permission/`)
- **Files**: `permission.go`, `options.go`, `errors.go`
- **Purpose**: Stable public interface, error conversion, configuration options
- **Key Interface**: `PermissionManager` (13 core methods)
- **Pattern**: Facade pattern wrapping internal implementation

#### 2. Application Layer (`internal/application/`)
- **Files**: `manager.go`, `manager_test.go`
- **Purpose**: Business logic orchestration, transaction management, caching coordination
- **Key Component**: `Manager` struct with full CRUD operations
- **Features**: Auto-save, cache invalidation, concurrent access control

#### 3. Domain Layer (`internal/domain/`)
- **Files**: `checker.go`, `matcher.go`, `cache.go`, `storage.go`
- **Purpose**: Core business logic, permission checking, pattern matching, caching
- **Key Components**:
  - `PermissionChecker`: Thread-safe permission validation
  - `PermissionMatcher`: Wildcard pattern matching logic
  - `PermissionCache`: TTL-based result caching

#### 4. Repository Layer (`internal/repository/`)
- **Files**: `storage.go`, `json_storage.go`, `json_storage_test.go`
- **Purpose**: Data persistence abstraction
- **Key Interface**: `Storage` with Load/Save/Exists/Close methods
- **Implementation**: JSON file-based storage with validation

#### 5. Shared Layer (`internal/shared/`)
- **Files**: `types.go`, `errors.go`
- **Purpose**: Common data structures, error definitions, configuration types
- **Key Types**: `PermissionData`, `PlayerData`, `Group`, `ManagerConfig`

## Development Commands

### Essential Commands

```bash
# Build the project
go build ./...

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./permission
go test ./internal/domain
go test ./internal/application

# Clean module dependencies
go mod tidy

# Format code
go fmt ./...

# Lint (if golangci-lint is available)
golangci-lint run
```

### Testing Strategy

The codebase uses comprehensive table-driven tests:

```bash
# Test structure follows package hierarchy
permission/permission_test.go          # Public API tests
internal/application/manager_test.go    # Business logic tests
internal/domain/checker_test.go         # Permission checking tests
internal/domain/matcher_test.go         # Pattern matching tests
internal/domain/cache_test.go           # Cache functionality tests
internal/repository/json_storage_test.go # Storage tests
```

**Test Patterns Used**:
- Setup/teardown with temporary files
- UUID-based test data
- Error condition testing
- Concurrent access testing
- Cache behavior validation

## Key Data Structures

### Core Types

```go
// Main permission data container
type PermissionData struct {
    Groups  map[string]*Group         `json:"groups"`
    Players map[uuid.UUID]*PlayerData `json:"players"`  
    Meta    *Metadata                 `json:"meta"`
}

// Player information with permissions and group memberships
type PlayerData struct {
    PlayerID    uuid.UUID `json:"player_id"`
    PlayerName  string    `json:"player_name"`
    Groups      []string  `json:"groups"`
    Permissions []string  `json:"permissions"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Permission group definition
type Group struct {
    Name        string   `json:"name"`
    Permissions []string `json:"permissions"`
}
```

### Configuration Types

```go
type ManagerConfig struct {
    AutoSave bool
    Storage  StorageConfig
    Cache    CacheConfig
}

type CacheConfig struct {
    TTL             time.Duration
    CleanupInterval time.Duration
    Enabled         bool
}
```

## Permission System Logic

### Wildcard Patterns

1. **Global Wildcard (`*`)**: Grants all permissions
2. **Prefix Wildcard (`prefix.*`)**: Grants all permissions starting with `prefix.`
3. **Exact Match**: Grants specific permission only

### Permission Resolution Order

1. Check player's direct permissions
2. Check permissions from all player's groups
3. Apply wildcard matching logic
4. Cache result for performance

### Default Groups

- **`default`**: `["chat.send", "world.interact"]` - Basic player permissions
- **`admin`**: `["*"]` - Administrative access (all permissions)

## Error Handling

The library uses a structured error system with error wrapping:

### Public Errors (re-exported from shared)
```go
var (
    ErrPlayerNotFound           // Player doesn't exist
    ErrPlayerAlreadyExists      // Player already exists  
    ErrGroupNotFound            // Group doesn't exist
    ErrGroupAlreadyExists       // Group already exists
    ErrSystemGroupProtected     // Cannot delete system groups
    ErrPlayerPermissionNotFound // Player lacks specific permission
    ErrGroupPermissionNotFound  // Group lacks specific permission
    ErrInvalidPermission        // Invalid permission format
    ErrStorage                  // Storage operation failed
)
```

### Error Construction Pattern
Errors include contextual information using format strings:
```go
NewPlayerNotFoundError(playerID) // "player not found: player with ID {uuid}"
NewGroupNotFoundError(groupName) // "group not found: '{name}'"
```

## Common Development Patterns

### Creating Managers

```go
// Default configuration
mgr := permission.NewManager()

// Custom configuration with options
mgr := permission.NewManager(
    permission.WithStorage("./custom.json"),
    permission.WithCache(45*time.Second),
    permission.WithAutoSave(true),
    permission.WithCacheCleanup(2*time.Minute),
)
```

### Error Handling Pattern

```go
err := mgr.CreateGroup("vip", []string{"chat.color"})
if err != nil {
    if errors.Is(err, permission.ErrGroupAlreadyExists) {
        // Handle group exists case
    }
    // Handle other errors
}
```

### Thread Safety

All operations are thread-safe via RWMutex:
- Read operations (HasPermission, GetPlayerData) use RLock
- Write operations (CreateGroup, AddPlayer) use Lock
- Cache operations are independently thread-safe

## Integration Notes

### Dragonfly Server Integration
- Located in `internal/dragonfly/command/` (currently empty)
- Designed for integration with Dragonfly command system
- UUID-based player identification matches Minecraft standards

### Storage Format
- JSON file format for human readability
- Includes metadata with version tracking
- Atomic write operations for data consistency
- Automatic backup/recovery patterns supported

## Performance Considerations

### Caching Strategy
- TTL-based permission result caching
- Automatic cache invalidation on permission changes
- Background cleanup goroutine for expired entries
- Configurable cache settings per use case

### Memory Management
- Efficient UUID-based indexing
- Copy-on-read pattern for thread safety
- Automatic cleanup of expired cache entries
- Minimal memory allocation in hot paths

### Concurrency Design
- Read-heavy workload optimization
- Non-blocking cache access
- Goroutine-safe auto-save operations
- Lock-free permission checking where possible

## Extension Points

### Adding New Storage Backends
Implement the `Storage` interface:
```go
type Storage interface {
    Load() (*shared.PermissionData, error)
    Save(data *shared.PermissionData) error
    Exists() bool
    Close() error
}
```

### Custom Permission Matchers
Extend `PermissionMatcher` with new pattern types while maintaining the existing interface.

### Additional Configuration Options
Add new `Option` functions following the existing pattern:
```go
func WithNewFeature(value SomeType) Option {
    return func(config *shared.ManagerConfig) {
        config.NewFeature = value
    }
}
```

## Dependencies

- **Go 1.24+** (specified in go.mod)
- **github.com/google/uuid v1.6.0** - UUID generation and handling
- **Standard library only** - No external dependencies for core functionality

This architecture provides a solid foundation for permission management with clear separation of concerns, comprehensive testing, and excellent performance characteristics.