package dfpermission

import (
	"time"

	"github.com/google/uuid"
)

// 権限設定メタデータ
type Metadata struct {
	Version   string    `json:"version"`    // バージョン情報
	CreatedAt time.Time `json:"created_at"` // 権限データが作成された日時
	UpdatedAt time.Time `json:"updated_at"` // 権限データが最後に更新された日時
}

// 権限データ
type PermissionData struct {
	Groups  map[string]*Group         `json:"groups"`  // 権限グループのマップ
	Players map[uuid.UUID]*PlayerData `json:"players"` // プレイヤーデータのマップ
	Meta    *Metadata                 `json:"meta"`    // 権限設定メタデータ
}

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
