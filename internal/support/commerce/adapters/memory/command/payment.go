package command

import (
	"context"

	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/payment"
	"cqrs/internal/support/commerce/domain/repo"
)

type paymentRepo struct {
	data *memory.Data
}

// NewPaymentRepo creates a new PaymentRepo instance.
func NewPaymentRepo(d *memory.Data) repo.PaymentRepo {
	return &paymentRepo{data: d}
}

func (r *paymentRepo) FindByID(ctx context.Context, id string) (*payment.Payment, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *paymentRepo) FindByOrderID(ctx context.Context, orderID string) (*payment.Payment, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *paymentRepo) CreatePayment(ctx context.Context, p *payment.Payment) (*payment.Payment, error) {
	// TODO: implement database insert
	return p, nil
}

func (r *paymentRepo) UpdatePayment(ctx context.Context, p *payment.Payment) (*payment.Payment, error) {
	// TODO: implement database update
	return p, nil
}

func (r *paymentRepo) DeletePayment(ctx context.Context, id string) error {
	// TODO: implement database soft delete
	return nil
}
