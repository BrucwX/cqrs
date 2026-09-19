package command

import (
	"context"
	"payment/internal/domain/aggregate/order"
)

type OrderCommand interface {
	MustGetOrder(ctx context.Context, orderID int64) (order.Order, error)
}
