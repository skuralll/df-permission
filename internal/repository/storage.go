package repository

import (
	"time"

	"github.com/google/uuid"
	permission "github.com/skuralll/df-permission"
)

type Storage interface {
	// ストレージバックエンドからPermissionDataを取得して返す
	Load() (*permission.PermissionData, error)

	// PermissionDataをストレージバックエンドに保存
	Save(data *permission.PermissionData) error

	// ストレージバックエンドに権限データが存在するか確認
	Exists() bool

	// 必要に応じてストレージバックエンドを安全に終了させる
	Close() error
}

// デフォルトのPermissionDataを返す
func NewDefaultPermissionData() *permission.PermissionData {
	return &permission.PermissionData{
		Groups:  make(map[string]*permission.Group),
		Players: make(map[uuid.UUID]*permission.PlayerData),
		Meta: &permission.Metadata{
			Version:   permission.PermissionDataVersion,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
