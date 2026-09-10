package command

import (
	"cqrs/internal/core/product/domain/repo/command"
	"cqrs/internal/core/product/domain/repo/query"
)

// ProductUsecase is the product usecase.
type ProductUsecase struct {
	Query   query.ProductQuery
	Command command.ProductCommand
}

// NewProductUsecase creates a new ProductUsecase.
func NewProductUsecase(q query.ProductQuery, c command.ProductCommand) *ProductUsecase {
	return &ProductUsecase{Query: q, Command: c}
}
