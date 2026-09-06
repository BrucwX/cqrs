package payment

import (
	"context"

	"cqrs/internal/shared/types"
)

// PaymentRepo defines the payment repository interface.
type PaymentRepo interface {
	FindByID(ctx context.Context, id string) (*Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*Payment, error)
	ListPayments(ctx context.Context, opts ...types.ListOption) ([]*Payment, error)
	CreatePayment(ctx context.Context, p *Payment) (*Payment, error)
	UpdatePayment(ctx context.Context, p *Payment) (*Payment, error)
	DeletePayment(ctx context.Context, id string) error
}
