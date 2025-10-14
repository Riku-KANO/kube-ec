# GCP セットアップガイド

このドキュメントでは、Kube ECをGoogle Cloud Platform (GCP)にデプロイする手順を説明します。

## 前提条件

- GCPアカウント
- `gcloud` CLIがインストールされていること
- `kubectl` がインストールされていること
- Dockerがインストールされていること

## 1. GCPプロジェクトのセットアップ

```bash
# GCPプロジェクトの作成
export PROJECT_ID="kube-ec-project"
gcloud projects create $PROJECT_ID

# プロジェクトを設定
gcloud config set project $PROJECT_ID

# 必要なAPIを有効化
gcloud services enable \
  container.googleapis.com \
  containerregistry.googleapis.com \
  sqladmin.googleapis.com \
  compute.googleapis.com
```

## 2. GKE クラスターの作成

```bash
# GKEクラスターを作成
gcloud container clusters create kube-ec-cluster \
  --zone asia-northeast1-a \
  --num-nodes 3 \
  --machine-type e2-medium \
  --enable-autoscaling \
  --min-nodes 3 \
  --max-nodes 10

# クラスターに接続
gcloud container clusters get-credentials kube-ec-cluster \
  --zone asia-northeast1-a
```

## 3. Cloud SQL (PostgreSQL) のセットアップ

```bash
# Cloud SQLインスタンスを作成
gcloud sql instances create kube-ec-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=asia-northeast1

# rootパスワードを設定
gcloud sql users set-password postgres \
  --instance=kube-ec-db \
  --password=YOUR_SECURE_PASSWORD

# データベースを作成
gcloud sql databases create kube_ec_prod \
  --instance=kube-ec-db

# Cloud SQL Proxyを使用して接続
# または、プライベートIPでVPC接続を設定
```

## 4. Dockerイメージのビルドとプッシュ

```bash
# GCRに認証
gcloud auth configure-docker

# 各サービスのイメージをビルド
cd services/product
docker build -t gcr.io/$PROJECT_ID/product-service:latest .
docker push gcr.io/$PROJECT_ID/product-service:latest

cd ../user
docker build -t gcr.io/$PROJECT_ID/user-service:latest .
docker push gcr.io/$PROJECT_ID/user-service:latest

cd ../order
docker build -t gcr.io/$PROJECT_ID/order-service:latest .
docker push gcr.io/$PROJECT_ID/order-service:latest

cd ../payment
docker build -t gcr.io/$PROJECT_ID/payment-service:latest .
docker push gcr.io/$PROJECT_ID/payment-service:latest

cd ../gateway
docker build -t gcr.io/$PROJECT_ID/gateway-service:latest .
docker push gcr.io/$PROJECT_ID/gateway-service:latest

# フロントエンド
cd ../../web
docker build -t gcr.io/$PROJECT_ID/web-frontend:latest .
docker push gcr.io/$PROJECT_ID/web-frontend:latest
```

## 5. Kubernetes Secretsの作成

```bash
# データベース接続文字列
kubectl create secret generic db-secret \
  --from-literal=database-url='postgres://postgres:YOUR_PASSWORD@CLOUD_SQL_IP:5432/kube_ec_prod'

# JWT Secret
kubectl create secret generic jwt-secret \
  --from-literal=secret='your-jwt-secret-key-here'
```

## 6. Kubernetesマニフェストのデプロイ

```bash
cd deploy/kubernetes

# YAMLファイルのPROJECT_IDを置換
find . -name "*.yaml" -exec sed -i "s/YOUR_PROJECT_ID/$PROJECT_ID/g" {} +

# すべてのリソースをデプロイ
kubectl apply -k .

# デプロイメントの状態を確認
kubectl get deployments
kubectl get services
kubectl get pods
```

## 7. 外部IPアドレスの確認

```bash
# Gateway ServiceのIPを取得
kubectl get service gateway-service

# Web FrontendのIPを取得
kubectl get service web-frontend
```

## 8. データベースマイグレーション

```bash
# Productサービス用テーブル
kubectl exec -it deployment/product-service -- /bin/sh
# SQL実行

CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price_currency VARCHAR(3),
    price_amount BIGINT,
    stock_quantity INTEGER,
    category VARCHAR(100),
    sku VARCHAR(100) UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

# Userサービス用テーブル
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

# Orderサービス用テーブル
CREATE TABLE orders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    items JSONB,
    total_currency VARCHAR(3),
    total_amount BIGINT,
    status VARCHAR(50),
    payment_id VARCHAR(36),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

# Paymentサービス用テーブル
CREATE TABLE payments (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    amount_currency VARCHAR(3),
    amount BIGINT,
    status VARCHAR(50),
    method VARCHAR(50),
    transaction_id VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

## 9. 監視とログ

```bash
# ログの確認
kubectl logs -f deployment/gateway-service
kubectl logs -f deployment/product-service

# Cloud Loggingで確認
gcloud logging read "resource.type=k8s_container AND resource.labels.cluster_name=kube-ec-cluster"
```

## 10. スケーリング

```bash
# 手動スケーリング
kubectl scale deployment product-service --replicas=5

# Horizontal Pod Autoscaler (HPA) の設定
kubectl autoscale deployment product-service \
  --cpu-percent=70 \
  --min=2 \
  --max=10
```

## トラブルシューティング

### Podが起動しない場合

```bash
kubectl describe pod POD_NAME
kubectl logs POD_NAME
```

### サービスに接続できない場合

```bash
# サービスのエンドポイント確認
kubectl get endpoints

# Podのネットワーク確認
kubectl exec -it POD_NAME -- /bin/sh
```

## クリーンアップ

```bash
# すべてのリソースを削除
kubectl delete -k deploy/kubernetes/

# GKEクラスターを削除
gcloud container clusters delete kube-ec-cluster --zone asia-northeast1-a

# Cloud SQLインスタンスを削除
gcloud sql instances delete kube-ec-db
```

## コスト最適化

- 開発環境では、プリエンプティブルノードを使用
- 不要な環境は停止する
- Cloud SQLのサイズを適切に設定
- LoadBalancerの代わりにIngressを検討

## セキュリティ推奨事項

- Secretsは暗号化して管理
- Network Policyを設定
- RBAC (Role-Based Access Control) を設定
- TLS/HTTPSを有効化
- Workload Identityを使用
