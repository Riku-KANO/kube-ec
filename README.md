# Kube EC - Kubernetes ECサイト練習
※練習

Bazel + Kubernetes + Go + Next.js + GCP で構築したマイクロサービスアーキテクチャのECサイトです。

## アーキテクチャ

### バックエンド（Goマイクロサービス）
- **Product Service**: 商品管理（gRPC）
- **User Service**: ユーザー認証・管理（gRPC + JWT）
- **Order Service**: 注文処理（gRPC）
- **Payment Service**: 決済処理（gRPC）
- **Gateway Service**: REST API Gateway（Gin + gRPC client）

### フロントエンド
- **Next.js 15**: App Router構成
- **TailwindCSS**: スタイリング
- **Zustand**: 状態管理

### インフラ
- **Kubernetes (GKE)**: コンテナオーケストレーション
- **Cloud SQL (PostgreSQL)**: データベース
- **Bazel**: ビルドシステム

## 技術スタック

- **ビルドツール**: Bazel
- **バックエンド**: Go 1.25
- **通信**: gRPC, REST
- **フロントエンド**: Next.js 15, React 19, TypeScript
- **データベース**: PostgreSQL
- **認証**: JWT
- **インフラ**: Kubernetes, GCP

## セットアップ

詳細は [docs/gcp-setup.md](./docs/gcp-setup.md) を参照してください。

### ローカル開発

```bash
# 依存関係のインストール
cd services/product && go mod download
cd ../user && go mod download
# ... 他のサービスも同様

# フロントエンド
cd web && npm install

# サービスの起動
DATABASE_URL="postgres://..." GRPC_PORT=50051 go run .
```

### GCPへのデプロイ

```bash
# GKEクラスター作成
gcloud container clusters create kube-ec-cluster --zone asia-northeast1-a

# デプロイ
kubectl apply -k deploy/kubernetes/
```
