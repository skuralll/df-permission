# df-permission

Minecraft Bedrock Edition Dragonflyサーバー向けの権限管理ライブラリ

## 特徴

- **権限システム**: ワイルドカード権限（`*`、`prefix.*`）をサポート
- **グループベース管理**: グループによる権限管理
- **スレッドセーフ**: 並行処理に対応
- **設定**: カスタマイズ可能

## インストール

```bash
go get github.com/skuralll/df-permission
```

必要なGoバージョン: **Go 1.24+**

## クイックスタート

プレイヤーの参加時に自動で権限システムに登録する機能を提供します。

```go
package main

import (
    "log"
    
    "github.com/df-mc/dragonfly/server"
    "github.com/skuralll/df-permission/permission" 
    "github.com/skuralll/df-permission/dragonfly"
)

func main() {
    // 権限マネージャーを作成
    mgr := permission.NewManager()
    
    // オプション: 参加ハンドラーを作成
    joinHandler := dragonfly.NewJoinHandler(mgr, "default")
    
    // Dragonflyサーバーを設定...
    srv := server.New(nil, "logs")
    
    // サーバーの受け入れループ
    for p := range srv.Accept() {
        playerID := p.UUID()
        playerName := p.Name()
        
        // オプション: プレイヤー参加処理（自動登録、デフォルトグループ追加、名前更新）
        if err := joinHandler.HandlePlayerJoin(playerID, playerName); err != nil {
            log.Printf("Join handler error: %v", err)
        }
    }

    /* API使用例 */
    // プレイヤーをグループに追加
    playerID := uuid.New()
    mgr.AddPlayerToGroup(playerID, "Steve", "admin")
    
    // 権限をチェック
    if mgr.HasPermission(playerID, "any.permission") {
        // adminは全ての権限を持つ
        println("プレイヤーは権限を持っています")
    }
    
    // データを保存
    mgr.Save()
}

```

## 詳細な使用例

### カスタム設定

```go
mgr := permission.NewManager(
    permission.WithStorage("./custom_permissions.json"),
    permission.WithCache(45*time.Second),
    permission.WithAutoSave(true),
    permission.WithCacheCleanup(2*time.Minute),
)
```

### グループ管理

```go
// 新しいグループを作成
err := mgr.CreateGroup("vip", []string{"chat.color", "world.build.fast"})
if err != nil {
    if err == permission.ErrGroupAlreadyExists {
        // グループが既に存在
    }
}

// グループの権限を更新
err = mgr.UpdateGroup("vip", []string{"chat.color", "world.teleport"})

// グループを削除
err = mgr.DeleteGroup("vip")
```

### プレイヤー管理

```go
playerID := uuid.New()

// プレイヤーを作成
err := mgr.CreatePlayer(playerID, "PlayerName")

// プレイヤー名を更新
err = mgr.UpdatePlayerName(playerID, "NewPlayerName")

// プレイヤーをグループに追加
err = mgr.AddPlayerToGroup(playerID, "PlayerName", "vip")

// 個人権限を追加
err = mgr.AddPlayerPermission(playerID, "custom.permission")

// プレイヤーの全権限を設定
err = mgr.SetPlayerPermissions(playerID, []string{"perm1", "perm2"})

// プレイヤーの権限を取得
permissions := mgr.GetPlayerPermissions(playerID)
```

### 権限チェック

```go
// 基本的な権限チェック
if mgr.HasPermission(playerID, "chat.send") {
    // プレイヤーはチャットを送信可能
}

// ワイルドカード権限
if mgr.HasPermission(playerID, "world.build") {
    // "world.*" 権限を持つ場合もtrueになる
}
```

## 権限システム

### 権限パターン

- `*` : 全ての権限を付与（グローバルワイルドカード）
- `prefix.*` : `prefix.`で始まる全ての権限を付与
- `prefix.specific` : 特定の権限のみを付与

### デフォルトグループ

システム起動時に以下のグループが自動作成されます：

- **default**: `["chat.send", "world.interact"]`
  - ※ このライブラリ自体はチャットやワールド操作制限を行いません。
- **admin**: `["*"]` (全権限)

## API リファレンス

### PermissionManager インターフェース

```go
type PermissionManager interface {
    // 権限チェック
    HasPermission(playerID uuid.UUID, permission string) bool
    
    // プレイヤー管理
    CreatePlayer(playerID uuid.UUID, playerName string) error
    UpdatePlayerName(playerID uuid.UUID, newName string) error
    
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

### 設定オプション

```go
// ストレージファイルのパスを設定
permission.WithStorage("./permissions.json")

// オートセーブの有効/無効
permission.WithAutoSave(true)

// キャッシュの設定
permission.WithCache(30*time.Second)
permission.WithCacheEnabled(true)
permission.WithCacheCleanup(time.Minute)
```

### エラー処理

```go
var (
    ErrPlayerNotFound           // プレイヤーが見つからない
    ErrPlayerAlreadyExists      // プレイヤーが既に存在
    ErrGroupNotFound            // グループが見つからない
    ErrGroupAlreadyExists       // グループが既に存在
    ErrSystemGroupProtected     // システムグループは削除不可
    ErrPlayerPermissionNotFound // プレイヤー権限が見つからない
    ErrGroupPermissionNotFound  // グループ権限が見つからない
    ErrInvalidPermission        // 無効な権限
    ErrStorage                  // ストレージエラー
)
```

## テスト

```bash
go test ./...
```

特定のパッケージのテスト:
```bash
go test ./permission
go test ./internal/domain
```

## プロジェクト構造

```
df-permission/
├── permission/                 # 公開API
│   ├── permission.go           # メインインターフェース
│   ├── options.go              # 設定オプション
│   └── errors.go               # エラー定義
├── internal/                   # 内部実装
│   ├── domain/                 # ドメイン層
│   │   ├── checker.go          # 権限チェック
│   │   ├── matcher.go          # パターンマッチング
│   │   └── cache.go            # キャッシュ機能
│   ├── repository/             # データアクセス層
│   │   ├── storage.go          # ストレージインターフェース
│   │   └── json_storage.go     # JSON実装
│   ├── application/            # アプリケーション層
│   │   └── manager.go          # 統合管理
│   └── shared/                 # 共有定義
│       └── types.go            # データ型
└── docs/
    └── SPECIFICATION.md        # 詳細仕様
```

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) ファイルをご覧ください。

## 貢献

バグレポートや機能リクエストは [Issues](https://github.com/skuralll/df-permission/issues) まで。