# Gateway Service - DDD Architecture

## ディレクトリ構造

```
services/gateway/
├── cmd/
│   └── server/
│       └── main.go                    # エントリーポイント
├── internal/
│   ├── domain/                        # Domain Layer
│   │   ├── user/
│   │   │   ├── entity.go             # User エンティティ
│   │   │   ├── value_objects.go      # Email, Password など
│   │   │   └── repository.go         # Repository interface
│   │   └── errors/
│   │       └── errors.go             # ドメインエラー
│   │
│   ├── application/                   # Application Layer
│   │   └── user/
│   │       ├── service.go            # Application Service
│   │       ├── dto.go                # DTOs (Input/Output)
│   │       └── mapper.go             # Domain ⟷ DTO 変換
│   │
│   ├── infrastructure/                # Infrastructure Layer
│   │   ├── grpc/
│   │   │   ├── client.go             # gRPC client factory
│   │   │   └── user_repository.go    # User Repository 実装
│   │   └── config/
│   │       └── config.go             # 設定管理
│   │
│   ├── presentation/                  # Presentation Layer
│   │   └── http/
│   │       ├── router.go             # ルーティング設定
│   │       ├── middleware/
│   │       │   └── auth.go           # 認証ミドルウェア
│   │       └── handler/
│   │           ├── user_handler.go   # HTTPハンドラー
│   │           └── mapper.go         # OpenAPI ⟷ DTO 変換
│   │
│   └── api/                           # 生成されたOpenAPI型
│       └── user.gen.go
│
├── go.mod
└── go.sum
```

## レイヤー責務

### Domain Layer (内側)
- ビジネスロジックの中核
- 他のレイヤーに依存しない
- エンティティ、値オブジェクト、リポジトリインターフェース

### Application Layer
- ユースケースの実装
- ドメインオブジェクトを orchestrate
- トランザクション境界

### Infrastructure Layer
- 外部サービスとの通信 (gRPC)
- リポジトリの実装
- 技術的な詳細

### Presentation Layer (外側)
- HTTPリクエスト/レスポンス処理
- OpenAPI型からDTOへの変換
- バリデーション
