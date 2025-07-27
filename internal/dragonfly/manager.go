package dragonfly

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/skuralll/df-permission/permission"
)

type DragonflyManager struct {
	core permission.PermissionManager
	srv  *server.Server
}
