package main

import (
	"context"
	"encoding/json"
	"fmt"
	"godoc_mcp/internal/infrastructure"
	"godoc_mcp/internal/usecase"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// リポジトリの初期化
	repo := infrastructure.NewPackageRepository()

	// ユースケースの初期化
	useCase := usecase.NewPackageUseCase(repo)

	// MCPサーバーの設定
	s := server.NewMCPServer(
		"Go Package Search",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	tool := mcp.NewTool(
		"search_package",
		mcp.WithDescription("godoc_mcp is a tool for searching and getting details of Go packages"),
		mcp.WithString("packageName",
			mcp.Required(),
			mcp.Description("godoc_mcp search package"),
		),
	)

	// パッケージ検索ハンドラ
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		packageName := req.Params.Arguments["packageName"].(string)
		packages, err := useCase.SearchPackages(packageName)
		if err != nil {
			return nil, err
		}

		// 候補リストを見やすい形式に整形
		var results []map[string]string
		for i, pkg := range packages {
			results = append(results, map[string]string{
				"id":          fmt.Sprintf("%d", i),
				"name":        pkg.Name,
				"description": pkg.Description,
				"url":         pkg.URL,
			})
		}

		jsonBytes, err := json.Marshal(results)
		if err != nil {
			return nil, fmt.Errorf("JSON変換エラー: %w", err)
		}
		return mcp.NewToolResultText(fmt.Sprintf(`\n以下のパッケージが見つかりました。詳細を確認したいパッケージの情報を"details"ツールで取得できます：\n%s\n"details"ツールを使用する際は、上記のURLを指定してください。`, string(jsonBytes))), nil
	})

	// パッケージ詳細取得ツール
	s.AddTool(mcp.NewTool(
		"details",
		mcp.WithDescription(`指定されたパッケージの詳細情報を返します。

例えばWithDescriptionメソッドの場合は以下のように回答してください。

===
パッケージの詳細情報を確認したところ、WithDescriptionメソッドはToolOption型の関数として定義されており、以下の仕様を持っています：
Apply to main.go
このメソッドの主な特徴は：
引数：
description string: ツールの説明文を表す文字列
戻り値：
ToolOption: ツールの設定を変更する関数
機能：
ツールに説明文を追加するためのオプション関数
ツールの動作を人間が読みやすい形で説明するための説明文を設定
この説明文はツールの使用方法や目的を明確にするために使用される
===

"details"ツールで返されるマークダウンを以下のように抽出してください。
===
#### func [WithDescription](https://github.com/mark3labs/mcp-go/blob/v0.23.1/mcp/tools.go#L187) [¶](#WithDescription "Go to WithDescription") added in v0.5.1

===
func WithDescription(description string) ToolOption
===

WithDescription adds a description to the Tool. The description should provide a clear, human-readable explanation of what the tool does.
===

`),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("パッケージのURL"),
		),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url := req.Params.Arguments["url"].(string)
		pkg, err := useCase.GetPackageDetails(url)
		if err != nil {
			return nil, err
		}

		// 詳細情報を整形
		details := map[string]string{
			"name":          pkg.Name,
			"description":   pkg.Description,
			"version":       pkg.Version,
			"documentation": pkg.Documentation,
		}

		jsonBytes, err := json.Marshal(details)
		if err != nil {
			return nil, fmt.Errorf("JSON変換エラー: %w", err)
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
