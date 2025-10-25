# OpenAPI スキーマ駆動開発 セットアップ完了 ✅

## 🎉 完成したもの

### 1. OpenAPI定義
- `api/openapi/user.yaml` - User Service API の完全な定義

### 2. コード生成スクリプト
- `scripts/generate-proto.sh` - Protobuf から Go コード生成
- `scripts/generate-go.sh` - OpenAPI から Go 型生成
- `npm run codegen` - OpenAPI から TypeScript 型生成

### 3. Gateway Service (DDD アーキテクチャ)
```
services/gateway/
├── cmd/server/main.go              # エントリーポイント
├── internal/
│   ├── domain/                     # Domain Layer
│   │   ├── user/                  # User ドメイン
│   │   │   ├── entity.go
│   │   │   ├── value_objects.go
│   │   │   └── repository.go
│   │   └── errors/                # ドメインエラー
│   ├── application/                # Application Layer
│   │   └── user/                  # User ユースケース
│   │       ├── service.go
│   │       ├── dto.go
│   │       └── mapper.go
│   ├── infrastructure/             # Infrastructure Layer
│   │   └── grpc/                  # gRPC 実装
│   │       ├── client.go
│   │       └── user_repository.go
│   ├── presentation/               # Presentation Layer
│   │   └── http/
│   │       ├── router.go
│   │       └── handler/
│   │           ├── user_handler.go
│   │           └── mapper.go
│   └── api/                        # 生成された OpenAPI 型
│       └── user.gen.go
└── gateway.exe                      # ビルド済みバイナリ (23MB)
```

### 4. フロントエンド (Next.js)
```
web/
├── src/
│   ├── lib/api/
│   │   ├── client.ts              # openapi-fetch クライアント
│   │   └── user.ts                # User API 関数
│   ├── types/api/
│   │   └── user.ts                # 生成された TypeScript 型
│   └── app/example/
│       └── page.tsx               # 使用例
└── package.json                    # codegen スクリプト含む
```

## 🚀 使い方

### コード生成

```bash
# 全ての型を生成
npm run codegen

# Go 型のみ
npm run codegen:go

# TypeScript 型のみ
npm run codegen:ts

# Protobuf から Go コード生成
bash scripts/generate-proto.sh
```

### Gateway Service 実行

```bash
cd services/gateway

# 環境変数を設定して実行
USER_SERVICE_ADDR=localhost:50051 PORT=8080 ./gateway.exe

# または Go で直接実行
USER_SERVICE_ADDR=localhost:50051 PORT=8080 go run cmd/server/main.go
```

### フロントエンド実行

```bash
cd web

# 依存関係インストール
npm install

# 開発サーバー起動
npm run dev
```

## 📝 API エンドポイント

Gateway は以下のエンドポイントを提供します:

- `POST /api/v1/auth/register` - ユーザー登録
- `POST /api/v1/auth/login` - ログイン
- `GET /api/v1/users/{id}` - ユーザー取得
- `PUT /api/v1/users/{id}` - ユーザー更新
- `DELETE /api/v1/users/{id}` - ユーザー削除
- `GET /health` - ヘルスチェック

## 💡 開発フロー

1. **OpenAPI 定義を編集**: `api/openapi/user.yaml`
2. **型を生成**: `npm run codegen`
3. **フロントエンドで使用**:
   ```typescript
   import { register, login } from '@/lib/api/user';

   const response = await register({
     email: 'user@example.com',
     password: 'password123',
     name: 'John Doe'
   });
   ```
4. **バックエンドはDDDレイヤーで実装済み**

## 🔧 技術スタック

- **Go 1.25+**: バックエンド
- **Gin**: HTTP フレームワーク
- **gRPC**: マイクロサービス間通信
- **OpenAPI 3.0**: API 定義
- **oapi-codegen**: Go 型生成
- **openapi-typescript**: TypeScript 型生成
- **openapi-fetch**: 型安全な fetch クライアント
- **Next.js 15**: フロントエンド
- **React 19**: UI

## ✅ ビルド確認済み

- ✅ Protobuf コード生成
- ✅ OpenAPI Go 型生成
- ✅ Gateway Service ビルド成功 (23MB)
- ✅ DDD レイヤー実装完了
- ✅ TypeScript 型生成設定完了
- ✅ フロントエンド API クライアント実装

## 📚 ドキュメント

- `README-CODEGEN.md` - コード生成の詳細
- `services/gateway/DDD-STRUCTURE.md` - DDD アーキテクチャ説明
