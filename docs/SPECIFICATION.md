# df-permission 仕様

## 概要

Minecraft Bedrock Edition Dragonflyサーバー向けの包括的な権限管理ライブラリです。

## 機能

### コア機能
- プレイヤーの権限チェック（キャッシュ対応）
- グループによる権限管理
- ワイルドカード権限 (* と prefix.* パターンのサポート)
- ファイルベースストレージ（JSON）
- スレッドセーフな操作

### 権限管理
- **グループ管理**: 作成、削除、更新、権限追加/削除
- **プレイヤー管理**: 作成、削除、グループ参加/離脱、個人権限設定
- **個人権限**: プレイヤーに直接権限を付与/削除
- **一括操作**: 権限の一括設定・取得

### パフォーマンス機能
- **権限チェックキャッシュ**: TTL対応、自動クリーンアップ
- **オートセーブ**: 変更時の自動保存
- **データ整合性**: キャッシュ無効化による一貫性保証

## プロジェクト構造
```
df-permission/
├── permission.go      # 公開API
├── types.go          # 公開データ型
├── options.go        # オプションパターン
├── errors.go         # 公開エラー定義
│
└── internal/
    ├── domain/       # ドメイン層 (権限システムの核心ロジック)
    │   ├── checker.go     # 権限チェックロジック
    │   ├── matcher.go     # パターンマッチング
    │   └── cache.go       # 権限チェックキャッシュ
    ├── repository/   # データアクセス層
    │   ├── storage.go     # Storage interface
    │   ├── json_storage.go # JSON実装
    │   └── errors.go      # ストレージエラー
    ├── application/  # アプリケーション層
    └── shared/       # 共有定義
        ├── types.go       # 共有データ型
        └── config.go      # 設定構造体
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