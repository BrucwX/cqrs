package repo

import (
	"context"

	"payment/internal/domain/aggregate/orderItem"
)

// OrderItemsCheckRepo 是一批行项预检服务（domain/service/orderItemsCheck）**自己要用的那几个方法**。
//
// 服务只用它取那批行项，别的都不要 —— 服务不背「完整仓储」的包袱；
// 以后服务真的要看商品 / 客户，再往这个接口上加方法。
type OrderItemsCheckRepo interface {
	MustGetOrderItems(ctx context.Context, orderItemIDs []int64) (orderItem.OrderItems, error)
}
