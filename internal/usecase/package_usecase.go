package usecase

import (
	"fmt"
	"godoc_mcp/internal/domain"
)

// PackageUseCase はパッケージ情報を取得するユースケース
type PackageUseCase struct {
	repo domain.PackageRepository
}

// NewPackageUseCase は新しいPackageUseCaseを作成する
func NewPackageUseCase(repo domain.PackageRepository) *PackageUseCase {
	return &PackageUseCase{
		repo: repo,
	}
}

// SearchPackages はパッケージを検索する
func (u *PackageUseCase) SearchPackages(query string) ([]domain.Package, error) {
	return u.repo.Search(query)
}

// GetPackageDetails はパッケージの詳細情報を取得する
func (u *PackageUseCase) GetPackageDetails(url string) (*domain.Package, error) {
	res, err := u.repo.GetDetails(url)
	fmt.Println(res)
	return res, err
}
