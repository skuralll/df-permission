# Codebase Structure

## Architecture Overview
The project follows a clean, layered architecture pattern:

```
┌─────────────────────┐
│   Public API Layer  │  permission/ (Public interface)
├─────────────────────┤
│ Application Layer   │  internal/application/ (Business orchestration)
├─────────────────────┤
│   Domain Layer      │  internal/domain/ (Core business logic)
├─────────────────────┤
│ Repository Layer    │  internal/repository/ (Data access)
├─────────────────────┤
│   Shared Layer      │  internal/shared/ (Common types & errors)
└─────────────────────┘
```

## Directory Structure
```
df-permission/
├── permission/                 # 公開API (Public API Layer)
│   ├── permission.go           # メインインターフェース (Main interface)
│   ├── options.go              # 設定オプション (Configuration options)
│   ├── errors.go               # エラー定義 (Error definitions)
│   └── *_test.go               # Public API tests
├── internal/                   # 内部実装 (Internal implementation)
│   ├── application/            # アプリケーション層 (Application Layer)
│   │   ├── manager.go          # 統合管理 (Business orchestration)
│   │   └── manager_test.go     # Business logic tests
│   ├── domain/                 # ドメイン層 (Domain Layer)
│   │   ├── checker.go          # 権限チェック (Permission checking)
│   │   ├── matcher.go          # パターンマッチング (Pattern matching)
│   │   ├── cache.go            # キャッシュ機能 (Caching)
│   │   ├── storage.go          # Domain storage interface
│   │   └── *_test.go           # Domain logic tests
│   ├── repository/             # データアクセス層 (Repository Layer)
│   │   ├── storage.go          # ストレージインターフェース (Storage interface)
│   │   ├── json_storage.go     # JSON実装 (JSON implementation)
│   │   └── json_storage_test.go # Storage tests
│   ├── dragonfly/              # Dragonfly integration (currently empty)
│   │   ├── events/             # Event handlers
│   │   └── command/            # Command handlers
│   └── shared/                 # 共有定義 (Shared definitions)
│       ├── types.go            # データ型 (Data types)
│       └── errors.go           # Error definitions
├── docs/                       # Documentation
├── .github/workflows/          # CI/CD workflows
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── .goreleaser.yaml            # Release configuration
├── README.md                   # Project documentation
└── CLAUDE.md                   # Development guide
```

## Key Components
- **PermissionManager Interface**: 13 core methods for permission management
- **Layered Architecture**: Clear separation of concerns with domain-driven design
- **Thread-Safe Design**: RWMutex-based concurrent access control
- **Caching System**: TTL-based permission result caching
- **Storage Abstraction**: JSON file-based storage with interface for extensibility