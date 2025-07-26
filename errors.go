package dfpermission

import (
	"errors"
	"strings"
)

// Public error definitions for the permission management system
var (
	// ErrPlayerNotFound is returned when a requested player does not exist
	ErrPlayerNotFound = errors.New("player not found")

	// ErrGroupNotFound is returned when a requested group does not exist
	ErrGroupNotFound = errors.New("group not found")

	// ErrGroupAlreadyExists is returned when trying to create a group that already exists
	ErrGroupAlreadyExists = errors.New("group already exists")

	// ErrSystemGroupProtected is returned when trying to delete a system group
	ErrSystemGroupProtected = errors.New("system group cannot be deleted")

	// ErrPlayerAlreadyExists is returned when trying to create a player that already exists
	ErrPlayerAlreadyExists = errors.New("player already exists")

	// ErrPermissionNotFound is returned when trying to remove a permission that doesn't exist
	ErrPermissionNotFound = errors.New("permission not found")
)

// convertError converts internal errors to public API errors.
// This provides a stable error interface and hides internal implementation details.
func convertError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Convert player-related errors
	if strings.Contains(errMsg, "player") {
		if strings.Contains(errMsg, "not found") {
			return ErrPlayerNotFound
		}
		if strings.Contains(errMsg, "already exists") {
			return ErrPlayerAlreadyExists
		}
	}

	// Convert group-related errors
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

	// Convert permission-related errors
	if strings.Contains(errMsg, "permission") && strings.Contains(errMsg, "not") {
		return ErrPermissionNotFound
	}

	// Return original error if no conversion is needed
	return err
}
