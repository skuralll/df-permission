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

#### 公開API
- ✅ **permission.go**: 公開APIエントリーポイント
- ✅ **errors.go**: 公開エラー定義
- ✅ **PermissionManager**: ラッパー方式インターフェース

### 🚧 未実装機能

#### 高度な機能
- ❌ **options.go**: オプションパターン実装
- ❌ **QuickStart()**: 簡単セットアップ関数

## プロジェクト構造

### 実装済み構造
```
df-permission/
├── permission.go               # ✅ 公開API（PermissionManagerインターフェース）
├── errors.go                   # ✅ 公開エラー定義（6つのエラー型）
├── types.go                    # ✅ 公開データ型（設定構造体の再エクスポート）
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
        ├── types.go            # ✅ 共有データ型（設定構造体含む）
```

### 未実装の拡張機能
```
df-permission/
├── options.go         # ❌ オプションパターン（未実装）
└── quickstart.go      # ❌ 簡単セットアップ（未実装）
```

## 使用例

### 公開API使用方法（推奨）
```go
package main

import (
    "github.com/google/uuid"
    "github.com/skuralll/df-permission"
    "time"
)

func main() {
    // 1. 設定を作成
    config := dfpermission.ManagerConfig{
        AutoSave: true,
        Storage: dfpermission.StorageConfig{
            Path: "./permissions.json",
        },
        Cache: dfpermission.CacheConfig{
            Enabled:         true,
            TTL:             30 * time.Second,
            CleanupInterval: time.Minute,
        },
    }
    
    // 2. 公開APIでマネージャーを作成
    mgr := dfpermission.NewManager(config)
    
    // 3. グループを作成
    err := mgr.CreateGroup("vip", []string{"chat.color", "world.build.fast"})
    if err != nil {
        // エラーハンドリング（公開エラー型）
        if err == dfpermission.ErrGroupAlreadyExists {
            // グループが既に存在
        }
    }
    
    // 4. プレイヤーをグループに追加
    playerID := uuid.New()
    err = mgr.AddPlayerToGroup(playerID, "Steve", "vip")
    if err != nil {
        // エラーハンドリング
    }
    
    // 5. 権限をチェック
    if mgr.HasPermission(playerID, "chat.color") {
        // プレイヤーはカラーチャットを使用可能
    }
    
    // 6. グループに権限を動的に追加
    err = mgr.AddPermissionToGroup("vip", "world.teleport")
    if err != nil {
        // エラーハンドリング
    }
    
    // 7. プレイヤーは新しい権限を自動的に取得
    if mgr.HasPermission(playerID, "world.teleport") {
        // VIPグループのプレイヤーはテレポート可能
    }
    
    // 8. 個人権限を追加
    err = mgr.AddPlayerPermission(playerID, "custom.permission")
    if err == dfpermission.ErrPlayerNotFound {
        // プレイヤーが存在しない
    }
    
    // 9. グループから権限を削除
    err = mgr.RemovePermissionFromGroup("vip", "world.build.fast")
    if err == dfpermission.ErrPermissionNotFound {
        // グループが権限を持っていない
    }
    
    // 10. データを保存
    mgr.Save()
}
```

### 内部API直接使用（開発・テスト用）
```go
package main

import (
    "github.com/google/uuid"
    "github.com/skuralll/df-permission/internal/application"
    "github.com/skuralll/df-permission/internal/shared"
    "time"
)

func main() {
    // 内部APIを直接使用（非推奨）
    config := shared.ManagerConfig{
        // 設定内容は同じ
    }
    mgr := application.NewManager(config)
    // 全ての内部メソッドにアクセス可能
}
```

### 未実装の拡張予定API
```go
package main

import (
    "github.com/skuralll/df-permission"
)

func main() {
    // 将来の拡張機能（未実装）
    
    // オプションパターン
    mgr := dfpermission.NewManager(
        dfpermission.WithStorage("./permissions.json"),
        dfpermission.WithCache(30*time.Second),
        dfpermission.WithAutoSave(true),
    )
    
    // クイックスタート
    mgr := dfpermission.QuickStart()
}
```

---

## 詳細仕様

### 権限システム
- `*` : 全ての権限を付与（グローバルワイルドカード）
- `prefix.*` : `prefix.`で始まる全ての権限を付与（プレフィックスワイルドカード）
- `prefix.specific` : 特定の権限のみを付与

