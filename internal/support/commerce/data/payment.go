package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/support/commerce/biz/payment"
	"github.com/go-kratos/kratos-layout/internal/shared/types"
)

type paymentRepo struct {
	data *Data
}

// NewPaymentRepo creates a new PaymentRepo instance.
func NewPaymentRepo(d *Data) payment.PaymentRepo {
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

func (r *paymentRepo) ListPayments(ctx context.Context, opts ...types.ListOption) ([]*payment.Payment, error) {
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
