package application

import (
	"errors"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/shared"
)

// 参加時のイベントハンドラ
func (mgr *Manager) OnPlayerJoin(playerID uuid.UUID, playerName string) error {
	// 新しいプレイヤーを登録（既存なら無視）
	err := mgr.CreatePlayer(playerID, playerName)
	if err != nil && !errors.Is(err, shared.ErrPlayerAlreadyExists) {
		return err
	}
	return nil
}

// 退出時のイベントハンドラ
func (mgr *Manager) OnPlayerLeave(playerID uuid.UUID) error {
	// キャッシュ削除処理
	mgr.cache.InvalidatePlayer(playerID)
	return nil
}
