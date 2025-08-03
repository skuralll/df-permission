# Design Patterns and Guidelines

## Architectural Patterns

### Clean Architecture
The project follows a strict layered architecture:
1. **Public API Layer** (`permission/`): Stable public interface
2. **Application Layer** (`internal/application/`): Business orchestration  
3. **Domain Layer** (`internal/domain/`): Core business logic
4. **Repository Layer** (`internal/repository/`): Data access
5. **Shared Layer** (`internal/shared/`): Common types and errors

### Key Design Patterns Used

#### 1. Options Pattern
Used for flexible configuration:
```go
mgr := permission.NewManager(
    permission.WithStorage("./custom.json"),
    permission.WithCache(45*time.Second),
    permission.WithAutoSave(true),
)
```

#### 2. Facade Pattern
Public API wraps internal implementation:
```go
type permissionManager struct {
    internal *application.Manager
}
```

#### 3. Repository Pattern
Storage abstraction with clean interface:
```go
type Storage interface {
    Load() (*shared.PermissionData, error)
    Save(data *shared.PermissionData) error
    Exists() bool
    Close() error
}
```

#### 4. Factory Pattern
Manager creation with dependency injection:
```go
func NewManager(opts ...Option) PermissionManager
```

## Thread Safety Guidelines

### Concurrency Design
- **Read-heavy workload optimization**: RWMutex usage
- **Thread-safe operations**: All public methods are safe for concurrent use
- **Cache thread safety**: Independent concurrent access support
- **Non-blocking permission checking**: Lock-free hot paths where possible

### Locking Strategy
- Read operations (`HasPermission`, `GetPlayerData`): Use `RLock`
- Write operations (`CreateGroup`, `AddPlayer`): Use `Lock`
- Cache operations: Independently thread-safe

## Error Handling Patterns

### Structured Errors
```go
var (
    ErrPlayerNotFound           // Player doesn't exist
    ErrPlayerAlreadyExists      // Player already exists  
    ErrGroupNotFound            // Group doesn't exist
    // ... more specific errors
)
```

### Error Construction
```go
NewPlayerNotFoundError(playerID) // "player not found: player with ID {uuid}"
NewGroupNotFoundError(groupName) // "group not found: '{name}'"
```

### Error Wrapping
- Use contextual information in error messages
- Convert internal errors to public errors at API boundary
- Maintain error chain for debugging

## Performance Guidelines

### Caching Strategy
- TTL-based permission result caching
- Automatic cache invalidation on permission changes
- Background cleanup for expired entries
- Configurable cache settings

### Memory Management
- Efficient UUID-based indexing
- Copy-on-read pattern for thread safety
- Minimal memory allocation in hot paths
- Automatic cleanup of expired cache entries

## Data Structure Guidelines

### Core Types
- Use `uuid.UUID` for player identification (Minecraft standard)
- Use `time.Time` for timestamps and TTL
- Use `map[string]` for string-based lookups
- Use slices for permission lists

### Immutability
- Return copies of internal data to prevent external modification
- Use value receivers where appropriate
- Protect internal state from external changes

## Testing Guidelines

### Test Structure
- Table-driven tests for comprehensive coverage
- Separate test files alongside source (`*_test.go`)
- Setup/teardown with temporary files and proper cleanup
- UUID-based test data for realistic scenarios

### Test Patterns
- Error condition testing for all failure modes
- Concurrent access testing for thread safety
- Cache behavior validation
- Integration testing between layers