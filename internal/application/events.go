package application

import (
	"errors"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/shared"
)

// イベント系のユーティリティ関数を定義

// 参加時のイベントハンドラ
func (mgr *Manager) OnPlayerJoin(playerID uuid.UUID, playerName string) error {
	// 新しいプレイヤーを登録（既存なら無視）
	err := mgr.CreatePlayer(playerID, playerName)
	if err != nil && !errors.Is(err, shared.ErrPlayerAlreadyExists) {
		return err
	}

	// プレイヤー名更新（名前変更対応）
	err = mgr.UpdatePlayerName(playerID, playerName)
	if err != nil && !errors.Is(err, shared.ErrPlayerNotFound) {
		return err
	}

	return nil
}

// 退出時のイベントハンドラ
func (mgr *Manager) OnPlayerLeave(playerID uuid.UUID, playerName string) error {
	// キャッシュ削除処理
	mgr.cache.InvalidatePlayer(playerID)
	return nil
}
