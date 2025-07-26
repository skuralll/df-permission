package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/shared"
)

func TestNewPermissionCache(t *testing.T) {
	config := shared.CacheConfig{
		TTL:             10 * time.Second,
		CleanupInterval: 5 * time.Second,
		Enabled:         true,
	}

	cache := NewPermissionCache(config)

	if cache.ttl != 10*time.Second {
		t.Errorf("Expected TTL %v, got %v", 10*time.Second, cache.ttl)
	}
	if cache.cleanupInterval != 5*time.Second {
		t.Errorf("Expected cleanup interval %v, got %v", 5*time.Second, cache.cleanupInterval)
	}
	if !cache.enabled {
		t.Error("Expected cache to be enabled")
	}

	cache.Close()
}

func TestPermissionCache_SetAndGet(t *testing.T) {
	config := shared.CacheConfig{
		TTL:     30 * time.Second,
		Enabled: true,
	}
	cache := NewPermissionCache(config)
	defer cache.Close()

	playerID := uuid.New()
	permission := "test.permission"

	// キャッシュにエントリなし
	result, found := cache.Get(playerID, permission)
	if found {
		t.Error("Expected cache miss, but got hit")
	}

	// キャッシュに結果を保存
	cache.Set(playerID, permission, true)

	// キャッシュから取得
	result, found = cache.Get(playerID, permission)
	if !found {
		t.Error("Expected cache hit, but got miss")
	}
	if !result {
		t.Error("Expected true, got false")
	}
}

func TestPermissionCache_Disabled(t *testing.T) {
	config := shared.CacheConfig{
		Enabled: false,
	}
	cache := NewPermissionCache(config)

	playerID := uuid.New()
	permission := "test.permission"

	// 無効化されたキャッシュでは何も保存されない
	cache.Set(playerID, permission, true)
	_, found := cache.Get(playerID, permission)

	if found {
		t.Error("Expected cache miss for disabled cache")
	}
}

func TestPermissionCache_TTL(t *testing.T) {
	config := shared.CacheConfig{
		TTL:     100 * time.Millisecond,
		Enabled: true,
	}
	cache := NewPermissionCache(config)
	defer cache.Close()

	playerID := uuid.New()
	permission := "test.permission"

	// キャッシュに保存
	cache.Set(playerID, permission, true)

	// すぐに取得（有効）
	result, found := cache.Get(playerID, permission)
	if !found || !result {
		t.Error("Expected cache hit immediately after set")
	}

	// TTL後に取得（期限切れ）
	time.Sleep(150 * time.Millisecond)
	result, found = cache.Get(playerID, permission)
	if found {
		t.Error("Expected cache miss after TTL expiration")
	}
}

func TestPermissionCache_Clear(t *testing.T) {
	config := shared.CacheConfig{
		TTL:     30 * time.Second,
		Enabled: true,
	}
	cache := NewPermissionCache(config)
	defer cache.Close()

	playerID := uuid.New()
	permission := "test.permission"

	// キャッシュに保存
	cache.Set(playerID, permission, true)

	// 取得可能
	_, found := cache.Get(playerID, permission)
	if !found {
		t.Error("Expected cache hit before clear")
	}

	// クリア
	cache.Clear()

	// 取得不可
	_, found = cache.Get(playerID, permission)
	if found {
		t.Error("Expected cache miss after clear")
	}
}

func TestPermissionCache_InvalidatePlayer(t *testing.T) {
	config := shared.CacheConfig{
		TTL:     30 * time.Second,
		Enabled: true,
	}
	cache := NewPermissionCache(config)
	defer cache.Close()

	playerID1 := uuid.New()
	playerID2 := uuid.New()
	permission1 := "test.permission1"
	permission2 := "test.permission2"

	// 複数のエントリを保存
	cache.Set(playerID1, permission1, true)
	cache.Set(playerID1, permission2, false)
	cache.Set(playerID2, permission1, true)

	// プレイヤー1のキャッシュを無効化
	cache.InvalidatePlayer(playerID1)

	// プレイヤー1のエントリは削除される
	_, found := cache.Get(playerID1, permission1)
	if found {
		t.Error("Expected cache miss for invalidated player")
	}
	_, found = cache.Get(playerID1, permission2)
	if found {
		t.Error("Expected cache miss for invalidated player")
	}

	// プレイヤー2のエントリは残る
	_, found = cache.Get(playerID2, permission1)
	if !found {
		t.Error("Expected cache hit for non-invalidated player")
	}
}

