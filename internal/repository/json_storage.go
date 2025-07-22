package repository

import (
	"sync"

	permission "github.com/skuralll/df-permission"
)

type JSONStorage struct {
	config permission.StorageConfig
	mutex  sync.RWMutex
}

var _ Storage = (*JSONStorage)(nil)

func NewJSONStorage(config permission.StorageConfig) *JSONStorage {
	return &JSONStorage{
		config: config,
		mutex:  sync.RWMutex{},
	}
}

// Close implements Storage.
func (j *JSONStorage) Close() error {
	panic("unimplemented")
}

// Exists implements Storage.
func (j *JSONStorage) Exists() bool {
	panic("unimplemented")
}

// Load implements Storage.
func (j *JSONStorage) Load() (*permission.PermissionData, error) {
	panic("unimplemented")
}

// Save implements Storage.
func (j *JSONStorage) Save(data *permission.PermissionData) error {
	panic("unimplemented")
}
