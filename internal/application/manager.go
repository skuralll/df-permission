package application

import (
	"sync"

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
