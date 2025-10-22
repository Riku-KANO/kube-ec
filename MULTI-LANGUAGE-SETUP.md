# マルチ言語プロジェクト構成 🌐

## プロジェクト構造

```
kube-ec/
├── MODULE.bazel              # Bazel Bzlmod設定（全言語対応）
├── BUILD.bazel               # ルートビルド定義
├── .bazelrc                  # Bazel設定
│
├── services/
│   ├── gateway/              # Go service
│   │   ├── go.mod           # Gateway固有の依存関係
│   │   ├── go.sum
│   │   └── BUILD.bazel
│   │
│   ├── analytics/            # 将来: Python service
│   │   ├── requirements.txt
│   │   └── BUILD.bazel
│   │
│   └── metrics/              # 将来: Rust service
│       ├── Cargo.toml
│       └── BUILD.bazel
│
├── proto/                    # 共通プロトコル定義
│   ├── go.mod               # Proto生成コード用
│   ├── user/
│   ├── product/
│   └── ...
│
└── web/                      # Next.js frontend
    ├── package.json
    └── BUILD.bazel
```

## 設計原則

### ✅ 各サービスが独立したビルド設定を持つ

**Go Services**:
- 各サービスディレクトリに`go.mod`を配置
- サービス固有の依存関係を管理

**Python Services（将来）**:
- `requirements.txt` または `pyproject.toml`
- 仮想環境を個別管理

**Rust Services（将来）**:
- `Cargo.toml`でクレート管理
- Workspace機能も利用可能

### ✅ Bazelで統一的にビルド

MODULE.bazelで全言語のツールチェーンを定義:

```python
# 現在
bazel_dep(name = "rules_go", version = "0.50.1")
bazel_dep(name = "gazelle", version = "0.39.1")

# 将来追加予定
# bazel_dep(name = "rules_python", version = "0.36.0")
# bazel_dep(name = "rules_rust", version = "0.53.0")
```

### ✅ rootにはビルドシステムの設定のみ

rootディレクトリには:
- ✅ `MODULE.bazel` - Bazel設定
- ✅ `BUILD.bazel` - ルートビルド定義
- ✅ `.bazelrc` - Bazel実行設定
- ✅ `package.json` - npm scripts（コード生成など）
- ❌ **言語固有のファイル (go.mod, requirements.txt, etc) は置かない**

## Go依存関係の管理

### Bazelの動作

```python
# MODULE.bazel
go_deps = use_extension("@gazelle//:extensions.bzl", "go_deps")
go_deps.from_file(go_mod = "//services/gateway:go.mod")
```

- Gateway serviceの`go.mod`を参照
- 他のGoサービスを追加する際は、そのgo.modも参照可能
- 各サービスは独立して依存関係を管理

### 通常のGo開発

```bash
cd services/gateway
go mod tidy
go build ./cmd/server
```

サービスディレクトリ内で通常のGo開発が可能

## 将来の拡張

### Python サービスの追加

1. **サービス作成**:
```bash
mkdir -p services/analytics
cd services/analytics
python -m venv .venv
pip install -r requirements.txt
```

2. **MODULE.bazel更新**:
```python
bazel_dep(name = "rules_python", version = "0.36.0")

python = use_extension("@rules_python//python:extensions.bzl", "python")
python.toolchain(python_version = "3.11")
```

3. **BUILD.bazel作成**:
```python
load("@rules_python//python:defs.bzl", "py_binary", "py_library")

py_binary(
    name = "analytics",
    srcs = ["main.py"],
    deps = [":lib"],
)
```

### Rust サービスの追加

1. **サービス作成**:
```bash
cd services
cargo new metrics
```

2. **MODULE.bazel更新**:
```python
bazel_dep(name = "rules_rust", version = "0.53.0")

rust = use_extension("@rules_rust//rust:extensions.bzl", "rust")
rust.toolchain(edition = "2021")
```

3. **BUILD.bazel作成**:
```python
load("@rules_rust//rust:defs.bzl", "rust_binary")

rust_binary(
    name = "metrics",
    srcs = ["src/main.rs"],
)
```

## ビルドコマンド

### 全サービスビルド
```bash
# Go services
bazel build //services/gateway:gateway

# Python service (将来)
bazel build //services/analytics:analytics

# Rust service (将来)
bazel build //services/metrics:metrics

# すべて
bazel build //services/...
```

### 個別開発
```bash
# Go
cd services/gateway && go run cmd/server/main.go

# Python (将来)
cd services/analytics && python main.py

# Rust (将来)
cd services/metrics && cargo run
```

## メリット

1. **言語の独立性**
   - 各言語のエコシステムをそのまま利用
   - 言語固有のツールが正常に動作

2. **段階的な導入**
   - 新しい言語を段階的に追加可能
   - 既存サービスに影響なし

3. **Bazelの統一性**
   - すべてのサービスを統一的にビルド
   - 依存関係の可視化
   - リモートキャッシュで高速ビルド

4. **チームの柔軟性**
   - GoチームはGoツールを使用
   - PythonチームはPythonツールを使用
   - Bazelで統合ビルド

## まとめ

✅ **rootにgo.modを置かない構成に変更完了**
✅ **各サービスが独立した設定を持つ**
✅ **Bazel Bzlmodで統一的にビルド可能**
✅ **将来のPython/Rust追加に対応**

これでマルチ言語マイクロサービスプロジェクトとして最適な構成になりました！
