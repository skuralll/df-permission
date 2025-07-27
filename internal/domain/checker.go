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

// プレイヤーの権限とグループの権限をチェックして、指定された権限があるかどうかとその詳細を返す
func (pc *PermissionChecker) HasPermissionDetailed(
	playerPermissions []string,
	playerGroups []string,
	groupsMap map[string][]string,
	permission string,
) (granted bool, source string, reason string) {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	// プレイヤーの個別権限をチェック
	if matched, pattern, matchType := pc.matcher.MatchAnyDetailed(playerPermissions, permission); matched {
		return true, "direct", "directly assigned permission (pattern: " + pattern + ", type: " + string(matchType) + ")"
	}

	// プレイヤーのグループ権限をチェック
	for _, groupName := range playerGroups {
		if groupPermissions, exists := groupsMap[groupName]; exists {
			if matched, pattern, matchType := pc.matcher.MatchAnyDetailed(groupPermissions, permission); matched {
				return true, groupName, "inherited from group " + groupName + " (pattern: " + pattern + ", type: " + string(matchType) + ")"
			}
		}
	}

	return false, "", "permission not found in direct permissions or any group"
}

// 権限のパターンを検証する
func (pc *PermissionChecker) ValidatePermission(permission string) bool {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	return pc.matcher.ValidatePattern(permission)
}

// 権限リストを検証する
func (pc *PermissionChecker) ValidatePermissions(permissions []string) (string, bool) {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	for _, permission := range permissions {
		if !pc.matcher.ValidatePattern(permission) {
			return permission, false
		}
	}
	return "", true
}
