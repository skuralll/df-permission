package repository

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJSONStorage_Exists(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewJSONStorage(StorageConfig{
		Path: filepath.Join(tempDir, "test.json"),
	})

	// ファイルが存在しない場合
	if storage.Exists() {
		t.Error("Expected file to not exist")
	}

	// ファイルを作成
	err := os.WriteFile(storage.config.Path, []byte("{}"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// ファイルが存在する場合
	if !storage.Exists() {
		t.Error("Expected file to exist")
	}
}

func TestJSONStorage_Save_Load(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewJSONStorage(StorageConfig{
		Path: filepath.Join(tempDir, "test.json"),
	})

	// テストデータを作成
	playerID := uuid.New()
	testData := &PermissionData{
		Groups: map[string]*Group{
			"admin": {
				Name:        "admin",
				Permissions: []string{"*"},
			},
		},
		Players: map[uuid.UUID]*PlayerData{
			playerID: {
				PlayerID:    playerID,
				PlayerName:  "TestPlayer",
				Groups:      []string{"admin"},
				Permissions: []string{"test.permission"},
				UpdatedAt:   time.Now(),
			},
		},
		Meta: &Metadata{
			Version:   PermissionDataVersion,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// 保存
	err := storage.Save(testData)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// ファイルが存在することを確認
	if !storage.Exists() {
		t.Error("File should exist after save")
	}

	// 読み込み
	loadedData, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// データの検証
	if loadedData.Meta.Version != PermissionDataVersion {
		t.Errorf("Expected version %s, got %s", PermissionDataVersion, loadedData.Meta.Version)
	}

	if len(loadedData.Groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(loadedData.Groups))
	}

	if adminGroup, exists := loadedData.Groups["admin"]; !exists {
		t.Error("Admin group not found")
	} else if len(adminGroup.Permissions) != 1 || adminGroup.Permissions[0] != "*" {
		t.Error("Admin group permissions mismatch")
	}

	if len(loadedData.Players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(loadedData.Players))
	}

	if player, exists := loadedData.Players[playerID]; !exists {
		t.Error("Test player not found")
	} else if player.PlayerName != "TestPlayer" {
		t.Errorf("Expected player name TestPlayer, got %s", player.PlayerName)
	}
}

func TestJSONStorage_Load_NonExistentFile(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewJSONStorage(StorageConfig{
		Path: filepath.Join(tempDir, "nonexistent.json"),
	})

	// 存在しないファイルの読み込み
	data, err := storage.Load()
	if err != nil {
		t.Fatalf("Load should not fail for non-existent file: %v", err)
	}

	// デフォルトデータが返されることを確認
	if data == nil {
		t.Fatal("Expected default data, got nil")
	}

	if data.Meta.Version != PermissionDataVersion {
		t.Errorf("Expected version %s, got %s", PermissionDataVersion, data.Meta.Version)
	}

	if len(data.Groups) != 0 {
		t.Errorf("Expected empty groups, got %d", len(data.Groups))
	}

	if len(data.Players) != 0 {
		t.Errorf("Expected empty players, got %d", len(data.Players))
	}
}
