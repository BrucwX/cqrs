package query

import (
	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/order"
)

type orderQuery struct {
	data *memory.Data
}

// NewOrderQuery creates a new OrderQuery instance.
func NewOrderQuery(d *memory.Data) *orderQuery {
	return &orderQuery{data: d}
}

func (q *orderQuery) GetByID(id string) (*order.Order, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *orderQuery) List() ([]*order.Order, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *orderQuery) ListByStudentID(studentID string) ([]*order.Order, error) {
	// TODO: implement database query
	return nil, nil
}
