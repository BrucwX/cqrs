package command

import (
	"context"
	"payment/internal/domain/aggregate/orderItem"
)

type OrderItemCommand interface {
	MustGetOrderItem(ctx context.Context, orderItemID int64) (orderItem.OrderItem, error)
	MustGetOrderItems(ctx context.Context, orderItemIDs []int64) (orderItem.OrderItems, error)
}
