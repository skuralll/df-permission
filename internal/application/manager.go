package application

import (
	"sync"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/domain"
	"github.com/skuralll/df-permission/internal/repository"
	"github.com/skuralll/df-permission/internal/shared"
)

type Manager struct {
	// config
	autoSave bool
	// internal state
	groups  map[string]*shared.Group
	players map[uuid.UUID]*shared.PlayerData
	storage repository.Storage
	cache   *domain.PermissionCache
	checker *domain.PermissionChecker
	mutex   sync.RWMutex
}
