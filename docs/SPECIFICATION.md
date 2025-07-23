# df-permission 仕様

## 概要

Minecraft Bedrock Edition Dragonflyサーバー向けのシンプルな権限管理システムです。

## 機能

- プレイヤーの権限チェック
- グループによる権限管理
- ワイルドカード権限 (* と prefix.* パターンのサポート)
- ファイルベースストレージ
- Dragonflyサーバーとの統合

## プロジェクト構造
```
df-permission/
├── permission.go      # Public API (Manager取得関数など最小限のみ公開、あとはinternalに隠蔽)
├── types.go          # 公開データ型
├── options.go        # オプションパターン
├── errors.go         # 公開エラー定義
│
└── internal/
    ├── domain/       # ビジネスロジック層 (権限システム本体)
    ├── repository/   # データアクセス層 (データR&W)
    ├── application/  # アプリケーション層 (Managerとpublic API)
    ├── dragonfly/    # Dragonfly統合層 (コマンドやイベントハンドラ)
    └── shared/
        ├── errors.go # 内部エラー定義
        └── utils.go  # ユーティリティ
```

## 使用例
```go
package main

import (
    "github.com/google/uuid"
    "github.com/skuralll/df-permission"
)

func main() {
    // 1. デフォルト設定でマネージャーを作成
    mgr := permission.QuickStart()
    
    // 2. グループと権限を作成
    mgr.CreateGroup("vip", []string{"chat.color", "world.build.fast"})
    mgr.CreateGroup("moderator", []string{"moderation.*"})
    
    // 3. プレイヤーをグループに追加
    playerID := uuid.New()
    mgr.AddPlayerToGroup(playerID, "Steve", "vip")
    
    // 4. 権限をチェック
    if mgr.HasPermission(playerID, "chat.color") {
        // プレイヤーはカラーチャットを使用可能
    }
    
    if mgr.HasPermission(playerID, "moderation.kick") {
        // プレイヤーは他のプレイヤーをキック可能（VIPでは false）
    }
}
```

---

## 権限システム
- `*` : 全ての権限を付与（グローバルワイルドカード）
- `prefix.*` : `prefix.`で始まる全ての権限を付与（プレフィックスワイルドカード）
- `prefix.specific` : 特定の権限のみを付与