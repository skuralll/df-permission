package shared

import (
	"time"

	"github.com/google/uuid"
)

// プレイヤーデータ
type PlayerData struct {
	PlayerID    uuid.UUID `json:"player_id"`   // プレイヤーのUUID
	PlayerName  string    `json:"player_name"` // プレイヤー名 (変更される可能性あり)
	Groups      []string  `json:"groups"`      // 所属権限グループのリスト
	Permissions []string  `json:"permissions"` // 個別に割り当てられた権限のリスト
	UpdatedAt   time.Time `json:"updated_at"`  // プレイヤーデータが最後に更新された日時
}

// 権限グループ
type Group struct {
	Name        string   `json:"name"`        // グループ名
	Permissions []string `json:"permissions"` // グループに割り当てられた権限のリスト
}