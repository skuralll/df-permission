package dragonfly

import (
	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/dragonfly/events"
	"github.com/skuralll/df-permission/permission"
)

type JoinHandler interface {
	HandlePlayerJoin(playerID uuid.UUID, playerName string) error
}

type joinHandler struct {
	internal *events.JoinHandler
}

func NewJoinHandler(mgr permission.PermissionManager, defaultGroup string) JoinHandler {
	return &joinHandler{
		internal: events.NewJoinHandler(mgr, defaultGroup),
	}
}

func (h *joinHandler) HandlePlayerJoin(playerID uuid.UUID, playerName string) error {
	return h.internal.HandlePlayerJoin(playerID, playerName)
}