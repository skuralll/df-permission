package dfpermission

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func createTestPermissionManager() PermissionManager {
	config := ManagerConfig{
		AutoSave: false,
		Storage: StorageConfig{
			Path: "/tmp/test_permission_public.json",
		},
		Cache: CacheConfig{
			Enabled:         false,
			TTL:             5 * time.Second,
			CleanupInterval: 10 * time.Second,
		},
	}
	return NewManager(config)
}

func cleanupPublicTest() {
	os.Remove("/tmp/test_permission_public.json")
}

func TestPermissionManager_Interface(t *testing.T) {
	defer cleanupPublicTest()
	mgr := createTestPermissionManager()

	if mgr == nil {
		t.Fatal("PermissionManager creation failed")
	}
}

func TestPermissionManager_HasPermission(t *testing.T) {
	defer cleanupPublicTest()
	mgr := createTestPermissionManager()
	playerID := uuid.New()

	// Non-existent player should not have permissions
	if mgr.HasPermission(playerID, "any.permission") {
		t.Error("Non-existent player should not have any permissions")
	}

	// Add player to admin group and test
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player to admin group: %v", err)
	}

	if !mgr.HasPermission(playerID, "any.permission") {
		t.Error("Admin should have all permissions")
	}
}

func TestPermissionManager_GroupManagement(t *testing.T) {
	defer cleanupPublicTest()
	mgr := createTestPermissionManager()

	// Create a new group
	err := mgr.CreateGroup("vip", []string{"chat.color", "world.build.fast"})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	// Try to create duplicate group
	err = mgr.CreateGroup("vip", []string{})
	if err != ErrGroupAlreadyExists {
		t.Errorf("Expected ErrGroupAlreadyExists, got %v", err)
	}

	// Delete the group
	err = mgr.DeleteGroup("vip")
	if err != nil {
		t.Fatalf("Failed to delete group: %v", err)
	}

	// Try to delete system group
	err = mgr.DeleteGroup("admin")
	if err != ErrSystemGroupProtected {
		t.Errorf("Expected ErrSystemGroupProtected, got %v", err)
	}
}

func TestPermissionManager_PlayerGroupMembership(t *testing.T) {
	defer cleanupPublicTest()
	mgr := createTestPermissionManager()
	playerID := uuid.New()

	// Create a test group
	err := mgr.CreateGroup("vip", []string{"chat.color"})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	// Add player to group
	err = mgr.AddPlayerToGroup(playerID, "TestPlayer", "vip")
	if err != nil {
		t.Fatalf("Failed to add player to group: %v", err)
	}

	// Check permission from group
	if !mgr.HasPermission(playerID, "chat.color") {
		t.Error("Player should inherit permission from group")
	}

	// Remove player from group
	err = mgr.RemovePlayerFromGroup(playerID, "vip")
	if err != nil {
		t.Fatalf("Failed to remove player from group: %v", err)
	}

	// Check permission is gone
	if mgr.HasPermission(playerID, "chat.color") {
		t.Error("Player should not have group permission after removal")
	}

	// Try to remove from non-existent group membership
	err = mgr.RemovePlayerFromGroup(playerID, "vip")
	if err == nil {
		t.Error("Should fail when removing player from group they're not in")
	}
}

func TestPermissionManager_PlayerPermissions(t *testing.T) {
	defer cleanupPublicTest()
	mgr := createTestPermissionManager()
	playerID := uuid.New()

	// Add player to ensure they exist
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "default")
	if err != nil {
		t.Fatalf("Failed to add player to default group: %v", err)
	}

	// Add individual permission
	err = mgr.AddPlayerPermission(playerID, "custom.permission")
	if err != nil {
		t.Fatalf("Failed to add player permission: %v", err)
	}

	// Check permission
	if !mgr.HasPermission(playerID, "custom.permission") {
		t.Error("Player should have the added permission")
	}

	// Remove permission
	err = mgr.RemovePlayerPermission(playerID, "custom.permission")
	if err != nil {
		t.Fatalf("Failed to remove player permission: %v", err)
	}

	// Check permission is gone
	if mgr.HasPermission(playerID, "custom.permission") {
		t.Error("Player should not have the removed permission")
	}

	// Try to remove non-existent permission
	err = mgr.RemovePlayerPermission(playerID, "non.existent")
	if err != ErrPermissionNotFound {
		t.Errorf("Expected ErrPermissionNotFound, got %v", err)
	}
}

func TestPermissionManager_ErrorConversion(t *testing.T) {
	defer cleanupPublicTest()
	mgr := createTestPermissionManager()
	playerID := uuid.New()

	// Test player not found error
	err := mgr.AddPlayerPermission(playerID, "test.permission")
	if err != ErrPlayerNotFound {
		t.Errorf("Expected ErrPlayerNotFound, got %v", err)
	}

	// Test group not found error
	err = mgr.AddPlayerToGroup(playerID, "TestPlayer", "non-existent")
	if err != ErrGroupNotFound {
		t.Errorf("Expected ErrGroupNotFound, got %v", err)
	}
}

func TestPermissionManager_GroupPermissionManagement(t *testing.T) {
	defer cleanupPublicTest()
	mgr := createTestPermissionManager()
	playerID := uuid.New()

	// Create a test group
	err := mgr.CreateGroup("testers", []string{"test.basic"})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	// Add player to group
	err = mgr.AddPlayerToGroup(playerID, "TestPlayer", "testers")
	if err != nil {
		t.Fatalf("Failed to add player to group: %v", err)
	}

	// Add permission to group
	err = mgr.AddPermissionToGroup("testers", "test.advanced")
	if err != nil {
		t.Fatalf("Failed to add permission to group: %v", err)
	}

	// Check that player now has the new permission
	if !mgr.HasPermission(playerID, "test.advanced") {
		t.Error("Player should have permission added to their group")
	}

	// Check that player still has original permission
	if !mgr.HasPermission(playerID, "test.basic") {
		t.Error("Player should still have original group permission")
	}

	// Remove permission from group
	err = mgr.RemovePermissionFromGroup("testers", "test.advanced")
	if err != nil {
		t.Fatalf("Failed to remove permission from group: %v", err)
	}

	// Check that player no longer has the removed permission
	if mgr.HasPermission(playerID, "test.advanced") {
		t.Error("Player should not have permission removed from their group")
	}

	// Check that player still has original permission
	if !mgr.HasPermission(playerID, "test.basic") {
		t.Error("Player should still have original group permission")
	}

	// Try to add permission to non-existent group
	err = mgr.AddPermissionToGroup("non-existent", "some.permission")
	if err != ErrGroupNotFound {
		t.Errorf("Expected ErrGroupNotFound, got %v", err)
	}

	// Try to remove permission that group doesn't have
	err = mgr.RemovePermissionFromGroup("testers", "non.existent")
	if err != ErrPermissionNotFound {
		t.Errorf("Expected ErrPermissionNotFound, got %v", err)
	}
}

func TestPermissionManager_Save(t *testing.T) {
	defer cleanupPublicTest()
	mgr := createTestPermissionManager()

	// Save should work without error
	err := mgr.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Check that file was created
	if _, err := os.Stat("/tmp/test_permission_public.json"); os.IsNotExist(err) {
		t.Error("Save should create the storage file")
	}
}