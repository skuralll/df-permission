package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/shared"
)

type Storage interface {
	// ストレージバックエンドからPermissionDataを取得して返す
	Load() (*shared.PermissionData, error)

	// PermissionDataをストレージバックエンドに保存
	Save(data *shared.PermissionData) error

	// ストレージバックエンドに権限データが存在するか確認
	Exists() bool

	// 必要に応じてストレージバックエンドを安全に終了させる
	Close() error
}

// デフォルトのPermissionDataを返す
func NewDefaultPermissionData() *shared.PermissionData {
	return &shared.PermissionData{
		Groups:  make(map[string]*shared.Group),
		Players: make(map[uuid.UUID]*shared.PlayerData),
		Meta: &shared.Metadata{
			Version:   shared.PermissionDataVersion,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// PermissionDataのバリデーションを行う
func ValidatePermissionData(data *shared.PermissionData) error {
	if data == nil {
		return shared.NewStorageError("validation", "permission data is nil")
	}

	if data.Groups == nil {
		return shared.NewStorageError("validation", "groups map is nil")
	}

	if data.Players == nil {
		return shared.NewStorageError("validation", "players map is nil")
	}

	if data.Meta == nil {
		return shared.NewStorageError("validation", "metadata is nil")
	}

	// メタデータのバリデーション
	if data.Meta.Version == "" {
		return shared.NewStorageError("validation", "metadata version is empty")
	}

	if data.Meta.CreatedAt.IsZero() {
		return shared.NewStorageError("validation", "metadata created_at is zero")
	}

	if data.Meta.UpdatedAt.IsZero() {
		return shared.NewStorageError("validation", "metadata updated_at is zero")
	}

	// グループのバリデーション
	for name, group := range data.Groups {
		if group == nil {
			return shared.NewStorageError("validation",
				fmt.Sprintf("group '%s' is nil", name))
		}
		if group.Name != name {
			return shared.NewStorageError("validation",
				fmt.Sprintf("group name mismatch: key='%s', name='%s'", name, group.Name))
		}
	}

	// プレイヤーデータのバリデーション
	for id, player := range data.Players {
		if player == nil {
			return shared.NewStorageError("validation",
				fmt.Sprintf("player '%s' is nil", id))
		}
		if player.PlayerID != id {
			return shared.NewStorageError("validation",
				fmt.Sprintf("player ID mismatch: key='%s', id='%s'", id, player.PlayerID))
		}
	}

	return nil
}
