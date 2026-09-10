package repo

import (
	"context"

	"cqrs/internal/support/commerce/domain/aggregate/payment"
)

// PaymentRepo defines the payment repository interface.
type PaymentRepo interface {
	FindByID(ctx context.Context, id string) (*payment.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*payment.Payment, error)
	CreatePayment(ctx context.Context, p *payment.Payment) (*payment.Payment, error)
	UpdatePayment(ctx context.Context, p *payment.Payment) (*payment.Payment, error)
	DeletePayment(ctx context.Context, id string) error
}
