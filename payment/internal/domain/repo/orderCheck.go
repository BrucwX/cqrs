package repo

import (
	"context"

	"payment/internal/domain/aggregate/customer"
	"payment/internal/domain/aggregate/order"
)

// OrderCheckRepo 是订单侧判定服务（domain/service/orderCheck）**自己要用的那几个方法**：
// 服务需要什么就声明什么，不去背「完整仓储」的包袱。
//
// 实现侧（adapters 的载体）只要签名对得上就能注入。
type OrderCheckRepo interface {
	MustGetCustomer(ctx context.Context, customerID int64) (customer.Customer, error)
	MustGetOrder(ctx context.Context, orderID int64) (order.Order, error)
}
