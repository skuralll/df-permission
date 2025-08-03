package events

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/shared"
	"github.com/skuralll/df-permission/permission"
)

type JoinHandler struct {
	permissionMgr permission.PermissionManager
	defaultGroup  string
}

func NewJoinHandler(mgr permission.PermissionManager) *JoinHandler {
	return &JoinHandler{
		permissionMgr: mgr,
		defaultGroup:  "default",
	}
}

func (h *JoinHandler) HandlePlayerJoin(playerID uuid.UUID, playerName string) error {
	if err := h.registerNewPlayer(playerID, playerName); err != nil {
		return err
	}

	if err := h.assignDefaultGroup(playerID, playerName); err != nil {
		return err
	}

	if err := h.updatePlayerName(playerID, playerName); err != nil {
		return fmt.Errorf("failed to update player name: %w", err)
	}

	return nil
}

func (h *JoinHandler) registerNewPlayer(playerID uuid.UUID, playerName string) error {
	if err := h.permissionMgr.CreatePlayer(playerID, playerName); err != nil {
		if !errors.Is(err, shared.ErrPlayerAlreadyExists) {
			return err
		}
	}
	return nil
}

func (h *JoinHandler) assignDefaultGroup(playerID uuid.UUID, playerName string) error {
	return h.permissionMgr.AddPlayerToGroup(playerID, playerName, h.defaultGroup)
}

func (h *JoinHandler) updatePlayerName(playerID uuid.UUID, playerName string) error {
	panic("unimplemented")
}
