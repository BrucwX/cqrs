package command

import "cqrs/internal/core/product/domain/repo"

// ProductUsecase is the product usecase.
type ProductUsecase struct {
	Repo repo.ProductRepo
}

// NewProductUsecase creates a new ProductUsecase.
func NewProductUsecase(r repo.ProductRepo) *ProductUsecase {
	return &ProductUsecase{Repo: r}
}
