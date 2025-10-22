# OpenAPI Code Generation

このプロジェクトでは、OpenAPI仕様からGoとTypeScriptの型定義を自動生成します。

## セットアップ

### 1. 必要なツールのインストール

```bash
# Go用のコードジェネレーターをインストール
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

# フロントエンド用の依存関係をインストール
cd web
npm install
```

### 2. コード生成

```bash
# プロジェクトルートから全ての型を生成
npm run codegen

# Go側の型のみ生成
npm run codegen:go

# TypeScript側の型のみ生成
npm run codegen:ts

# TypeScript型の自動生成（ファイル監視）
cd web
npm run codegen:watch
```

## ディレクトリ構成

```
kube-ec/
├── api/
│   └── openapi/
│       ├── user.yaml              # User Service API定義
│       ├── codegen-config.yaml    # Go生成設定
│       └── generated/
│           └── go/
│               └── user.gen.go    # 生成されたGo型とGinハンドラー
├── web/
│   └── src/
│       └── types/
│           └── api/
│               └── user.ts        # 生成されたTypeScript型
└── scripts/
    └── generate-go.sh             # Go生成スクリプト
```

## 使用方法

### Go (Gateway Service)

Gateway ServiceはDDD（ドメイン駆動設計）アーキテクチャで実装されています。

#### ディレクトリ構造
```
services/gateway/
├── cmd/server/main.go              # エントリーポイント
├── internal/
│   ├── domain/                     # ドメイン層
│   │   └── user/
│   │       ├── entity.go          # エンティティ
│   │       ├── value_objects.go   # 値オブジェクト
│   │       └── repository.go      # リポジトリインターフェース
│   ├── application/                # アプリケーション層
│   │   └── user/
│   │       ├── service.go         # ユースケース
│   │       ├── dto.go             # DTO
│   │       └── mapper.go          # DTO変換
│   ├── infrastructure/             # インフラ層
│   │   └── grpc/
│   │       ├── client.go          # gRPCクライアント
│   │       └── user_repository.go # リポジトリ実装
│   ├── presentation/               # プレゼンテーション層
│   │   └── http/
│   │       ├── router.go          # ルーティング
│   │       └── handler/
│   │           ├── user_handler.go # HTTPハンドラー
│   │           └── mapper.go       # OpenAPI型変換
│   └── api/                        # 生成されたOpenAPI型
│       └── user.gen.go
```

#### 実行方法
```bash
cd services/gateway
USER_SERVICE_ADDR=localhost:50051 PORT=8080 go run cmd/server/main.go
```

### TypeScript (Next.js Frontend)

このプロジェクトでは `openapi-fetch` を使用して、型安全なAPIクライアントを提供しています。

```typescript
import { register, login, getUser, updateUser } from '@/lib/api/user';

// 完全に型安全なAPI呼び出し
const response = await register({
  email: 'user@example.com',
  password: 'password123',
  name: 'John Doe',
  phone_number: '+1234567890'
});

// 型定義も直接インポート可能
import type { User, RegisterRequest } from '@/lib/api/user';

// または元の型定義から
import type { components } from '@/types/api/user';
type User = components['schemas']['User'];
```

#### カスタムAPIクライアントの作成

```typescript
import createClient from 'openapi-fetch';
import type { paths } from '@/types/api/user';

const client = createClient<paths>({
  baseUrl: 'http://localhost:8080/api/v1'
});

// 完全に型安全！
const { data, error } = await client.POST('/auth/register', {
  body: { email: '...', password: '...', name: '...' }
});
```

## OpenAPI仕様の編集

1. `api/openapi/user.yaml`を編集
2. `npm run codegen`を実行して型を再生成
3. 生成されたコードを確認してコミット

## 注意事項

- 生成されたファイル（`*.gen.go`, `web/src/types/api/*.ts`）は直接編集しないでください
- OpenAPI仕様を変更したら必ずコード生成を実行してください
- 生成されたファイルは`.gitignore`に追加されていません（意図的に含めています）
