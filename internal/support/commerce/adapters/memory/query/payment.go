package query

import (
	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/payment"
)

type paymentQuery struct {
	data *memory.Data
}

// NewPaymentQuery creates a new PaymentQuery instance.
func NewPaymentQuery(d *memory.Data) *paymentQuery {
	return &paymentQuery{data: d}
}

func (q *paymentQuery) GetByID(id string) (*payment.Payment, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *paymentQuery) List() ([]*payment.Payment, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *paymentQuery) ListByOrderID(orderID string) ([]*payment.Payment, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *paymentQuery) ListByStudentID(studentID string) ([]*payment.Payment, error) {
	// TODO: implement database query
	return nil, nil
}
