package infrastructure

import (
	"fmt"
	"net/http"
	"strings"

	"godoc_mcp/internal/domain"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/PuerkitoBio/goquery"
)

// PackageRepository はパッケージ情報を取得するリポジトリの実装
type PackageRepository struct {
	baseURL string
}

// NewPackageRepository は新しいPackageRepositoryを作成する
func NewPackageRepository() *PackageRepository {
	return &PackageRepository{
		baseURL: "https://pkg.go.dev",
	}
}

// Search はパッケージを検索する
func (r *PackageRepository) Search(query string) ([]domain.Package, error) {
	url := fmt.Sprintf("%s/search?q=%s", r.baseURL, query)
	doc, err := r.scrape(url)
	if err != nil {
		return nil, err
	}

	var packages []domain.Package
	doc.Find("a[data-gtmc=\"search result\"]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			packages = append(packages, domain.Package{
				Name: strings.TrimPrefix(href, "/"),
				URL:  r.baseURL + href,
			})
		}
	})

	return packages, nil
}

// GetDetails はパッケージの詳細情報を取得する
func (r *PackageRepository) GetDetails(url string) (*domain.Package, error) {
	doc, err := r.scrape(url)
	if err != nil {
		return nil, err
	}

	// パッケージ情報を抽出
	pkg := &domain.Package{
		URL: url,
	}

	// タイトルを取得
	pkg.Name = doc.Find("h1").First().Text()

	// 説明を取得
	pkg.Description = doc.Find(".Documentation-overview").First().Text()

	// バージョン情報を取得
	pkg.Version = doc.Find(".DetailsHeader-version").First().Text()

	// ドキュメントのヘッダー部分をマークダウンに変換する
	headerSelector := ".go-Main-headerContent"
	headerContent := doc.Find(headerSelector)
	if headerContent.Length() == 0 {
		return nil, fmt.Errorf("header content not found")
	}

	// ドキュメントの本文をマークダウンに変換する
	bodySelector := ".Documentation-content"
	bodyContent := doc.Find(bodySelector)
	if bodyContent.Length() == 0 {
		return nil, fmt.Errorf("body content not found")
	}

	// 不要な要素を削除
	bodyContent.Find("script").Remove()
	bodyContent.Find("style").Remove()
	bodyContent.Find(".Documentation-indexColumn").Remove()
	bodyContent.Find(".Documentation-sidebar").Remove()

	// コンバーターの設定
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(
				commonmark.WithStrongDelimiter("**"),
				commonmark.WithEmDelimiter("*"),
			),
		),
	)

	// ヘッダーをマークダウンに変換
	headerHtml, err := headerContent.Html()
	if err != nil {
		return nil, err
	}
	headerHtml = strings.ReplaceAll(headerHtml, `href="/`, fmt.Sprintf(`href="%s/`, r.baseURL))
	headerMarkdown, err := conv.ConvertString(headerHtml)
	if err != nil {
		return nil, err
	}
	headerMarkdown = strings.TrimSpace(headerMarkdown)

	// ボディをマークダウンに変換
	bodyHtml, err := bodyContent.Html()
	if err != nil {
		return nil, err
	}
	bodyHtml = strings.ReplaceAll(bodyHtml, `href="/`, fmt.Sprintf(`href="%s/`, r.baseURL))
	bodyMarkdown, err := conv.ConvertString(bodyHtml)
	if err != nil {
		return nil, err
	}
	bodyMarkdown = strings.TrimSpace(bodyMarkdown)

	// ヘッダーとボディを結合
	pkg.Documentation = fmt.Sprintf(`# Overview

%s

# Documentation

%s`, headerMarkdown, bodyMarkdown)

	return pkg, nil
}

// scrape は指定されたURLからデータをスクレイピングする
func (r *PackageRepository) scrape(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
