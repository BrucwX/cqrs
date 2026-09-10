package query

import (
	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/discount"
)

type discountQuery struct {
	data *memory.Data
}

// NewDiscountQuery creates a new DiscountQuery instance.
func NewDiscountQuery(d *memory.Data) *discountQuery {
	return &discountQuery{data: d}
}

func (q *discountQuery) GetByID(id string) (*discount.Discount, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *discountQuery) List() ([]*discount.Discount, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *discountQuery) GetByCode(code string) (*discount.Discount, error) {
	// TODO: implement database query
	return nil, nil
}
