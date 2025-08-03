package application

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/shared"
)

func createTestManager() *Manager {
	config := shared.ManagerConfig{
		AutoSave: false,
		Storage: shared.StorageConfig{
			Path: "/tmp/test_permissions.json",
		},
		Cache: shared.CacheConfig{
			Enabled:         false,
			TTL:             5 * time.Second,
			CleanupInterval: 10 * time.Second,
		},
	}
	return NewManager(config)
}

func cleanup() {
	os.Remove("/tmp/test_permissions.json")
}

func TestManagerCreation(t *testing.T) {
	defer cleanup()
	mgr := createTestManager()

	if mgr == nil {
		t.Fatal("Manager creation failed")
	}

	groups := mgr.GetAllGroups()
	if len(groups) != 2 {
		t.Fatalf("Expected 2 default groups, got %d", len(groups))
	}

	if _, exists := groups["default"]; !exists {
		t.Error("Default group not found")
	}

	if _, exists := groups["admin"]; !exists {
		t.Error("Admin group not found")
	}
}

func TestHasPermission(t *testing.T) {
	defer cleanup()
	mgr := createTestManager()
	playerID := uuid.New()

	if mgr.HasPermission(playerID, "any.permission") {
		t.Error("Non-existent player should not have any permissions")
	}

	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player to admin group: %v", err)
	}

	if !mgr.HasPermission(playerID, "any.permission") {
		t.Error("Admin should have all permissions")
	}

	if !mgr.HasPermission(playerID, "moderation.kick") {
		t.Error("Admin should have moderation.kick permission")
	}
}

func TestGroupManagement(t *testing.T) {
	defer cleanup()
	mgr := createTestManager()

	err := mgr.CreateGroup("vip", []string{"chat.color", "world.build.fast"})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	group := mgr.GetGroup("vip")
	if group == nil {
		t.Fatal("Created group not found")
	}

	if group.Name != "vip" {
		t.Errorf("Expected group name 'vip', got '%s'", group.Name)
	}

	if len(group.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(group.Permissions))
	}

	err = mgr.CreateGroup("vip", []string{})
	if err == nil {
		t.Error("Creating duplicate group should fail")
	}

	err = mgr.DeleteGroup("vip")
	if err != nil {
		t.Fatalf("Failed to delete group: %v", err)
	}

	if mgr.GetGroup("vip") != nil {
		t.Error("Group should be deleted")
	}
}

func TestPlayerManagement(t *testing.T) {
	defer cleanup()
	mgr := createTestManager()
	playerID := uuid.New()

	if mgr.PlayerExists(playerID) {
		t.Error("Player should not exist initially")
	}

	err := mgr.CreatePlayer(playerID, "TestPlayer")
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	if !mgr.PlayerExists(playerID) {
		t.Error("Player should exist after creation")
	}

	player := mgr.GetPlayerData(playerID)
	if player == nil {
		t.Fatal("Player data should not be nil")
	}

	if player.PlayerName != "TestPlayer" {
		t.Errorf("Expected player name 'TestPlayer', got '%s'", player.PlayerName)
	}

	err = mgr.CreatePlayer(playerID, "DuplicatePlayer")
	if err == nil {
		t.Error("Creating duplicate player should fail")
	}

	err = mgr.RemovePlayer(playerID)
	if err != nil {
		t.Fatalf("Failed to remove player: %v", err)
	}

	if mgr.PlayerExists(playerID) {
		t.Error("Player should not exist after removal")
	}
}

func TestManager_UpdatePlayerName(t *testing.T) {
	defer cleanup()
	mgr := createTestManager()
	playerID := uuid.New()

	err := mgr.CreatePlayer(playerID, "OriginalName")
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	err = mgr.UpdatePlayerName(playerID, "UpdatedName")
	if err != nil {
		t.Fatalf("Failed to update player name: %v", err)
	}

	player := mgr.GetPlayerData(playerID)
	if player == nil {
		t.Fatal("Player data should not be nil")
	}

	if player.PlayerName != "UpdatedName" {
		t.Errorf("Expected player name 'UpdatedName', got '%s'", player.PlayerName)
	}

	err = mgr.UpdatePlayerName(playerID, "UpdatedName")
	if err != nil {
		t.Fatalf("Updating with same name should not fail: %v", err)
	}

	nonExistentPlayerID := uuid.New()
	err = mgr.UpdatePlayerName(nonExistentPlayerID, "NewName")
	if err == nil {
		t.Error("Updating non-existent player should fail")
	}

	if err != nil && !errors.Is(err, shared.ErrPlayerNotFound) {
		t.Errorf("Expected PlayerNotFoundError, got %T: %v", err, err)
	}
}

