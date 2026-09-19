package order

import "context"

type OrderCommand interface {
	MustGetOrder(ctx context.Context, orderID int64) (Order, error)
	// Save 落库订单（新建或覆盖，由实现按主键判）。
	//
	// 参数是接口 `Order`：待支付 / 已支付 / 已取消 / 已关单都能存，
	// 适配器用 type switch 取状态名写进 status 列。
	Save(ctx context.Context, o Order) error
}
