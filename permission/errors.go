package permission

import (
	"github.com/skuralll/df-permission/internal/shared"
)

// パーミッション管理システムの公開エラー定義
// 全て shared パッケージから再エクスポート
var (
	// ストレージ関連エラー
	ErrStorage = shared.ErrStorage

	// プレイヤー関連エラー
	ErrPlayerNotFound           = shared.ErrPlayerNotFound
	ErrPlayerAlreadyExists      = shared.ErrPlayerAlreadyExists
	ErrPlayerPermissionNotFound = shared.ErrPlayerPermissionNotFound
	ErrPlayerNotInGroup         = shared.ErrPlayerNotInGroup

	// グループ関連エラー
	ErrGroupNotFound           = shared.ErrGroupNotFound
	ErrGroupAlreadyExists      = shared.ErrGroupAlreadyExists
	ErrSystemGroupProtected    = shared.ErrSystemGroupProtected
	ErrGroupPermissionNotFound = shared.ErrGroupPermissionNotFound

	// 権限関連エラー
	ErrInvalidPermission = shared.ErrInvalidPermission

	// 下位互換性のため (使用非推奨)
	ErrPermissionNotFound = shared.ErrPlayerPermissionNotFound
)
