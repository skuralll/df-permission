package domain

import "github.com/skuralll/df-permission/internal/shared"

type PermissionRepository interface {
	SavePermissions(data *shared.PermissionData) error
	LoadPermissions() (*shared.PermissionData, error)
}
