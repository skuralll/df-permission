package domain

import (
	"sync"

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
