package domain

import (
	dfpermission "github.com/skuralll/df-permission"
	"github.com/skuralll/df-permission/internal/repository"
	"github.com/skuralll/df-permission/internal/shared"
)

type PermissionRepository struct {
	storage repository.Storage
}

func NewPermissionRepository(config dfpermission.StorageConfig) *PermissionRepository {
	return &PermissionRepository{
		storage: repository.NewJSONStorage(config),
	}
}

func (p *PermissionRepository) Save(data *shared.PermissionData) error {
	return p.storage.Save(data)
}

func (p *PermissionRepository) Load() (*shared.PermissionData, error) {
	return p.storage.Load()
}