func TestPermissionCache_InvalidatePermission(t *testing.T) {
	config := shared.CacheConfig{
		TTL:     30 * time.Second,
		Enabled: true,
	}
	cache := NewPermissionCache(config)
	defer cache.Close()

	playerID1 := uuid.New()
	playerID2 := uuid.New()
	permission1 := "test.permission1"
	permission2 := "test.permission2"

	// 複数のエントリを保存
	cache.Set(playerID1, permission1, true)
	cache.Set(playerID2, permission1, false)
	cache.Set(playerID1, permission2, true)

	// permission1を無効化
	cache.InvalidatePermission(permission1)

	// permission1のエントリは削除される
	_, found := cache.Get(playerID1, permission1)
	if found {
		t.Error("Expected cache miss for invalidated permission")
	}
	_, found = cache.Get(playerID2, permission1)
	if found {
		t.Error("Expected cache miss for invalidated permission")
	}

	// permission2のエントリは残る
	_, found = cache.Get(playerID1, permission2)
	if !found {
		t.Error("Expected cache hit for non-invalidated permission")
	}
}

func TestPermissionCache_EnabledState(t *testing.T) {
	config := shared.CacheConfig{
		TTL:     30 * time.Second,
		Enabled: false,
	}
	cache := NewPermissionCache(config)

	if cache.IsEnabled() {
		t.Error("Expected cache to be disabled")
	}

	// 有効化
	cache.SetEnabled(true)
	if !cache.IsEnabled() {
		t.Error("Expected cache to be enabled after SetEnabled(true)")
	}

	// 無効化
	cache.SetEnabled(false)
	if cache.IsEnabled() {
		t.Error("Expected cache to be disabled after SetEnabled(false)")
	}
}

func TestPermissionCache_AutoCleanup(t *testing.T) {
	config := shared.CacheConfig{
		TTL:             50 * time.Millisecond,
		CleanupInterval: 30 * time.Millisecond,
		Enabled:         true,
	}
	cache := NewPermissionCache(config)
	defer cache.Close()

	playerID1 := uuid.New()
	playerID2 := uuid.New()
	permission := "test.permission"

	// 複数のエントリを保存
	cache.Set(playerID1, permission, true)
	cache.Set(playerID2, permission, false)

	// エントリが存在することを確認
	_, found1 := cache.Get(playerID1, permission)
	_, found2 := cache.Get(playerID2, permission)
	if !found1 || !found2 {
		t.Error("Expected both cache entries to be present")
	}

	// TTL経過まで待機
	time.Sleep(60 * time.Millisecond)

	// 自動クリーンアップが実行されるまで待機
	time.Sleep(50 * time.Millisecond)

	// エントリが自動削除されることを確認
	_, found1 = cache.Get(playerID1, permission)
	_, found2 = cache.Get(playerID2, permission)
	if found1 || found2 {
		t.Error("Expected cache entries to be automatically cleaned up")
	}
}

func TestPermissionCache_CleanupStops(t *testing.T) {
	config := shared.CacheConfig{
		TTL:             100 * time.Millisecond,
		CleanupInterval: 20 * time.Millisecond,
		Enabled:         true,
	}
	cache := NewPermissionCache(config)

	playerID := uuid.New()
	permission := "test.permission"

	// エントリを保存
	cache.Set(playerID, permission, true)

	// キャッシュを閉じる
	cache.Close()

	// 少し待機してクリーンアップが停止することを確認
	time.Sleep(100 * time.Millisecond)

	// キャッシュが無効化されているか確認
	if cache.IsEnabled() {
		t.Error("Expected cache to be disabled after Close()")
	}
}

func TestPermissionCache_SetEnabledCleanup(t *testing.T) {
	config := shared.CacheConfig{
		TTL:             50 * time.Millisecond,
		CleanupInterval: 25 * time.Millisecond,
		Enabled:         false,
	}
	cache := NewPermissionCache(config)

	playerID := uuid.New()
	permission := "test.permission"

	// 無効状態では何も保存されない
	cache.Set(playerID, permission, true)
	_, found := cache.Get(playerID, permission)
	if found {
		t.Error("Expected no cache entry when disabled")
	}

	// 有効化
	cache.SetEnabled(true)

	// エントリを保存
	cache.Set(playerID, permission, true)
	_, found = cache.Get(playerID, permission)
	if !found {
		t.Error("Expected cache entry after enabling")
	}

	// TTL経過まで待機
	time.Sleep(60 * time.Millisecond)

	// 自動クリーンアップが実行されるまで待機
	time.Sleep(40 * time.Millisecond)

	// エントリが自動削除されることを確認
	_, found = cache.Get(playerID, permission)
	if found {
		t.Error("Expected cache entry to be automatically cleaned up after enabling")
	}

	// 無効化
	cache.SetEnabled(false)

	// 再度エントリを保存しようとしても保存されない
	cache.Set(playerID, permission, true)
	_, found = cache.Get(playerID, permission)
	if found {
		t.Error("Expected no cache entry after disabling")
	}
}
