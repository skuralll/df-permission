package dfpermission

import (
	"errors"
	"strings"
)

// パーミッション管理システムの公開エラー定義
var (
	// 存在しないプレイヤー参照時のエラー
	ErrPlayerNotFound = errors.New("player not found")

	// 存在しないグループ参照時のエラー
	ErrGroupNotFound = errors.New("group not found")

	// 既存グループ作成時のエラー
	ErrGroupAlreadyExists = errors.New("group already exists")

	// システムグループ削除時のエラー
	ErrSystemGroupProtected = errors.New("system group cannot be deleted")

	// 既存プレイヤー作成時のエラー
	ErrPlayerAlreadyExists = errors.New("player already exists")

	// 存在しないパーミッション削除時のエラー
	ErrPermissionNotFound = errors.New("permission not found")
)

// 内部エラーを公開APIエラーへ変換
func convertError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// プレイヤー関連エラー変換
	if strings.Contains(errMsg, "player") {
		if strings.Contains(errMsg, "not found") {
			return ErrPlayerNotFound
		}
		if strings.Contains(errMsg, "already exists") {
			return ErrPlayerAlreadyExists
		}
	}

	// グループ関連エラー変換
	if strings.Contains(errMsg, "group") {
		if strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "does not exist") {
			return ErrGroupNotFound
		}
		if strings.Contains(errMsg, "already exists") {
			return ErrGroupAlreadyExists
		}
		if strings.Contains(errMsg, "cannot delete system group") || strings.Contains(errMsg, "system group") {
			return ErrSystemGroupProtected
		}
	}

	// パーミッション関連エラー変換
	if strings.Contains(errMsg, "permission") && strings.Contains(errMsg, "not") {
		return ErrPermissionNotFound
	}

	// 変換不要時は元のエラーを返却
	return err
}
