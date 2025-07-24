package domain

import "sync"

// パーミッションチェックをスレッドセーフに扱う
type PermissionChecker struct {
	mutex   sync.RWMutex
	matcher *PermissionMatcher
}

func NewPermissionChecker() *PermissionChecker {
	return &PermissionChecker{
		mutex:   sync.RWMutex{},
		matcher: NewPermissionMatcher(),
	}
}

// プレイヤーの権限とグループの権限をチェックして、指定された権限があるかどうかを返す
func (pc *PermissionChecker) HasPermission(
	playerPermissions []string,
	playerGroups []string,
	groupsMap map[string][]string,
	permission string,
) bool {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	// プレイヤーの個別権限をチェック
	if pc.matcher.MatchAny(playerPermissions, permission) {
		return true
	}

	// プレイヤーのグループ権限をチェック
	for _, groupName := range playerGroups {
		if groupPermissions, exists := groupsMap[groupName]; exists {
			if pc.matcher.MatchAny(groupPermissions, permission) {
				return true
			}
		}
	}

	return false
}
