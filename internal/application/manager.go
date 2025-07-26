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
	// mgr.initializeDefaultGroups()

	// 既存データをロード
	// mgr.loadData()

	return mgr
}
