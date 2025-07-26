# df-permission 仕様

## 概要

Minecraft Bedrock Edition Dragonflyサーバー向けの包括的な権限管理ライブラリです。

## 実装状況

### ✅ 実装済み機能

#### コア機能
- ✅ プレイヤーの権限チェック（キャッシュ対応）
- ✅ グループによる権限管理
- ✅ ワイルドカード権限（`*` と `prefix.*` パターンのサポート）
- ✅ ファイルベースストレージ（JSON）
- ✅ スレッドセーフな操作

#### 権限管理
- ✅ **グループ管理**: 作成、削除、更新、権限追加/削除
- ✅ **プレイヤー管理**: 作成、削除、グループ参加/離脱、個人権限設定
- ✅ **個人権限**: プレイヤーに直接権限を付与/削除
- ✅ **一括操作**: 権限の一括設定・取得

#### パフォーマンス機能
- ✅ **権限チェックキャッシュ**: TTL対応、自動クリーンアップ
- ✅ **オートセーブ**: 変更時の自動保存
- ✅ **データ整合性**: キャッシュ無効化による一貫性保証

#### 内部アーキテクチャ
- ✅ **Domain層**: 権限チェック、パターンマッチング、キャッシュ
- ✅ **Repository層**: JSONストレージとデータアクセス
- ✅ **Application層**: 統合管理機能（Manager）
- ✅ **Shared**: 共通型定義と設定

### 🚧 未実装機能

#### 公開API
- ❌ **permission.go**: 公開APIエントリーポイント
- ❌ **options.go**: オプションパターン実装
- ❌ **errors.go**: 公開エラー定義
- ❌ **QuickStart()**: 簡単セットアップ関数

## プロジェクト構造

### 現在の実装済み構造
```
df-permission/
├── types.go                    # ✅ 公開データ型（ManagerConfig）
├── docs/
│   └── SPECIFICATION.md        # ✅ 仕様書
│
└── internal/                   # ✅ 内部実装
    ├── domain/                 # ✅ ドメイン層
    │   ├── checker.go          # ✅ 権限チェックロジック
    │   ├── matcher.go          # ✅ パターンマッチング
    │   ├── cache.go            # ✅ 権限チェックキャッシュ
    │   └── storage.go          # ✅ ドメイン層ストレージ（リポジトリブリッジ）
    ├── repository/             # ✅ データアクセス層
    │   ├── storage.go          # ✅ Storage interface
    │   ├── json_storage.go     # ✅ JSON実装
    │   └── errors.go           # ✅ ストレージエラー
    ├── application/            # ✅ アプリケーション層
    │   └── manager.go          # ✅ 統合管理機能
    └── shared/                 # ✅ 共有定義
        ├── types.go            # ✅ 共有データ型
        └── config.go           # ✅ 設定構造体
```

### 予定されている公開API構造
```
df-permission/
├── permission.go      # ❌ 公開API（未実装）
├── options.go         # ❌ オプションパターン（未実装）
├── errors.go          # ❌ 公開エラー定義（未実装）
└── types.go           # ✅ 公開データ型（実装済み）
```

## 使用例

### 現在の実装での使用方法
```go
package main

import (
    "github.com/google/uuid"
    dfpermission "github.com/skuralll/df-permission"
    "github.com/skuralll/df-permission/internal/application"
    "github.com/skuralll/df-permission/internal/shared"
    "time"
)

func main() {
    // 1. 設定を作成
    config := dfpermission.ManagerConfig{
        AutoSave: true,
        Storage: shared.StorageConfig{
            Path: "./permissions.json",
        },
        Cache: shared.CacheConfig{
            Enabled:         true,
            TTL:             30 * time.Second,
            CleanupInterval: time.Minute,
        },
    }
    
    // 2. マネージャーを作成
    mgr := application.NewManager(config)
    
    // 3. グループを作成
    mgr.CreateGroup("vip", []string{"chat.color", "world.build.fast"})
    mgr.CreateGroup("moderator", []string{"moderation.*"})
    
    // 4. プレイヤーをグループに追加
    playerID := uuid.New()
    mgr.AddPlayerToGroup(playerID, "Steve", "vip")
    
    // 5. 権限をチェック
    if mgr.HasPermission(playerID, "chat.color") {
        // プレイヤーはカラーチャットを使用可能
    }
    
    if mgr.HasPermission(playerID, "moderation.kick") {
        // プレイヤーは他のプレイヤーをキック可能（VIPでは false）
    }
    
    // 6. データを手動保存
    mgr.Save()
}
```