func TestPlayerPermissions(t *testing.T) {
	defer cleanup()
	mgr := createTestManager()
	playerID := uuid.New()

	err := mgr.CreatePlayer(playerID, "TestPlayer")
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	err = mgr.AddPlayerPermission(playerID, "custom.permission")
	if err != nil {
		t.Fatalf("Failed to add player permission: %v", err)
	}

	if !mgr.HasPermission(playerID, "custom.permission") {
		t.Error("Player should have the added permission")
	}

	permissions := mgr.GetPlayerPermissions(playerID)
	if len(permissions) != 1 {
		t.Errorf("Expected 1 permission, got %d", len(permissions))
	}

	err = mgr.RemovePlayerPermission(playerID, "custom.permission")
	if err != nil {
		t.Fatalf("Failed to remove player permission: %v", err)
	}

	if mgr.HasPermission(playerID, "custom.permission") {
		t.Error("Player should not have the removed permission")
	}
}

func TestPlayerGroupMembership(t *testing.T) {
	defer cleanup()
	mgr := createTestManager()
	playerID := uuid.New()

	err := mgr.CreateGroup("vip", []string{"chat.color"})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	err = mgr.AddPlayerToGroup(playerID, "TestPlayer", "vip")
	if err != nil {
		t.Fatalf("Failed to add player to group: %v", err)
	}

	groups := mgr.GetPlayerGroups(playerID)
	if len(groups) != 1 || groups[0] != "vip" {
		t.Errorf("Expected player to be in 'vip' group, got %v", groups)
	}

	if !mgr.HasPermission(playerID, "chat.color") {
		t.Error("Player should inherit permission from group")
	}

	err = mgr.RemovePlayerFromGroup(playerID, "vip")
	if err != nil {
		t.Fatalf("Failed to remove player from group: %v", err)
	}

	groups = mgr.GetPlayerGroups(playerID)
	if len(groups) != 0 {
		t.Errorf("Expected player to have no groups, got %v", groups)
	}

	if mgr.HasPermission(playerID, "chat.color") {
		t.Error("Player should not have group permission after removal")
	}
}

func TestWildcardPermissions(t *testing.T) {
	defer cleanup()
	mgr := createTestManager()
	playerID := uuid.New()

	err := mgr.CreateGroup("moderator", []string{"moderation.*"})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	err = mgr.AddPlayerToGroup(playerID, "TestPlayer", "moderator")
	if err != nil {
		t.Fatalf("Failed to add player to group: %v", err)
	}

	if !mgr.HasPermission(playerID, "moderation.kick") {
		t.Error("Player should have moderation.kick via wildcard")
	}

	if !mgr.HasPermission(playerID, "moderation.ban") {
		t.Error("Player should have moderation.ban via wildcard")
	}

	if mgr.HasPermission(playerID, "chat.color") {
		t.Error("Player should not have chat.color (not in moderation.*)")
	}
}

func TestSystemGroupProtection(t *testing.T) {
	defer cleanup()
	mgr := createTestManager()

	err := mgr.DeleteGroup("default")
	if err == nil {
		t.Error("Should not be able to delete default group")
	}

	err = mgr.DeleteGroup("admin")
	if err == nil {
		t.Error("Should not be able to delete admin group")
	}
}

func TestCacheManagement(t *testing.T) {
	defer cleanup()
	config := shared.ManagerConfig{
		AutoSave: false,
		Storage: shared.StorageConfig{
			Path: "/tmp/test_permissions.json",
		},
		Cache: shared.CacheConfig{
			Enabled:         true,
			TTL:             5 * time.Second,
			CleanupInterval: 10 * time.Second,
		},
	}
	mgr := NewManager(config)

	mgr.ClearCache()

	mgr.SetCacheEnabled(false)
	mgr.SetCacheEnabled(true)
}