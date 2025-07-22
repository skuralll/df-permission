package repository

import (
	"sync"

	permission "github.com/skuralll/df-permission"
)

type JSONStorage struct {
	config permission.StorageConfig
	mutex  sync.RWMutex
}

var _ JSONStorage = Storage{}
