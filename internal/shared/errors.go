package shared

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ストレージ関連エラー
var (
	ErrStorage = errors.New("storage error")
)

// プレイヤー関連エラー
var (
	ErrPlayerNotFound           = errors.New("player not found")
	ErrPlayerAlreadyExists      = errors.New("player already exists")
	ErrPlayerPermissionNotFound = errors.New("player does not have permission")
	ErrPlayerNotInGroup         = errors.New("player is not in group")
)

// グループ関連エラー
var (
	ErrGroupNotFound           = errors.New("group not found")
	ErrGroupAlreadyExists      = errors.New("group already exists")
	ErrSystemGroupProtected    = errors.New("system group cannot be deleted")
	ErrGroupPermissionNotFound = errors.New("group does not have permission")
)

// 権限関連エラー
var (
	ErrInvalidPermission = errors.New("invalid permission format")
)

// ストレージエラー生成関数
func NewStorageError(operation, message string) error {
	return fmt.Errorf("%w: %s failed: %s", ErrStorage, operation, message)
}

// プレイヤーエラー生成関数
func NewPlayerNotFoundError(playerID uuid.UUID) error {
	return fmt.Errorf("%w: player with ID %s", ErrPlayerNotFound, playerID.String())
}

func NewPlayerAlreadyExistsError(playerID uuid.UUID) error {
	return fmt.Errorf("%w: player with ID %s", ErrPlayerAlreadyExists, playerID.String())
}

func NewPlayerPermissionNotFoundError(permission string) error {
	return fmt.Errorf("%w: '%s'", ErrPlayerPermissionNotFound, permission)
}

func NewPlayerNotInGroupError(groupName string) error {
	return fmt.Errorf("%w: '%s'", ErrPlayerNotInGroup, groupName)
}

// グループエラー生成関数
func NewGroupNotFoundError(groupName string) error {
	return fmt.Errorf("%w: '%s'", ErrGroupNotFound, groupName)
}

func NewGroupAlreadyExistsError(groupName string) error {
	return fmt.Errorf("%w: '%s'", ErrGroupAlreadyExists, groupName)
}

func NewSystemGroupProtectedError(groupName string) error {
	return fmt.Errorf("%w: '%s'", ErrSystemGroupProtected, groupName)
}

func NewGroupPermissionNotFoundError(groupName, permission string) error {
	return fmt.Errorf("%w: '%s' does not have permission '%s'", ErrGroupPermissionNotFound, groupName, permission)
}

// 権限エラー生成関数
func NewInvalidPermissionError(permission string) error {
	return fmt.Errorf("%w: '%s'", ErrInvalidPermission, permission)
}
