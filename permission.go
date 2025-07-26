package dfpermission

import (
	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/application"
)

// PermissionManager provides a public interface for permission management operations.
// It wraps the internal manager implementation and exposes only essential methods.
type PermissionManager interface {
	// Permission checking
	HasPermission(playerID uuid.UUID, permission string) bool

	// Group management
	CreateGroup(name string, permissions []string) error
	DeleteGroup(name string) error

	// Player-group relationships
	AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error
	RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error

	// Individual player permissions
	AddPlayerPermission(playerID uuid.UUID, permission string) error
	RemovePlayerPermission(playerID uuid.UUID, permission string) error

	// System operations
	Save() error
}

// permissionManager is the concrete implementation that wraps internal.Manager
type permissionManager struct {
	internal *application.Manager
}

// NewManager creates a new PermissionManager with the given configuration.
// It wraps the internal manager and provides a stable public API.
func NewManager(config ManagerConfig) PermissionManager {
	internalMgr := application.NewManager(config)
	return &permissionManager{
		internal: internalMgr,
	}
}

// HasPermission checks if a player has a specific permission.
// Returns true if the player has the permission either directly or through group membership.
func (p *permissionManager) HasPermission(playerID uuid.UUID, permission string) bool {
	return p.internal.HasPermission(playerID, permission)
}

// CreateGroup creates a new permission group with the specified permissions.
// Returns an error if the group already exists.
func (p *permissionManager) CreateGroup(name string, permissions []string) error {
	err := p.internal.CreateGroup(name, permissions)
	return convertError(err)
}

// DeleteGroup removes a permission group.
// Returns an error if the group doesn't exist or is a system group.
func (p *permissionManager) DeleteGroup(name string) error {
	err := p.internal.DeleteGroup(name)
	return convertError(err)
}

// AddPlayerToGroup adds a player to a permission group.
// Creates the player if they don't exist.
func (p *permissionManager) AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error {
	err := p.internal.AddPlayerToGroup(playerID, playerName, groupName)
	return convertError(err)
}

// RemovePlayerFromGroup removes a player from a permission group.
// Returns an error if the player is not in the group.
func (p *permissionManager) RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error {
	err := p.internal.RemovePlayerFromGroup(playerID, groupName)
	return convertError(err)
}

// AddPlayerPermission grants a specific permission directly to a player.
// Returns an error if the player doesn't exist.
func (p *permissionManager) AddPlayerPermission(playerID uuid.UUID, permission string) error {
	err := p.internal.AddPlayerPermission(playerID, permission)
	return convertError(err)
}

// RemovePlayerPermission removes a specific permission from a player.
// Returns an error if the player doesn't have the permission.
func (p *permissionManager) RemovePlayerPermission(playerID uuid.UUID, permission string) error {
	err := p.internal.RemovePlayerPermission(playerID, permission)
	return convertError(err)
}

// Save persists the current permission data to storage.
func (p *permissionManager) Save() error {
	err := p.internal.Save()
	return convertError(err)
}