package repo

import (
	"context"

	"cqrs/internal/support/commerce/domain/aggregate/order"
)

// OrderRepo defines the order repository interface.
type OrderRepo interface {
	FindByID(ctx context.Context, id string) (*order.Order, error)
	FindByStudentID(ctx context.Context, studentID string) ([]*order.Order, error)
	CreateOrder(ctx context.Context, o *order.Order) (*order.Order, error)
	UpdateOrder(ctx context.Context, o *order.Order) (*order.Order, error)
	DeleteOrder(ctx context.Context, id string) error
}
