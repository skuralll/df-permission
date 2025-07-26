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
