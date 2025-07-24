package domain

import "strings"

// Permission文字列のパターンマッチングを扱う
type PermissionMatcher struct{}

// パターンマッチのタイプ
type MatchType string

const (
	GlobalWildcard MatchType = "global_wildcard"
	ExactMatch     MatchType = "exact_match"
	PrefixWildcard MatchType = "prefix_wildcard"
	UnknownMatch   MatchType = "unknown"
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

// パターンが有効なPermission文字列かどうかを確認
func (pm *PermissionMatcher) ValidatePattern(pattern string) bool {
	// パターンが空文字列の場合は無効
	if pattern == "" {
		return false
	}

	// グローバルワイルドカードは常に有効
	if pattern == "*" {
		return true
	}

	// プレフィックスワイルドカード (e.g. "world.*") の場合
	if strings.HasSuffix(pattern, ".*") {
		prefix := pattern[:len(pattern)-2]
		// プレフィックスが空文字列の場合は無効
		if prefix == "" {
			return false
		}
		return pm.validatePermissionString(prefix)
	}

	// ワイルドカードなしのパターンの場合、文字列全体が有効なPermission文字列であるかを確認
	return pm.validatePermissionString(pattern)
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

// 計算量の少ないパターン順に並べ替える
// 完全一致 -> プレフィックス -> ワイルドカード
func (pm *PermissionMatcher) OptimizePatterns(patterns []string) []string {
	if len(patterns) <= 1 {
		return patterns
	}

	optimized := make([]string, 0, len(patterns))
	wildcards := make([]string, 0)
	prefixes := make([]string, 0)
	exacts := make([]string, 0)

	// パターン別に分類
	for _, pattern := range patterns {
		if pattern == "*" {
			wildcards = append(wildcards, pattern)
		} else if strings.HasSuffix(pattern, ".*") {
			prefixes = append(prefixes, pattern)
		} else {
			exacts = append(exacts, pattern)
		}
	}

	// 計算量の少ない順に結合
	optimized = append(optimized, exacts...)
	optimized = append(optimized, prefixes...)
	optimized = append(optimized, wildcards...)

	return optimized
}
