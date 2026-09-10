package service

import "cqrs/internal/core/product/app/command"

// ProductService is a product service.
type ProductService struct {
	uc *command.ProductUsecase
}

// NewProductService creates a new ProductService.
func NewProductService(uc *command.ProductUsecase) *ProductService {
	return &ProductService{uc: uc}
}
