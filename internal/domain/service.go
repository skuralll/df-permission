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
	// config
	autoSave bool
	// internal state
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

	svc := &PermissionService{
		autoSave: config.AutoSave,
		groups:   make(map[string]*shared.Group),
		players:  make(map[uuid.UUID]*shared.PlayerData),
		storage:  storage,
		cache:    cache,
		checker:  checker,
		mutex:    sync.RWMutex{},
	}

	// ストレージからデータを読み込む
	svc.initializeDefaultGroups()

	// 既存データをロード
	svc.loadData()

	return svc
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

// ストレージからパーミッションデータを再読み込みする
func (svc *PermissionService) Reload() error {
	return svc.loadData()
}

// すべてのキャッシュされたパーミッション結果をクリアする
func (svc *PermissionService) ClearCache() {
	if svc.cache != nil {
		svc.cache.Clear()
	}
}

// キャッシュの有効・無効を切り替える
func (svc *PermissionService) SetCacheEnabled(enabled bool) {
	if svc.cache != nil {
		svc.cache.SetEnabled(enabled)
	}
}

// オートセーブの有効・無効を切り替える
func (svc *PermissionService) SetAutoSave(enabled bool) {
	svc.autoSave = enabled
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
	if svc.autoSave {
		go svc.Save()
	}

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

			if svc.autoSave {
				go svc.Save()
			}
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

	if svc.autoSave {
		go svc.Save()
	}

	return nil
}

// パーミッショングループを削除する
func (svc *PermissionService) DeleteGroup(name string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	// デフォルトグループは削除できない
	if name == "default" || name == "admin" {
		return fmt.Errorf("cannot delete system group '%s'", name)
	}

	// グループが存在するかチェック
	if _, exists := svc.groups[name]; !exists {
		return fmt.Errorf("group '%s' does not exist", name)
	}

	// グループを削除
	delete(svc.groups, name)

	// このグループに所属していたプレイヤーからグループを削除
	for _, player := range svc.players {
		for i, groupName := range player.Groups {
			if groupName == name {
				player.Groups = append(player.Groups[:i], player.Groups[i+1:]...)
				player.UpdatedAt = time.Now()
				
				// このプレイヤーのキャッシュを無効化
				if svc.cache != nil && svc.cache.IsEnabled() {
					svc.cache.InvalidatePlayer(player.PlayerID)
				}
				break
			}
		}
	}

	// オートセーブ
	if svc.autoSave {
		go svc.Save()
	}

	return nil
}

// すべてのパーミッショングループのコピーを返す
func (svc *PermissionService) GetAllGroups() map[string]*shared.Group {
	svc.mutex.RLock()
	defer svc.mutex.RUnlock()

	result := make(map[string]*shared.Group)
	for name, group := range svc.groups {
		result[name] = &shared.Group{
			Name:        group.Name,
			Permissions: append([]string{}, group.Permissions...),
		}
	}
	return result
}

// パーミッショングループの権限を更新する
func (svc *PermissionService) UpdateGroup(name string, permissions []string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	// グループが存在するかチェック
	group, exists := svc.groups[name]
	if !exists {
		return fmt.Errorf("group '%s' does not exist", name)
	}

	// 権限を更新
	group.Permissions = append([]string{}, permissions...)

	// このグループに所属するプレイヤーのキャッシュを無効化
	for _, player := range svc.players {
		for _, groupName := range player.Groups {
			if groupName == name {
				if svc.cache != nil && svc.cache.IsEnabled() {
					svc.cache.InvalidatePlayer(player.PlayerID)
				}
				break
			}
		}
	}

	// オートセーブ
	if svc.autoSave {
		go svc.Save()
	}

	return nil
}

// 特定のパーミッショングループを取得する
func (svc *PermissionService) GetGroup(name string) *shared.Group {
	svc.mutex.RLock()
	defer svc.mutex.RUnlock()

	group, exists := svc.groups[name]
	if !exists {
		return nil
	}

	return &shared.Group{
		Name:        group.Name,
		Permissions: append([]string{}, group.Permissions...),
	}
}

// プレイヤーデータのコピーを返す
func (svc *PermissionService) GetPlayerData(playerID uuid.UUID) *shared.PlayerData {
	svc.mutex.RLock()
	defer svc.mutex.RUnlock()

	player, exists := svc.players[playerID]
	if !exists {
		return nil
	}

	return &shared.PlayerData{
		PlayerID:    player.PlayerID,
		PlayerName:  player.PlayerName,
		Groups:      append([]string{}, player.Groups...),
		Permissions: append([]string{}, player.Permissions...),
		UpdatedAt:   player.UpdatedAt,
	}
}

