package dfpermission

import "github.com/skuralll/df-permission/internal/shared"

// PermissionServiceの設定
type ManagerConfig struct {
	AutoSave bool
	Storage  shared.StorageConfig
	Cache    shared.CacheConfig
}
