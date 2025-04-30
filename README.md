# godoc_mcp

Goのパッケージ検索とメソッド仕様取得のためのMCPサーバー

## 概要

godoc_mcpは、Goのパッケージを検索し、そのパッケージに含まれるメソッドの詳細な仕様を取得できるMCPサーバーです。MCP（Model Control Protocol）を使用して、AIモデルと連携して動作します。

## 機能

- パッケージ検索
  - パッケージ名による検索
  - 検索結果の表示（パッケージ名、説明、URL）
- パッケージ詳細取得
  - メソッドの詳細仕様の取得
  - バージョン情報の表示
  - ドキュメントの取得

## 必要条件

- Go 1.16以上
- MCPクライアント

## セットアップ手順

### 1. MCPサーバーのクローン

```bash
git clone https://github.com/yourusername/godoc_mcp.git
cd godoc_mcp
```

### 2. MCPサーバーのビルド

```bash
go build -o godoc_mcp cmd/server/main.go
```

### 3. Cursorの設定

1. Cursorを起動し、設定画面を開きます
2. MCPサーバーの設定で、以下の情報を入力します：
   - サーバー名: godoc_mcp
   - ホスト: localhost
   - ポート: 8080（デフォルト）
   - パス: /mcp

## 使用方法

1. サーバーの起動
```bash
./godoc_mcp
```

2. パッケージの検索
```bash
# 例：fmtパッケージの検索
search_package packageName="fmt"
```

3. パッケージ詳細の取得
```bash
# 例：fmtパッケージの詳細取得
details url="https://pkg.go.dev/fmt"
```

## プロジェクト構造

```
godoc_mcp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── package.go
│   ├── infrastructure/
│   │   └── repository.go
│   └── usecase/
│       └── package.go
└── README.md
```