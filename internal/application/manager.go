package application

import (
	"sync"
	"time"

	"github.com/google/uuid"
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

func NewManager(config shared.ManagerConfig) *Manager {
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
		return shared.NewGroupNotFoundError(groupName)
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
		return shared.ErrPlayerNotFound
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

	return shared.NewPlayerNotInGroupError(groupName)
}

// 指定されたパーミッションを持つ新しいパーミッショングループを作成する
// グループがすでに存在する場合はエラーを返す
func (mgr *Manager) CreateGroup(name string, permissions []string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// 権限リストの検証
	if invalidPerm, valid := mgr.checker.ValidatePermissions(permissions); !valid {
		return shared.NewInvalidPermissionError(invalidPerm)
	}

	if _, exists := mgr.groups[name]; exists {
		return shared.NewGroupAlreadyExistsError(name)
	}

	mgr.groups[name] = &shared.Group{
		Name:        name,
		Permissions: permissions,
	}

	if mgr.autoSave {
		go mgr.Save()
	}

	return nil
}

// パーミッショングループを削除する
func (mgr *Manager) DeleteGroup(name string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// デフォルトグループは削除できない
	if name == "default" || name == "admin" {
		return shared.NewSystemGroupProtectedError(name)
	}

	// グループが存在するかチェック
	if _, exists := mgr.groups[name]; !exists {
		return shared.NewGroupNotFoundError(name)
	}

	// グループを削除
	delete(mgr.groups, name)

	// このグループに所属していたプレイヤーからグループを削除
	for _, player := range mgr.players {
		for i, groupName := range player.Groups {
			if groupName == name {
				player.Groups = append(player.Groups[:i], player.Groups[i+1:]...)
				player.UpdatedAt = time.Now()

				// このプレイヤーのキャッシュを無効化
				if mgr.cache != nil && mgr.cache.IsEnabled() {
					mgr.cache.InvalidatePlayer(player.PlayerID)
				}
				break
			}
		}
	}

	// オートセーブ
	if mgr.autoSave {
		go mgr.Save()
	}

	return nil
}

// すべてのパーミッショングループのコピーを返す
func (mgr *Manager) GetAllGroups() map[string]*shared.Group {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	result := make(map[string]*shared.Group)
	for name, group := range mgr.groups {
		result[name] = &shared.Group{
			Name:        group.Name,
			Permissions: append([]string{}, group.Permissions...),
		}
	}
	return result
}

// パーミッショングループの権限を更新する
func (mgr *Manager) UpdateGroup(name string, permissions []string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// 権限リストの検証
	if invalidPerm, valid := mgr.checker.ValidatePermissions(permissions); !valid {
		return shared.NewInvalidPermissionError(invalidPerm)
	}

	// グループが存在するかチェック
	group, exists := mgr.groups[name]
	if !exists {
		return shared.NewGroupNotFoundError(name)
	}

	// 権限を更新
	group.Permissions = append([]string{}, permissions...)

	// このグループに所属するプレイヤーのキャッシュを無効化
	for _, player := range mgr.players {
		for _, groupName := range player.Groups {
			if groupName == name {
				if mgr.cache != nil && mgr.cache.IsEnabled() {
					mgr.cache.InvalidatePlayer(player.PlayerID)
				}
				break
			}
		}
	}

	// オートセーブ
	if mgr.autoSave {
		go mgr.Save()
	}

	return nil
}

// 特定のパーミッショングループを取得する
func (mgr *Manager) GetGroup(name string) *shared.Group {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	group, exists := mgr.groups[name]
	if !exists {
		return nil
	}

	return &shared.Group{
		Name:        group.Name,
		Permissions: append([]string{}, group.Permissions...),
	}
}

// プレイヤーデータのコピーを返す
func (mgr *Manager) GetPlayerData(playerID uuid.UUID) *shared.PlayerData {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	player, exists := mgr.players[playerID]
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

// プレイヤーデータを完全に削除する
func (mgr *Manager) RemovePlayer(playerID uuid.UUID) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// プレイヤーが存在するかチェック
	if _, exists := mgr.players[playerID]; !exists {
		return shared.NewPlayerNotFoundError(playerID)
	}

	// プレイヤーデータを削除
	delete(mgr.players, playerID)

	// キャッシュを無効化
	if mgr.cache != nil && mgr.cache.IsEnabled() {
		mgr.cache.InvalidatePlayer(playerID)
	}

	// オートセーブ
	if mgr.autoSave {
		go mgr.Save()
	}

	return nil
}

// プレイヤーが存在するかチェックする
func (mgr *Manager) PlayerExists(playerID uuid.UUID) bool {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	_, exists := mgr.players[playerID]
	return exists
}

// すべてのプレイヤーのリストを取得する
func (mgr *Manager) GetAllPlayers() map[uuid.UUID]*shared.PlayerData {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	result := make(map[uuid.UUID]*shared.PlayerData)
	for id, player := range mgr.players {
		result[id] = &shared.PlayerData{
			PlayerID:    player.PlayerID,
			PlayerName:  player.PlayerName,
			Groups:      append([]string{}, player.Groups...),
			Permissions: append([]string{}, player.Permissions...),
			UpdatedAt:   player.UpdatedAt,
		}
	}
	return result
}

// プレイヤーが所属するグループのリストを取得する
func (mgr *Manager) GetPlayerGroups(playerID uuid.UUID) []string {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	player, exists := mgr.players[playerID]
	if !exists {
		return []string{}
	}

	return append([]string{}, player.Groups...)
}

// プレイヤーを新規作成する
func (mgr *Manager) CreatePlayer(playerID uuid.UUID, playerName string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// プレイヤーが既に存在するかチェック
	if _, exists := mgr.players[playerID]; exists {
		return shared.NewPlayerAlreadyExistsError(playerID)
	}

	// プレイヤーデータを作成
	mgr.players[playerID] = &shared.PlayerData{
		PlayerID:    playerID,
		PlayerName:  playerName,
		Groups:      []string{},
		Permissions: []string{},
		UpdatedAt:   time.Now(),
	}

	// オートセーブ
	if mgr.autoSave {
		go mgr.Save()
	}

	return nil
}

// プレイヤーに個別権限を追加する
func (mgr *Manager) AddPlayerPermission(playerID uuid.UUID, permission string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// 権限フォーマットのバリデーション
	if !mgr.checker.ValidatePermission(permission) {
		return shared.NewInvalidPermissionError(permission)
	}

	// プレイヤーデータを取得または作成
	player, exists := mgr.players[playerID]
	if !exists {
		return shared.NewPlayerNotFoundError(playerID)
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
	if mgr.cache != nil && mgr.cache.IsEnabled() {
		mgr.cache.InvalidatePlayer(playerID)
	}

	// オートセーブ
	if mgr.autoSave {
		go mgr.Save()
	}

	return nil
}

// プレイヤーから個別権限を削除する
func (mgr *Manager) RemovePlayerPermission(playerID uuid.UUID, permission string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	player, exists := mgr.players[playerID]
	if !exists {
		return shared.NewPlayerNotFoundError(playerID)
	}

	// 権限を探して削除
	for i, perm := range player.Permissions {
		if perm == permission {
			player.Permissions = append(player.Permissions[:i], player.Permissions[i+1:]...)
			player.UpdatedAt = time.Now()

			// キャッシュを無効化
			if mgr.cache != nil && mgr.cache.IsEnabled() {
				mgr.cache.InvalidatePlayer(playerID)
			}

			// オートセーブ
			if mgr.autoSave {
				go mgr.Save()
			}

			return nil
		}
	}

	return shared.NewPlayerPermissionNotFoundError(permission)
}

// プレイヤーの権限を一括設定する
func (mgr *Manager) SetPlayerPermissions(playerID uuid.UUID, permissions []string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// 権限リストの検証
	if invalidPerm, valid := mgr.checker.ValidatePermissions(permissions); !valid {
		return shared.NewInvalidPermissionError(invalidPerm)
	}

	player, exists := mgr.players[playerID]
	if !exists {
		return shared.NewPlayerNotFoundError(playerID)
	}

	// 権限を置き換え
	player.Permissions = append([]string{}, permissions...)
	player.UpdatedAt = time.Now()

	// キャッシュを無効化
	if mgr.cache != nil && mgr.cache.IsEnabled() {
		mgr.cache.InvalidatePlayer(playerID)
	}

	// オートセーブ
	if mgr.autoSave {
		go mgr.Save()
	}

	return nil
}

// プレイヤーの全有効権限を取得する（グループ権限と個別権限を含む）
func (mgr *Manager) GetPlayerPermissions(playerID uuid.UUID) []string {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	player, exists := mgr.players[playerID]
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
		if group, exists := mgr.groups[groupName]; exists {
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
func (mgr *Manager) AddPermissionToGroup(groupName, permission string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// 権限フォーマットのバリデーション
	if !mgr.checker.ValidatePermission(permission) {
		return shared.NewInvalidPermissionError(permission)
	}

	group, exists := mgr.groups[groupName]
	if !exists {
		return shared.NewGroupNotFoundError(groupName)
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
	for _, player := range mgr.players {
		for _, gName := range player.Groups {
			if gName == groupName {
				if mgr.cache != nil && mgr.cache.IsEnabled() {
					mgr.cache.InvalidatePlayer(player.PlayerID)
				}
				break
			}
		}
	}

	// オートセーブ
	if mgr.autoSave {
		go mgr.Save()
	}

	return nil
}

// グループから権限を削除する
func (mgr *Manager) RemovePermissionFromGroup(groupName, permission string) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	group, exists := mgr.groups[groupName]
	if !exists {
		return shared.NewGroupNotFoundError(groupName)
	}

	// 権限を探して削除
	for i, perm := range group.Permissions {
		if perm == permission {
			group.Permissions = append(group.Permissions[:i], group.Permissions[i+1:]...)

			// このグループに所属するプレイヤーのキャッシュを無効化
			for _, player := range mgr.players {
				for _, gName := range player.Groups {
					if gName == groupName {
						if mgr.cache != nil && mgr.cache.IsEnabled() {
							mgr.cache.InvalidatePlayer(player.PlayerID)
						}
						break
					}
				}
			}

			// オートセーブ
			if mgr.autoSave {
				go mgr.Save()
			}

			return nil
		}
	}

	return shared.NewGroupPermissionNotFoundError(groupName, permission)
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
