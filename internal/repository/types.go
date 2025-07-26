package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/skuralll/df-permission/internal/shared"
)

// 権限データバージョン
const PermissionDataVersion = "1.0"

// 権限設定メタデータ
type Metadata struct {
	Version   string    `json:"version"`    // バージョン情報
	CreatedAt time.Time `json:"created_at"` // 権限データが作成された日時
	UpdatedAt time.Time `json:"updated_at"` // 権限データが最後に更新された日時
}

// 権限データ
type PermissionData struct {
	Groups  map[string]*shared.Group         `json:"groups"`  // 権限グループのマップ
	Players map[uuid.UUID]*shared.PlayerData `json:"players"` // プレイヤーデータのマップ
	Meta    *Metadata                        `json:"meta"`    // 権限設定メタデータ
}
