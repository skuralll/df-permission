package dfpermission

import "github.com/skuralll/df-permission/internal/shared"

// Re-export configuration types from internal/shared for public API
type (
	StorageConfig = shared.StorageConfig
	CacheConfig   = shared.CacheConfig
	ManagerConfig = shared.ManagerConfig
)
