package domain

import "strings"

// Permission文字列のパターンマッチングを扱う
type PermissionMatcher struct{}

// パターンマッチのタイプ
type MatchType int

const (
	GlobalWildcard MatchType = iota
	ExactMatch
	PrefixWildcard
	UnknownMatch
)

// PermissionMatcherの新しいインスタンスを作成
func NewPermissionMatcher() *PermissionMatcher {
	return &PermissionMatcher{}
}

// パターンとターゲットのPermission文字列を比較し、マッチするかどうかを返す
func (pm *PermissionMatcher) Match(pattern, target string) bool {
	// グルーバルワイルドカード (すべてを許可)
	if pattern == "*" {
		return true
	}

	// 完全一致
	if pattern == target {
		return true
	}

	// プレフィックスワイルドカード (e.g. "world.build" は "world.*" にマッチ)
	if strings.HasSuffix(pattern, ".*") {
		prefix := pattern[:len(pattern)-2]
		return strings.HasPrefix(target, prefix+".") || target == prefix
	}

	return false
}

// 複数のパターンに対して、いずれかにマッチするかを確認
func (pm *PermissionMatcher) MatchAny(patterns []string, target string) bool {
	for _, pattern := range patterns {
		if pm.Match(pattern, target) {
			return true
		}
	}
	return false
}

// 複数のパターンに対して、いずれかにマッチする場合、マッチしたパターンとそのタイプを返す
func (pm *PermissionMatcher) MatchAnyDetailed(patterns []string, target string) (matched bool, matchedPattern string, matchType MatchType) {
	for _, pattern := range patterns {
		if pm.Match(pattern, target) {
			matchType := pm.getMatchType(pattern, target)
			return true, pattern, matchType
		}
	}
	return false, "", UnknownMatch
}

// マッチの種類を返す
func (pm *PermissionMatcher) getMatchType(pattern, target string) MatchType {
	if pattern == "*" {
		return GlobalWildcard
	}
	if pattern == target {
		return ExactMatch
	}
	if strings.HasSuffix(pattern, ".*") {
		return PrefixWildcard
	}
	return UnknownMatch
}

// 権限文字列をバリデーション
// - 英数字、ドット、アンダースコア、ハイフンのみを含むことができる
// - ドットで始まったり終わったりしてはいけない。
// - 連続したドットも許可されない。
func (pm *PermissionMatcher) validatePermissionString(permission string) bool {
	if permission == "" {
		return false
	}

	// 文字をチェック
	for _, r := range permission {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '.' || r == '_' || r == '-') {
			return false
		}
	}

	// ドットで始まったり終わったりしてはいけない
	if strings.HasPrefix(permission, ".") || strings.HasSuffix(permission, ".") {
		return false
	}

	// 連続したドットがないかチェック
	if strings.Contains(permission, "..") {
		return false
	}

	return true
}
