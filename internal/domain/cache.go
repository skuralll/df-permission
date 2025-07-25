package domain

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// 権限チェック結果をメモリに保存し、パフォーマンス向上を狙うためのキャッシュシステム
// スレッドセーフ, TTL対応, 自動クリーンアップ
type PermissionCache struct {
	data            map[string]*cacheEntry
	ttl             time.Duration
	cleanupInterval time.Duration
	mutex           sync.RWMutex
	stopCleanup     chan struct{}
	enabled         bool
}

// キャッシュエントリ
type cacheEntry struct {
	result    bool
	expiresAt time.Time
	createdAt time.Time
}

// キャッシュの設定
type CacheConfig struct {
	TTL             time.Duration
	CleanupInterval time.Duration
	Enabled         bool
}

func NewPermissionCache(config CacheConfig) *PermissionCache {
	// デフォルト設定
	if config.TTL == 0 {
		config.TTL = 30 * time.Second
	}
	if config.CleanupInterval == 0 {
		config.CleanupInterval = time.Minute
	}

	cache := &PermissionCache{
		data:            make(map[string]*cacheEntry),
		ttl:             config.TTL,
		cleanupInterval: config.CleanupInterval,
		mutex:           sync.RWMutex{},
		stopCleanup:     make(chan struct{}),
		enabled:         config.Enabled,
	}

	// 設定されていれば自動クリーンアップを開始
	if config.Enabled {
		go cache.cleanupLoop()
	}

	return cache
}

// 指定されたプレイヤーIDとパーミッションに対するキャッシュ済みの権限チェック結果を返す
// TTL: キャッシュエントリの有効期限を確認し、期限切れの場合はキャッシュミスとして扱う
func (pc *PermissionCache) Get(playerID uuid.UUID, permission string) (result bool, found bool) {
	if !pc.enabled {
		return false, false
	}

	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	key := pc.generateKey(playerID, permission)
	entry, exists := pc.data[key]

	if !exists {
		return false, false
	}

	// キャッシュエントリが期限切れかどうかを確認
	if time.Now().After(entry.expiresAt) {
		return false, false
	}
	return entry.result, true
}

// 権限チェック結果をTTL付きでキャッシュに保存
func (pc *PermissionCache) Set(playerID uuid.UUID, permission string, result bool) {
	if !pc.enabled {
		return
	}

	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	key := pc.generateKey(playerID, permission)

	pc.data[key] = &cacheEntry{
		result:    result,
		expiresAt: time.Now().Add(pc.ttl),
		createdAt: time.Now(),
	}
}

// 全てのキャッシュエントリを削除
func (pc *PermissionCache) Clear() {
	if !pc.enabled {
		return
	}

	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	pc.data = make(map[string]*cacheEntry)
}

// 指定されたプレイヤーに関連するすべてのキャッシュエントリを削除
func (pc *PermissionCache) InvalidatePlayer(playerID uuid.UUID) {
	if !pc.enabled {
		return
	}

	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	playerPrefix := playerID.String() + ":"
	removed := 0

	for key := range pc.data {
		if len(key) > len(playerPrefix) && key[:len(playerPrefix)] == playerPrefix {
			delete(pc.data, key)
			removed++
		}
	}
}

// 指定されたパーミッションに関連するすべてのキャッシュエントリを削除
func (pc *PermissionCache) InvalidatePermission(permission string) {
	if !pc.enabled {
		return
	}

	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	permissionSuffix := ":" + permission
	removed := 0

	for key := range pc.data {
		if len(key) > len(permissionSuffix) &&
			key[len(key)-len(permissionSuffix):] == permissionSuffix {
			delete(pc.data, key)
			removed++
		}
	}
}

// バックグラウンドのクリーンアップ用ゴルーチンを停止し、リソースを解放する
func (pc *PermissionCache) Close() {
	if !pc.enabled {
		return
	}

	// クリーンアップ用ゴルーチンに停止を通知
	select {
	case pc.stopCleanup <- struct{}{}:
	default:
		// チャネルがすでに閉じられている可能性がある
	}

	pc.Clear()
	pc.enabled = false
}

// キャッシングが有効かどうかを取得
func (pc *PermissionCache) IsEnabled() bool {
	return pc.enabled
}

// キャッシュの有効/無効を切り替える
func (pc *PermissionCache) SetEnabled(enabled bool) {
	if pc.enabled == enabled {
		return
	}

	pc.enabled = enabled

	if enabled {
		// 自動クリーンアップを開始
		go pc.cleanupLoop()
	} else {
		// 自動クリーンアップを停止、キャッシュクリア
		pc.Close()
	}
}

// プレイヤーIDとパーミッションから一意なキャッシュキーを生成
// フォーマット: "{playerUUID}:{permission}"
func (pc *PermissionCache) generateKey(playerID uuid.UUID, permission string) string {
	return fmt.Sprintf("%s:%s", playerID.String(), permission)
}

// 自動クリーンアップループ
func (pc *PermissionCache) cleanupLoop() {
	ticker := time.NewTicker(pc.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pc.cleanup()
		case <-pc.stopCleanup:
			return
		}
	}
}

// クリーンアップ
func (pc *PermissionCache) cleanup() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	now := time.Now()
	evicted := 0

	for key, entry := range pc.data {
		if now.After(entry.expiresAt) {
			delete(pc.data, key)
			evicted++
		}
	}
}