// プレイヤーに個別権限を追加する
func (svc *PermissionService) AddPlayerPermission(playerID uuid.UUID, permission string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	// プレイヤーデータを取得または作成
	player, exists := svc.players[playerID]
	if !exists {
		return fmt.Errorf("player with ID %s not found", playerID.String())
	}

	// 既に権限を持っているかチェック
	for _, perm := range player.Permissions {
		if perm == permission {
			return nil // 既に権限を持っている場合は何もしない
		}
	}

	// 権限を追加
	player.Permissions = append(player.Permissions, permission)
	player.UpdatedAt = time.Now()

	// キャッシュを無効化
	if svc.cache != nil && svc.cache.IsEnabled() {
		svc.cache.InvalidatePlayer(playerID)
	}

	// オートセーブ
	if svc.autoSave {
		go svc.Save()
	}

	return nil
}

// プレイヤーから個別権限を削除する
func (svc *PermissionService) RemovePlayerPermission(playerID uuid.UUID, permission string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	player, exists := svc.players[playerID]
	if !exists {
		return fmt.Errorf("player with ID %s not found", playerID.String())
	}

	// 権限を探して削除
	for i, perm := range player.Permissions {
		if perm == permission {
			player.Permissions = append(player.Permissions[:i], player.Permissions[i+1:]...)
			player.UpdatedAt = time.Now()

			// キャッシュを無効化
			if svc.cache != nil && svc.cache.IsEnabled() {
				svc.cache.InvalidatePlayer(playerID)
			}

			// オートセーブ
			if svc.autoSave {
				go svc.Save()
			}

			return nil
		}
	}

	return fmt.Errorf("player does not have permission '%s'", permission)
}

// プレイヤーの権限を一括設定する
func (svc *PermissionService) SetPlayerPermissions(playerID uuid.UUID, permissions []string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	player, exists := svc.players[playerID]
	if !exists {
		return fmt.Errorf("player with ID %s not found", playerID.String())
	}

	// 権限を置き換え
	player.Permissions = append([]string{}, permissions...)
	player.UpdatedAt = time.Now()

	// キャッシュを無効化
	if svc.cache != nil && svc.cache.IsEnabled() {
		svc.cache.InvalidatePlayer(playerID)
	}

	// オートセーブ
	if svc.autoSave {
		go svc.Save()
	}

	return nil
}

// プレイヤーの全有効権限を取得する（グループ権限と個別権限を含む）
func (svc *PermissionService) GetPlayerPermissions(playerID uuid.UUID) []string {
	svc.mutex.RLock()
	defer svc.mutex.RUnlock()

	player, exists := svc.players[playerID]
	if !exists {
		return []string{}
	}

	// 権限の重複を防ぐためのmap
	permissionSet := make(map[string]bool)

	// 個別権限を追加
	for _, perm := range player.Permissions {
		permissionSet[perm] = true
	}

	// グループ権限を追加
	for _, groupName := range player.Groups {
		if group, exists := svc.groups[groupName]; exists {
			for _, perm := range group.Permissions {
				permissionSet[perm] = true
			}
		}
	}

	// mapからsliceに変換
	var result []string
	for perm := range permissionSet {
		result = append(result, perm)
	}

	return result
}

// グループに権限を追加する
func (svc *PermissionService) AddPermissionToGroup(groupName, permission string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	group, exists := svc.groups[groupName]
	if !exists {
		return fmt.Errorf("group '%s' does not exist", groupName)
	}

	// 既に権限を持っているかチェック
	for _, perm := range group.Permissions {
		if perm == permission {
			return nil // 既に権限を持っている場合は何もしない
		}
	}

	// 権限を追加
	group.Permissions = append(group.Permissions, permission)

	// このグループに所属するプレイヤーのキャッシュを無効化
	for _, player := range svc.players {
		for _, gName := range player.Groups {
			if gName == groupName {
				if svc.cache != nil && svc.cache.IsEnabled() {
					svc.cache.InvalidatePlayer(player.PlayerID)
				}
				break
			}
		}
	}

	// オートセーブ
	if svc.autoSave {
		go svc.Save()
	}

	return nil
}

// グループから権限を削除する
func (svc *PermissionService) RemovePermissionFromGroup(groupName, permission string) error {
	svc.mutex.Lock()
	defer svc.mutex.Unlock()

	group, exists := svc.groups[groupName]
	if !exists {
		return fmt.Errorf("group '%s' does not exist", groupName)
	}

	// 権限を探して削除
	for i, perm := range group.Permissions {
		if perm == permission {
			group.Permissions = append(group.Permissions[:i], group.Permissions[i+1:]...)

			// このグループに所属するプレイヤーのキャッシュを無効化
			for _, player := range svc.players {
				for _, gName := range player.Groups {
					if gName == groupName {
						if svc.cache != nil && svc.cache.IsEnabled() {
							svc.cache.InvalidatePlayer(player.PlayerID)
						}
						break
					}
				}
			}

			// オートセーブ
			if svc.autoSave {
				go svc.Save()
			}

			return nil
		}
	}

	return fmt.Errorf("group '%s' does not have permission '%s'", groupName, permission)
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
