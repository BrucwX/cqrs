package service

import (
	"github.com/go-kratos/kratos-layout/internal/core/product/biz/product"
)

// ProductService is a product service.
type ProductService struct {
	uc *product.ProductUsecase
}

// NewProductService creates a new ProductService.
func NewProductService(uc *product.ProductUsecase) *ProductService {
	return &ProductService{uc: uc}
}
