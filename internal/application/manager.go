package application

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	dfpermission "github.com/skuralll/df-permission"
	"github.com/skuralll/df-permission/internal/domain"
	"github.com/skuralll/df-permission/internal/shared"
)

type Manager struct {
	// config
	autoSave bool
	// internal state
	storage domain.PermissionRepository
	groups  map[string]*shared.Group
	players map[uuid.UUID]*shared.PlayerData
	cache   *domain.PermissionCache
	checker *domain.PermissionChecker
	mutex   sync.RWMutex
}

func NewManager(config dfpermission.ManagerConfig) *Manager {
	storage := *domain.NewPermissionRepository(config.Storage)
	cache := domain.NewPermissionCache(config.Cache)
	checker := domain.NewPermissionChecker()

	mgr := &Manager{
		autoSave: config.AutoSave,
		storage:  storage,
		groups:   make(map[string]*shared.Group),
		players:  make(map[uuid.UUID]*shared.PlayerData),
		cache:    cache,
		checker:  checker,
		mutex:    sync.RWMutex{},
	}

	// ストレージからデータを読み込む
	mgr.initializeDefaultGroups()

	// 既存データをロード
	mgr.loadData()

	return mgr
}

// 現在のパーミッションデータをストレージに保存
func (mgr *Manager) Save() error {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	data := &shared.PermissionData{
		Groups:  mgr.groups,
		Players: mgr.players,
		Meta: &shared.Metadata{
			Version:   shared.PermissionDataVersion,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return mgr.storage.Save(data)
}

// ストレージからパーミッションデータを再読み込みする
func (mgr *Manager) Reload() error {
	return mgr.loadData()
}

// すべてのキャッシュされたパーミッション結果をクリアする
func (mgr *Manager) ClearCache() {
	if mgr.cache != nil {
		mgr.cache.Clear()
	}
}

// キャッシュの有効・無効を切り替える
func (mgr *Manager) SetCacheEnabled(enabled bool) {
	if mgr.cache != nil {
		mgr.cache.SetEnabled(enabled)
	}
}

// オートセーブの有効・無効を切り替える
func (mgr *Manager) SetAutoSave(enabled bool) {
	mgr.autoSave = enabled
}

// プレイヤーが特定のパーミッションを持っているかどうかを確認する
// 結果はキャッシュされる
func (mgr *Manager) HasPermission(playerID uuid.UUID, permission string) bool {
	// キャッシュから結果を取得
	if mgr.cache != nil && mgr.cache.IsEnabled() {
		if result, found := mgr.cache.Get(playerID, permission); found {
			return result
		}
	}

	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	// プレイヤーデータ取得
	player, exists := mgr.players[playerID]
	if !exists {
		// キャッシュ
		if mgr.cache != nil && mgr.cache.IsEnabled() {
			mgr.cache.Set(playerID, permission, false)
		}
		return false
	}

	// チェッカー用のグループマップを構築
	groupsMap := make(map[string][]string)
	for name, group := range mgr.groups {
		groupsMap[name] = group.Permissions
	}

	// チェッカーを使用してパーミッションを確認
	result := mgr.checker.HasPermission(
		player.Permissions,
		player.Groups,
		groupsMap,
		permission,
	)

	// 結果をキャッシュ
	if mgr.cache != nil && mgr.cache.IsEnabled() {
		mgr.cache.Set(playerID, permission, result)
	}

	return result
}

// プレイヤーをパーミッショングループに追加
func (mgr *Manager) AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// すでにグループが存在するか確認
	if _, exists := mgr.groups[groupName]; !exists {
		return fmt.Errorf("group '%s' does not exist", groupName)
	}

	// プレイヤーデータを取得または作成
	player, exists := mgr.players[playerID]
	if !exists {
		player = &shared.PlayerData{
			PlayerID:    playerID,
			PlayerName:  playerName,
			Groups:      []string{},
			Permissions: []string{},
			UpdatedAt:   time.Now(),
		}
		mgr.players[playerID] = player
	}

	// グループにすでに参加しているか確認
	for _, group := range player.Groups {
		if group == groupName {
			return nil
		}
	}

	// グループを追加
	player.Groups = append(player.Groups, groupName)
	player.UpdatedAt = time.Now()

	// キャッシュを無効化
	if mgr.cache != nil && mgr.cache.IsEnabled() {
		mgr.cache.InvalidatePlayer(playerID)
	}

	// セーブする
	if mgr.autoSave {
		go mgr.Save()
	}

	return nil
}

// プレイヤーをパーミッショングループから削除する
// プレイヤーがそのグループにいない場合はエラーを返す
func (mgr *Manager) RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	player, exists := mgr.players[playerID]
	if !exists {
		return fmt.Errorf("player not found")
	}

	// グループを探して削除
	for i, group := range player.Groups {
		if group == groupName {
			player.Groups = append(player.Groups[:i], player.Groups[i+1:]...)
			player.UpdatedAt = time.Now()

			// このプレイヤーのキャッシュを無効化（グループ所属が変更されたため）
			if mgr.cache != nil && mgr.cache.IsEnabled() {
				mgr.cache.InvalidatePlayer(playerID)
			}

			if mgr.autoSave {
				go mgr.Save()
			}
			return nil
		}
	}

	return fmt.Errorf("player is not in group '%s'", groupName)
}

// =============================================================================
// 内部実装
// =============================================================================

// デフォルトのパーミッショングループを作成する
func (mgr *Manager) initializeDefaultGroups() {
	// デフォルトグループが存在しない場合は作成
	if _, exists := mgr.groups["default"]; !exists {
		mgr.groups["default"] = &shared.Group{
			Name:        "default",
			Permissions: []string{"chat.send", "world.interact"},
		}
	}

	// 管理者グループが存在しない場合は作成
	if _, exists := mgr.groups["admin"]; !exists {
		mgr.groups["admin"] = &shared.Group{
			Name:        "admin",
			Permissions: []string{"*"},
		}
	}
}

// ストレージからパーミッションデータを読み込む
func (mgr *Manager) loadData() error {
	data, err := mgr.storage.Load()
	if err != nil {
		return err
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// グループを読み込む
	if data.Groups != nil {
		for name, group := range data.Groups {
			mgr.groups[name] = group
		}
	}

	// プレイヤーを読み込む
	if data.Players != nil {
		for id, player := range data.Players {
			mgr.players[id] = player
		}
	}

	return nil
}