### 公開API（permission.go）

#### PermissionManagerインターフェース
```go
type PermissionManager interface {
    // 権限チェック
    HasPermission(playerID uuid.UUID, permission string) bool
    
    // グループ管理
    CreateGroup(name string, permissions []string) error
    DeleteGroup(name string) error
    UpdateGroup(name string, permissions []string) error
    
    // プレイヤー-グループ関係
    AddPlayerToGroup(playerID uuid.UUID, playerName, groupName string) error
    RemovePlayerFromGroup(playerID uuid.UUID, groupName string) error
    
    // 個人権限管理
    AddPlayerPermission(playerID uuid.UUID, permission string) error
    RemovePlayerPermission(playerID uuid.UUID, permission string) error
    SetPlayerPermissions(playerID uuid.UUID, permissions []string) error
    GetPlayerPermissions(playerID uuid.UUID) []string
    
    // グループ権限管理
    AddPermissionToGroup(groupName, permission string) error
    RemovePermissionFromGroup(groupName, permission string) error
    
    // システム操作
    Save() error
}
```

#### 公開エラー型（errors.go）
```go
var (
    ErrPlayerNotFound       = errors.New("player not found")
    ErrGroupNotFound        = errors.New("group not found")
    ErrGroupAlreadyExists   = errors.New("group already exists")
    ErrSystemGroupProtected = errors.New("system group cannot be deleted")
    ErrPlayerAlreadyExists  = errors.New("player already exists")
    ErrPermissionNotFound   = errors.New("permission not found")
)
```

#### 公開APIの特徴
- **13個のコアメソッド**: 権限管理に必要な機能を厳選
- **エラー変換**: 内部エラーを6つの公開エラー型に変換
- **安定したインターフェース**: 内部実装変更に影響されない
- **キャッシュ対応**: 権限チェックの高速化

#### ファクトリー関数
- `NewManager(config ManagerConfig) PermissionManager`
  - 設定からPermissionManagerを作成
  - 内部実装をラップして安定したAPIを提供

### 内部API（internal/application/manager.go）

#### 拡張機能（公開APIでは利用不可）
- `CreatePlayer(playerID uuid.UUID, playerName string) error`
- `RemovePlayer(playerID uuid.UUID) error`
- `PlayerExists(playerID uuid.UUID) bool`
- `GetPlayerData(playerID uuid.UUID) *shared.PlayerData`
- `GetAllPlayers() map[uuid.UUID]*shared.PlayerData`
- `GetGroup(name string) *shared.Group`
- `GetAllGroups() map[string]*shared.Group`
- `GetPlayerGroups(playerID uuid.UUID) []string`
- `Reload() error`
- `ClearCache()`
- `SetCacheEnabled(enabled bool)`
- `SetAutoSave(enabled bool)`

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
    AutoSave bool              // 変更時の自動保存
    Storage  StorageConfig     // ストレージ設定
    Cache    CacheConfig       // キャッシュ設定
}
```

#### StorageConfig
```go
type StorageConfig struct {
    Path string // JSONファイルのパス（例: "./permissions.json"）
}
```

#### CacheConfig
```go
type CacheConfig struct {
    TTL             time.Duration // キャッシュの有効期限（例: 30*time.Second）
    CleanupInterval time.Duration // クリーンアップ間隔（例: time.Minute）
    Enabled         bool          // キャッシュ有効/無効
}
```

## アーキテクチャ設計

### レイヤード構造
1. **公開API層**: `permission.go`, `errors.go`, `types.go`
   - 安定したインターフェース
   - エラー変換とハンドリング
   - 型の再エクスポート

2. **アプリケーション層**: `internal/application/`
   - ビジネスロジックの統合
   - トランザクション管理
   - キャッシュ制御

3. **ドメイン層**: `internal/domain/`
   - 権限チェックロジック
   - パターンマッチング
   - キャッシュ機能

4. **リポジトリ層**: `internal/repository/`
   - データ永続化
   - ストレージ抽象化

5. **共有層**: `internal/shared/`
   - 共通データ型
   - 設定構造体

### 設計原則
- **依存性逆転**: 上位レイヤーが下位レイヤーのインターフェースに依存
- **単一責任**: 各レイヤーが明確な責任を持つ
- **開放閉鎖**: 拡張に開放、修正に閉鎖
- **インターフェース分離**: 必要最小限のメソッドのみ公開