package domain

import "testing"

func TestPermissionChecker_HasPermission(t *testing.T) {
	pc := NewPermissionChecker()

	tests := []struct {
		name              string
		playerPermissions []string
		playerGroups      []string
		groupsMap         map[string][]string
		permission        string
		want              bool
	}{
		{
			// プレイヤーの直接権限による完全一致テスト
			name:              "direct permission match",
			playerPermissions: []string{"world.build"},
			playerGroups:      []string{},
			groupsMap:         map[string][]string{},
			permission:        "world.build",
			want:              true,
		},
		{
			// プレイヤーの直接権限によるワイルドカード一致テスト
			name:              "direct permission wildcard",
			playerPermissions: []string{"world.*"},
			playerGroups:      []string{},
			groupsMap:         map[string][]string{},
			permission:        "world.build",
			want:              true,
		},
		{
			// グループ権限による権限付与テスト
			name:              "group permission",
			playerPermissions: []string{},
			playerGroups:      []string{"admin"},
			groupsMap:         map[string][]string{"admin": {"*"}},
			permission:        "world.build",
			want:              true,
		},
		{
			// 権限が見つからない場合のテスト
			name:              "no permission",
			playerPermissions: []string{"chat.color"},
			playerGroups:      []string{"user"},
			groupsMap:         map[string][]string{"user": {"chat.*"}},
			permission:        "world.build",
			want:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pc.HasPermission(tt.playerPermissions, tt.playerGroups, tt.groupsMap, tt.permission)
			if got != tt.want {
				t.Errorf("HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionChecker_ValidatePermission(t *testing.T) {
	pc := NewPermissionChecker()

	tests := []struct {
		permission string
		want       bool
	}{
		{"world.build", true},  // 有効な権限文字列
		{"*", true},            // グローバルワイルドカード
		{"world.*", true},      // プレフィックスワイルドカード
		{"", false},            // 空文字列は無効
		{".invalid", false},    // ドットで始まる文字列は無効
	}

	for _, tt := range tests {
		got := pc.ValidatePermission(tt.permission)
		if got != tt.want {
			t.Errorf("ValidatePermission(%q) = %v, want %v", tt.permission, got, tt.want)
		}
	}
}

func TestPermissionChecker_ValidatePermissions(t *testing.T) {
	pc := NewPermissionChecker()

	// 全て有効な場合
	invalidPerm, valid := pc.ValidatePermissions([]string{"world.build", "chat.*", "*"})
	if !valid || invalidPerm != "" {
		t.Errorf("ValidatePermissions(valid) = %q, %v, want empty, true", invalidPerm, valid)
	}

	// 無効な権限がある場合
	invalidPerm, valid = pc.ValidatePermissions([]string{"world.build", ".invalid", "chat.*"})
	if valid || invalidPerm != ".invalid" {
		t.Errorf("ValidatePermissions(invalid) = %q, %v, want .invalid, false", invalidPerm, valid)
	}
}