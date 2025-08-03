package events

import (
	"errors"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/shared"
	"github.com/skuralll/df-permission/permission"
)

type JoinHandler struct {
	permissionMgr permission.PermissionManager
	defaultGroup  string
}

func NewJoinHandler(mgr permission.PermissionManager, defaultGroup string) *JoinHandler {
	return &JoinHandler{
		permissionMgr: mgr,
		defaultGroup:  defaultGroup,
	}
}

// プレイヤーが参加したときの処理
func (h *JoinHandler) HandlePlayerJoin(playerID uuid.UUID, playerName string) error {
	if err := h.registerNewPlayer(playerID, playerName); err != nil {
		return err
	}

	if err := h.assignDefaultGroup(playerID, playerName); err != nil {
		return err
	}

	if err := h.updatePlayerName(playerID, playerName); err != nil {
		return err
	}

	return nil
}

// 新しいプレイヤーを登録
func (h *JoinHandler) registerNewPlayer(playerID uuid.UUID, playerName string) error {
	if err := h.permissionMgr.CreatePlayer(playerID, playerName); err != nil {
		if !errors.Is(err, shared.ErrPlayerAlreadyExists) {
			return err
		}
	}
	return nil
}

// デフォルトグループにプレイヤーを追加
func (h *JoinHandler) assignDefaultGroup(playerID uuid.UUID, playerName string) error {
	if h.defaultGroup == "" {
		return nil // デフォルトグループが設定されていない場合は何もしない
	}
	// プレイヤーをデフォルトグループに追加
	return h.permissionMgr.AddPlayerToGroup(playerID, playerName, h.defaultGroup)
}

func (h *JoinHandler) updatePlayerName(playerID uuid.UUID, playerName string) error {
	if err := h.permissionMgr.UpdatePlayerName(playerID, playerName); err != nil {
		if !errors.Is(err, shared.ErrPlayerNotFound) {
			return err
		}
	}
	return nil
}