### 予定されている公開API
```go
package main

import (
    "github.com/google/uuid"
    "github.com/skuralll/df-permission"
)

func main() {
    // 1. デフォルト設定でマネージャーを作成（未実装）
    mgr := permission.QuickStart()
    
    // 以下、現在の実装と同様のAPI...
}
```

---

## 詳細仕様

### 権限システム
- `*` : 全ての権限を付与（グローバルワイルドカード）
- `prefix.*` : `prefix.`で始まる全ての権限を付与（プレフィックスワイルドカード）
- `prefix.specific` : 特定の権限のみを付与

### 実装済みAPI（internal/application/manager.go）

#### 権限チェック
- `HasPermission(playerID uuid.UUID, permission string) bool`
  - プレイヤーが特定の権限を持っているかチェック
  - キャッシュされた結果を返す

#### プレイヤー管理
- `CreatePlayer(playerID uuid.UUID, playerName string) error`
- `RemovePlayer(playerID uuid.UUID) error`
- `PlayerExists(playerID uuid.UUID) bool`
- `GetPlayerData(playerID uuid.UUID) *shared.PlayerData`
- `GetAllPlayers() map[uuid.UUID]*shared.PlayerData`

#### グループ管理
- `CreateGroup(name string, permissions []string) error`
- `DeleteGroup(name string) error`
- `UpdateGroup(name string, permissions []string) error`
- `GetGroup(name string) *shared.Group`
- `GetAllGroups() map[string]*shared.Group`

#### プレイヤー-グループ関係
- `AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error`
- `RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error`
- `GetPlayerGroups(playerID uuid.UUID) []string`

#### 個人権限管理
- `AddPlayerPermission(playerID uuid.UUID, permission string) error`
- `RemovePlayerPermission(playerID uuid.UUID, permission string) error`
- `SetPlayerPermissions(playerID uuid.UUID, permissions []string) error`
- `GetPlayerPermissions(playerID uuid.UUID) []string`

#### グループ権限管理
- `AddPermissionToGroup(groupName, permission string) error`
- `RemovePermissionFromGroup(groupName, permission string) error`

#### システム管理
- `Save() error` - データを手動保存
- `Reload() error` - データを再読み込み
- `ClearCache()` - キャッシュをクリア
- `SetCacheEnabled(enabled bool)` - キャッシュの有効/無効切り替え
- `SetAutoSave(enabled bool)` - オートセーブの有効/無効切り替え

### デフォルトグループ
システム起動時に以下のデフォルトグループが作成されます：
- **default**: `["chat.send", "world.interact"]`
- **admin**: `["*"]` (全権限)

### データ構造

#### PermissionData
```go
type PermissionData struct {
    Groups  map[string]*Group         `json:"groups"`
    Players map[uuid.UUID]*PlayerData `json:"players"`
    Meta    *Metadata                 `json:"meta"`
}
```

#### PlayerData
```go
type PlayerData struct {
    PlayerID    uuid.UUID `json:"player_id"`
    PlayerName  string    `json:"player_name"`
    Groups      []string  `json:"groups"`
    Permissions []string  `json:"permissions"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### Group
```go
type Group struct {
    Name        string   `json:"name"`
    Permissions []string `json:"permissions"`
}
```

### 設定オプション

#### ManagerConfig
```go
type ManagerConfig struct {
    AutoSave bool
    Storage  shared.StorageConfig
    Cache    shared.CacheConfig
}
```

#### StorageConfig
```go
type StorageConfig struct {
    Path string // JSONファイルのパス
}
```

#### CacheConfig
```go
type CacheConfig struct {
    TTL             time.Duration // キャッシュの有効期限
    CleanupInterval time.Duration // クリーンアップ間隔
    Enabled         bool          // キャッシュ有効/無効
}
```