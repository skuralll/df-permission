package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	// ストレージバックエンドからPermissionDataを取得して返す
	Load() (*PermissionData, error)

	// PermissionDataをストレージバックエンドに保存
	Save(data *PermissionData) error

	// ストレージバックエンドに権限データが存在するか確認
	Exists() bool

	// 必要に応じてストレージバックエンドを安全に終了させる
	Close() error
}

// デフォルトのPermissionDataを返す
func NewDefaultPermissionData() *PermissionData {
	return &PermissionData{
		Groups:  make(map[string]*Group),
		Players: make(map[uuid.UUID]*PlayerData),
		Meta: &Metadata{
			Version:   PermissionDataVersion,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// PermissionDataのバリデーションを行う
func ValidatePermissionData(data *PermissionData) error {
	if data == nil {
		return NewStorageError("validation", "permission data is nil")
	}

	if data.Groups == nil {
		return NewStorageError("validation", "groups map is nil")
	}

	if data.Players == nil {
		return NewStorageError("validation", "players map is nil")
	}

	if data.Meta == nil {
		return NewStorageError("validation", "metadata is nil")
	}

	// メタデータのバリデーション
	if data.Meta.Version == "" {
		return NewStorageError("validation", "metadata version is empty")
	}

	if data.Meta.CreatedAt.IsZero() {
		return NewStorageError("validation", "metadata created_at is zero")
	}

	if data.Meta.UpdatedAt.IsZero() {
		return NewStorageError("validation", "metadata updated_at is zero")
	}

	// グループのバリデーション
	for name, group := range data.Groups {
		if group == nil {
			return NewStorageError("validation",
				fmt.Sprintf("group '%s' is nil", name))
		}
		if group.Name != name {
			return NewStorageError("validation",
				fmt.Sprintf("group name mismatch: key='%s', name='%s'", name, group.Name))
		}
	}

	// プレイヤーデータのバリデーション
	for id, player := range data.Players {
		if player == nil {
			return NewStorageError("validation",
				fmt.Sprintf("player '%s' is nil", id))
		}
		if player.PlayerID != id {
			return NewStorageError("validation",
				fmt.Sprintf("player ID mismatch: key='%s', id='%s'", id, player.PlayerID))
		}
	}

	return nil
}
