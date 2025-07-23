package repository

import (
	"encoding/json"
	"os"
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

func (j *JSONStorage) Load() (*permission.PermissionData, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	// ストレージファイルが存在しない場合はデフォルトのPermissionDataを返す
	if !j.Exists() {
		return NewDefaultPermissionData(), nil
	}

	// ファイルからデータを読み込む
	data, err := os.ReadFile(j.config.Path)
	if err != nil {
		return nil, NewStorageError("load", err.Error())
	}

	// データをPermissionDataに変換
	var permData permission.PermissionData
	if err := json.Unmarshal(data, &permData); err != nil {
		return nil, NewStorageError("parse", err.Error())
	}

	// バリデーション
	if err := ValidatePermissionData(&permData); err != nil {
		return nil, NewStorageError("validation", err.Error())
	}

	return nil, nil
}

func (j *JSONStorage) Exists() bool {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	_, err := os.Stat(j.config.Path)
	return err == nil
}

// Save implements Storage.
func (j *JSONStorage) Save(data *permission.PermissionData) error {
	panic("unimplemented")
}

// Close implements Storage.
func (j *JSONStorage) Close() error {
	panic("unimplemented")
}
