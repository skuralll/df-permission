package events

import (
	"fmt"

	"github.com/google/uuid"
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
		return fmt.Errorf("failed to register player: %w", err)
	}

	if err := h.assignDefaultGroup(playerID); err != nil {
		return fmt.Errorf("failed to assign default group: %w", err)
	}

	if h.isFirstTimeJoin(playerID) {
		h.handleFirstTimeJoin(playerID)
	}

	if err := h.updatePlayerName(playerID, playerName); err != nil {
		return fmt.Errorf("failed to update player name: %w", err)
	}

	return nil
}

func (h *JoinHandler) registerNewPlayer(playerID uuid.UUID, playerName string) error {
	panic("unimplemented")
}

func (h *JoinHandler) assignDefaultGroup(playerID uuid.UUID) error {
	panic("unimplemented")
}

func (h *JoinHandler) isFirstTimeJoin(playerID uuid.UUID) bool {
	panic("unimplemented")
}

func (h *JoinHandler) handleFirstTimeJoin(playerID uuid.UUID) {
	panic("unimplemented")
}

func (h *JoinHandler) updatePlayerName(playerID uuid.UUID, playerName string) error {
	panic("unimplemented")
}
