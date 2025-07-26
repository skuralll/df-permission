package dfpermission

import "github.com/skuralll/df-permission/internal/shared"

// 公開API用にinternal/sharedから設定型を再エクスポート
type (
	StorageConfig = shared.StorageConfig
	CacheConfig   = shared.CacheConfig
	ManagerConfig = shared.ManagerConfig
)
