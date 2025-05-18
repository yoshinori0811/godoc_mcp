# godoc_mcp

Goのパッケージ検索とメソッド仕様取得のためのMCPサーバー

## 概要

godoc_mcpは、Goのパッケージを検索し、そのパッケージに含まれるメソッドの詳細な仕様を取得できるMCPサーバーです。MCP（Model Context Protocol）を使用して、AIモデルと連携して動作します。

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
git clone https://github.com/yoshinori0811/godoc_mcp.git
cd godoc_mcp
```

### 2. MCPサーバーのビルド

```bash
go build -o godoc_mcp.exe cmd/server/main.go
```

### 3. Cursorの設定

1. Cursorを起動し、設定画面（Cursor Settings）を開きます。
   - **Windows:** `Ctrl + ,` を押す
   - **Mac:** `Cmd + ,` を押す
2. 左側のメニューから「MCP」タブを選択し、「+ Add new global MCP server」をクリックします。
3. 以下の内容で `mcp.json` を作成します。　

```json
{
  "mcpServers": {
    "godoc": {
      "command": "<godoc_mcp.exeのフルパス>"
    }
  }
}
```

※ `<godoc_mcp.exeのフルパス>` には、ビルドした `godoc_mcp.exe` の絶対パスを指定してください。

<img src="https://raw.github.com/wiki/yoshinori0811/godoc_mcp/images/godoc_mcp_demo.gif" style="width: 100%;" >

## 使用方法

<img src="https://raw.github.com/wiki/yoshinori0811/godoc_mcp/images/mcp_settings.png" style="width: 100%;" >
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