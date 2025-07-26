package dfpermission

import (
	"time"
)

// ストレージファイルの設定
type StorageConfig struct {
	// 保存するファイルのパス
	Path string
}

// キャッシュの設定
type CacheConfig struct {
	TTL             time.Duration
	CleanupInterval time.Duration
	Enabled         bool
}

// PermissionServiceの設定
type ManagerConfig struct {
	AutoSave bool
	Storage  StorageConfig
	Cache    CacheConfig
}
