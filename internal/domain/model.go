package domain

// Package はGoパッケージの情報を表すドメインモデル
type Package struct {
	Name          string
	Description   string
	URL           string
	Version       string
	Documentation string
}

// PackageRepository はパッケージ情報を取得するためのリポジトリインターフェース
type PackageRepository interface {
	Search(query string) ([]Package, error)
	GetDetails(url string) (*Package, error)
}
