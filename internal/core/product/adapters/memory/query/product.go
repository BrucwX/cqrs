package query

import (
	"cqrs/internal/core/product/adapters/memory"
	"cqrs/internal/core/product/domain/aggregate/product"
)

type productQuery struct {
	data *memory.Data
}

// NewProductQuery creates a new ProductQuery instance.
func NewProductQuery(d *memory.Data) *productQuery {
	return &productQuery{data: d}
}

func (q *productQuery) GetByID(id string) (*product.Product, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *productQuery) List() ([]*product.Product, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *productQuery) ListByCategory(category product.ProductCategory) ([]*product.Product, error) {
	// TODO: implement database query
	return nil, nil
}
