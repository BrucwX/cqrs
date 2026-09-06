package order

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/shared/types"
)

// OrderRepo defines the order repository interface.
type OrderRepo interface {
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByStudentID(ctx context.Context, studentID string) ([]*Order, error)
	ListOrders(ctx context.Context, opts ...types.ListOption) ([]*Order, error)
	CreateOrder(ctx context.Context, o *Order) (*Order, error)
	UpdateOrder(ctx context.Context, o *Order) (*Order, error)
	DeleteOrder(ctx context.Context, id string) error
}
