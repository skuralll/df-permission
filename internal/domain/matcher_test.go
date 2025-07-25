package domain

import "testing"

func TestPermissionMatcher_Match(t *testing.T) {
	pm := NewPermissionMatcher()

	tests := []struct {
		pattern string
		target  string
		want    bool
	}{
		{"*", "anything", true},                                           // グローバルワイル드カード（全権限を許可）
		{"exact", "exact", true},                                          // 完全一致
		{"exact", "different", false},                                     // 完全一致でマッチしない
		{"world.*", "world.build", true},                                  // プレフィックスワイルドカードでマッチ
		{"world.*", "world", true},                                        // プレフィックスワイルドカードで基底マッチ
		{"world.*", "other.build", false},                                 // プレフィックスワイルドカードでマッチしない
		{"moderation.admin.*", "moderation.admin.kick", true},             // 長いプレフィックスでマッチ
		{"moderation.admin.*", "moderation.admin.ban.permanent", true},    // ネストした長いプレフィックスでマッチ
		{"moderation.admin.*", "moderation.admin", true},                  // 長いプレフィックスで基底マッチ
		{"moderation.admin.*", "moderation.user.kick", false},             // 長いプレフィックスでマッチしない
		{"server.config.database.*", "server.config.database.read", true}, // データベース設定権限でマッチ
		{"server.config.database.*", "server.config.network.read", false}, // 異なる設定権限でマッチしない
		{"chat.color.rainbow", "chat.color.rainbow", true},                // 長い完全一致
		{"chat.color.rainbow", "chat.color.red", false},                   // 長い完全一致でマッチしない
	}

	for _, tt := range tests {
		got := pm.Match(tt.pattern, tt.target)
		if got != tt.want {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.target, got, tt.want)
		}
	}
}

func TestPermissionMatcher_ValidatePattern(t *testing.T) {
	pm := NewPermissionMatcher()

	tests := []struct {
		pattern string
		want    bool
	}{
		{"", false},                // 空文字列は無効
		{"*", true},                // グローバルワイルドカードは有効
		{"world.build", true},      // 有効な権限文字列
		{"world.*", true},          // 有効なプレフィックスワイルドカード
		{".invalid", false},        // ドットで始まる文字列は無効
		{"invalid.", false},        // ドットで終わる文字列は無効
		{"invalid..double", false}, // 連続したドットは無効
	}

	for _, tt := range tests {
		got := pm.ValidatePattern(tt.pattern)
		if got != tt.want {
			t.Errorf("ValidatePattern(%q) = %v, want %v", tt.pattern, got, tt.want)
		}
	}
}
