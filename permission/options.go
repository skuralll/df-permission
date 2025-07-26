package permission

import (
	"time"

	"github.com/skuralll/df-permission/internal/shared"
)

// Option は PermissionManager の設定をカスタマイズするための関数型
type Option func(*shared.ManagerConfig)

// ストレージファイルのパスを設定するオプション
func WithStorage(path string) Option {
	return func(config *shared.ManagerConfig) {
		config.Storage.Path = path
	}
}

// オートセーブの有効/無効を設定するオプション
func WithAutoSave(enabled bool) Option {
	return func(config *shared.ManagerConfig) {
		config.AutoSave = enabled
	}
}

// キャッシュの有効/無効を設定するオプション
func WithCacheEnabled(enabled bool) Option {
	return func(config *shared.ManagerConfig) {
		config.Cache.Enabled = enabled
	}
}

// キャッシュのTTL（有効期限）を設定するオプション
// キャッシュを有効にし、指定されたTTLを設定する
func WithCache(ttl time.Duration) Option {
	return func(config *shared.ManagerConfig) {
		config.Cache.Enabled = true
		config.Cache.TTL = ttl
	}
}

// キャッシュのクリーンアップ間隔を設定するオプション
func WithCacheCleanup(interval time.Duration) Option {
	return func(config *shared.ManagerConfig) {
		config.Cache.CleanupInterval = interval
	}
}

// デフォルト設定を返すヘルパー関数
func defaultConfig() shared.ManagerConfig {
	return shared.ManagerConfig{
		AutoSave: true,
		Storage: shared.StorageConfig{
			Path: "./permissions.json",
		},
		Cache: shared.CacheConfig{
			Enabled:         true,
			TTL:             30 * time.Second,
			CleanupInterval: time.Minute,
		},
	}
}

// オプションを適用してManagerConfigを構築する内部関数
func buildConfig(opts ...Option) shared.ManagerConfig {
	config := defaultConfig()
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
