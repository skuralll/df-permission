package permission

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func cleanupOptionsTest() {
	os.Remove("/tmp/test_options.json")
	os.Remove("/tmp/custom_path.json")
}

func TestNewManager_DefaultConfig(t *testing.T) {
	defer cleanupOptionsTest()

	// デフォルト設定でマネージャーを作成
	mgr := NewManager()
	if mgr == nil {
		t.Fatal("NewManager should create a valid manager with default config")
	}

	// デフォルト設定の動作確認（基本的な操作ができるかテスト）
	playerID := uuid.New()
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player to admin group with default config: %v", err)
	}

	if !mgr.HasPermission(playerID, "any.permission") {
		t.Error("Admin should have all permissions with default config")
	}
}

func TestWithStorage(t *testing.T) {
	defer cleanupOptionsTest()

	customPath := "/tmp/custom_path.json"
	mgr := NewManager(
		WithStorage(customPath),
	)

	playerID := uuid.New()
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// 保存してファイルが作成されるか確認
	err = mgr.Save()
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// カスタムパスにファイルが作成されているか確認
	if _, err := os.Stat(customPath); os.IsNotExist(err) {
		t.Errorf("Custom storage path file should be created at %s", customPath)
	}
}

func TestWithAutoSave(t *testing.T) {
	defer cleanupOptionsTest()

	// オートセーブ無効でマネージャーを作成
	mgr := NewManager(
		WithStorage("/tmp/test_options.json"),
		WithAutoSave(false),
	)

	playerID := uuid.New()
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// オートセーブが無効なので、ファイルは作成されないはず
	time.Sleep(100 * time.Millisecond) // 少し待機
	if _, err := os.Stat("/tmp/test_options.json"); !os.IsNotExist(err) {
		t.Error("File should not be auto-saved when AutoSave is disabled")
	}
}

func TestWithCacheEnabled(t *testing.T) {
	defer cleanupOptionsTest()

	// キャッシュ無効でマネージャーを作成
	mgr := NewManager(
		WithStorage("/tmp/test_options.json"),
		WithCacheEnabled(false),
	)

	playerID := uuid.New()
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// 権限チェックは正常に動作するはず（キャッシュなしでも）
	if !mgr.HasPermission(playerID, "any.permission") {
		t.Error("Admin should have all permissions even with cache disabled")
	}
}

func TestWithCache(t *testing.T) {
	defer cleanupOptionsTest()

	customTTL := 5 * time.Second
	mgr := NewManager(
		WithStorage("/tmp/test_options.json"),
		WithCache(customTTL),
	)

	playerID := uuid.New()
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// WithCacheを使うとキャッシュが有効になるはず
	if !mgr.HasPermission(playerID, "any.permission") {
		t.Error("Admin should have all permissions with custom cache TTL")
	}
}

func TestWithCacheCleanup(t *testing.T) {
	defer cleanupOptionsTest()

	customInterval := 30 * time.Second
	mgr := NewManager(
		WithStorage("/tmp/test_options.json"),
		WithCacheCleanup(customInterval),
	)

	playerID := uuid.New()
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// 基本的な動作確認
	if !mgr.HasPermission(playerID, "any.permission") {
		t.Error("Admin should have all permissions with custom cleanup interval")
	}
}

func TestMultipleOptions(t *testing.T) {
	defer cleanupOptionsTest()

	// 複数のオプションを組み合わせてテスト
	mgr := NewManager(
		WithStorage("/tmp/test_options.json"),
		WithAutoSave(true),
		WithCache(45*time.Second),
		WithCacheCleanup(2*time.Minute),
	)

	playerID := uuid.New()
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	if !mgr.HasPermission(playerID, "any.permission") {
		t.Error("Admin should have all permissions with multiple options")
	}

	// オートセーブが有効なので、少し待つとファイルが作成されるはず
	time.Sleep(100 * time.Millisecond)
	if _, err := os.Stat("/tmp/test_options.json"); os.IsNotExist(err) {
		t.Error("File should be auto-saved when AutoSave is enabled")
	}
}

func TestOptionOrder(t *testing.T) {
	defer cleanupOptionsTest()

	// オプションの適用順序をテスト（後から適用されるオプションが優先される）
	mgr := NewManager(
		WithAutoSave(false),
		WithAutoSave(true), // この設定が最終的に適用される
		WithStorage("/tmp/test_options.json"),
	)

	playerID := uuid.New()
	err := mgr.AddPlayerToGroup(playerID, "TestPlayer", "admin")
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// 最後のWithAutoSave(true)が適用されているはず
	time.Sleep(100 * time.Millisecond)
	if _, err := os.Stat("/tmp/test_options.json"); os.IsNotExist(err) {
		t.Error("File should be auto-saved as the last WithAutoSave(true) should take effect")
	}
}
