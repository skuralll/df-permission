package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

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
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// バリデーション
	if err := ValidatePermissionData(data); err != nil {
		return NewStorageError("validation", err.Error())
	}

	// ディレクトリの存在確認, なければ作成
	dir := filepath.Dir(j.config.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return NewStorageError("mkdir", err.Error())
	}

	// メタデータの更新
	data.Meta.UpdatedAt = time.Now()
	if data.Meta.CreatedAt.IsZero() {
		data.Meta.CreatedAt = time.Now()
	}

	// JSONへ変換
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return NewStorageError("marshal", err.Error())
	}

	// 一時ファイルを使用してAtomicに保存
	tempPath := j.config.Path + ".tmp"
	if err := os.WriteFile(tempPath, jsonData, 0644); err != nil {
		return NewStorageError("write", err.Error())
	}
	if err := os.Rename(tempPath, j.config.Path); err != nil {
		os.Remove(tempPath)
		return NewStorageError("rename", err.Error())
	}

	// ディスクに強制的に書き込む
	if err := j.syncFile(); err != nil {
		return NewStorageError("sync", err.Error())
	}

	return nil
}

// Close implements Storage.
func (j *JSONStorage) Close() error {
	panic("unimplemented")
}

// ディスクに強制的に書き込む
func (j *JSONStorage) syncFile() error {
	file, err := os.OpenFile(j.config.Path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	return file.Sync()
}
