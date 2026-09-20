package repo

import (
	"context"

	"payment/internal/domain/aggregate/order"
)

// OrderCommand 订单的写侧仓储。
type OrderCommand interface {
	MustGetOrder(ctx context.Context, orderID int64) (order.Order, error)
	// Save 落库订单（新建或覆盖，由实现按主键判）。
	//
	// 参数是接口 `order.Order`：待支付 / 已支付 / 已取消 / 已关单都能存，
	// 适配器用 type switch 取状态名写进 status 列。
	Save(ctx context.Context, o order.Order) error
}
