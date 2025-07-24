package domain

import "testing"

func TestPermissionMatcher_Match(t *testing.T) {
	pm := NewPermissionMatcher()

	tests := []struct {
		pattern string
		target  string
		want    bool
	}{
		{"*", "anything", true},
		{"exact", "exact", true},
		{"exact", "different", false},
		{"world.*", "world.build", true},
		{"world.*", "world", true},
		{"world.*", "other.build", false},
		{"moderation.admin.*", "moderation.admin.kick", true},
		{"moderation.admin.*", "moderation.admin.ban.permanent", true},
		{"moderation.admin.*", "moderation.admin", true},
		{"moderation.admin.*", "moderation.user.kick", false},
		{"server.config.database.*", "server.config.database.read", true},
		{"server.config.database.*", "server.config.network.read", false},
		{"chat.color.rainbow", "chat.color.rainbow", true},
		{"chat.color.rainbow", "chat.color.red", false},
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
		{"", false},
		{"*", true},
		{"world.build", true},
		{"world.*", true},
		{".invalid", false},
		{"invalid.", false},
		{"invalid..double", false},
	}

	for _, tt := range tests {
		got := pm.ValidatePattern(tt.pattern)
		if got != tt.want {
			t.Errorf("ValidatePattern(%q) = %v, want %v", tt.pattern, got, tt.want)
		}
	}
}