package domain

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/repository"
	"github.com/skuralll/df-permission/internal/shared"
)

type PermissionService struct {
	groups  map[string]*shared.Group
	players map[uuid.UUID]*shared.PlayerData
	storage repository.Storage
	cache   *PermissionCache
	checker *PermissionChecker
	mutex   sync.RWMutex
}

func NewPermissionService(config shared.ServiceConfig) *PermissionService {
	storage := repository.NewJSONStorage(config.Storage)
	cache := NewPermissionCache(config.Cache)
	checker := NewPermissionChecker()

	return &PermissionService{
		groups:  make(map[string]*shared.Group),
		players: make(map[uuid.UUID]*shared.PlayerData),
		storage: storage,
		cache:   cache,
		checker: checker,
		mutex:   sync.RWMutex{},
	}
}

// プレイヤーが特定のパーミッションを持っているかどうかを確認する
// 結果はキャッシュされる
func (svc *PermissionService) HasPermission(playerID uuid.UUID, permission string) bool {
	// キャッシュから結果を取得
	if svc.cache != nil && svc.cache.IsEnabled() {
		if result, found := svc.cache.Get(playerID, permission); found {
			return result
		}
	}

	svc.mutex.RLock()
	defer svc.mutex.RUnlock()

	// プレイヤーデータ取得
	player, exists := svc.players[playerID]
	if !exists {
		// キャッシュ
		if svc.cache != nil && svc.cache.IsEnabled() {
			svc.cache.Set(playerID, permission, false)
		}
		return false
	}

	// チェッカー用のグループマップを構築
	groupsMap := make(map[string][]string)
	for name, group := range svc.groups {
		groupsMap[name] = group.Permissions
	}

	// チェッカーを使用してパーミッションを確認
	result := svc.checker.HasPermission(
		player.Permissions,
		player.Groups,
		groupsMap,
		permission,
	)

	// 結果をキャッシュ
	if svc.cache != nil && svc.cache.IsEnabled() {
		svc.cache.Set(playerID, permission, result)
	}

	return result
}

// プレイヤーをパーミッショングループに追加
func (svc *PermissionService) AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	// すでにグループが存在するか確認
	if _, exists := svc.groups[groupName]; !exists {
		return fmt.Errorf("group '%s' does not exist", groupName)
	}

	// プレイヤーデータを取得または作成
	player, exists := svc.players[playerID]
	if !exists {
		player = &shared.PlayerData{
			PlayerID:    playerID,
			PlayerName:  playerName,
			Groups:      []string{},
			Permissions: []string{},
			UpdatedAt:   time.Now(),
		}
		svc.players[playerID] = player
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
	if svc.cache != nil && svc.cache.IsEnabled() {
		svc.cache.InvalidatePlayer(playerID)
	}

	// セーブする
	// if svc.settings.AutoSave {
	go svc.Save()
	// }

	return nil
}

// プレイヤーをパーミッショングループから削除する
// プレイヤーがそのグループにいない場合はエラーを返す
func (svc *PermissionService) RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	player, exists := svc.players[playerID]
	if !exists {
		return fmt.Errorf("player not found")
	}

	// グループを探して削除
	for i, group := range player.Groups {
		if group == groupName {
			player.Groups = append(player.Groups[:i], player.Groups[i+1:]...)
			player.UpdatedAt = time.Now()

			// このプレイヤーのキャッシュを無効化（グループ所属が変更されたため）
			if svc.cache != nil && svc.cache.IsEnabled() {
				svc.cache.InvalidatePlayer(playerID)
			}

			// if m.settings.AutoSave {
			go svc.Save()
			// }
			return nil
		}
	}

	return fmt.Errorf("player is not in group '%s'", groupName)
}

// 指定されたパーミッションを持つ新しいパーミッショングループを作成する
// グループがすでに存在する場合はエラーを返す
func (svc *PermissionService) CreateGroup(name string, permissions []string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	if _, exists := svc.groups[name]; exists {
		return fmt.Errorf("group '%s' already exists", name)
	}

	svc.groups[name] = &shared.Group{
		Name:        name,
		Permissions: permissions,
	}

	// if m.settings.AutoSave {
	go svc.Save()
	// }

	return nil
}

// 現在のパーミッションデータをストレージに保存
func (svc *PermissionService) Save() error {
	svc.mutex.RLock()
	defer svc.mutex.RUnlock()

	data := &shared.PermissionData{
		Groups:  svc.groups,
		Players: svc.players,
		Meta: &shared.Metadata{
			Version:   shared.PermissionDataVersion,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return svc.storage.Save(data)
}

// =============================================================================
// 内部実装
// =============================================================================

// デフォルトのパーミッショングループを作成する
func (svc *PermissionService) initializeDefaultGroups() {
	// デフォルトグループが存在しない場合は作成
	if _, exists := svc.groups["default"]; !exists {
		svc.groups["default"] = &shared.Group{
			Name:        "default",
			Permissions: []string{"chat.send", "world.interact"},
		}
	}

	// 管理者グループが存在しない場合は作成
	if _, exists := svc.groups["admin"]; !exists {
		svc.groups["admin"] = &shared.Group{
			Name:        "admin",
			Permissions: []string{"*"},
		}
	}
}

// ストレージからパーミッションデータを読み込む
func (svc *PermissionService) loadData() error {
	data, err := svc.storage.Load()
	if err != nil {
		return err
	}

	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	// グループを読み込む
	if data.Groups != nil {
		for name, group := range data.Groups {
			svc.groups[name] = group
		}
	}

	// プレイヤーを読み込む
	if data.Players != nil {
		for id, player := range data.Players {
			svc.players[id] = player
		}
	}

	return nil
}
