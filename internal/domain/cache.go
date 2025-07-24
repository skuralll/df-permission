package domain

import (
	"sync"
	"time"
)

type PermissionCache struct {
	data            map[string]*cacheEntry
	ttl             time.Duration
	cleanupInterval time.Duration
	mutex           sync.RWMutex
	stopCleanup     chan struct{}
	enabled         bool
}

type cacheEntry struct {
	result    bool
	expiresAt time.Time
	createdAt time.Time
}

type CacheConfig struct {
	TTL             time.Duration
	CleanupInterval time.Duration
	Enabled         bool
}
