package repo

import (
	"context"

	"payment/internal/domain/aggregate/orderItem"
)

// OrderItemCommand 订单行项的写侧仓储。
type OrderItemCommand interface {
	MustGetOrderItem(ctx context.Context, orderItemID int64) (orderItem.OrderItem, error)
	MustGetOrderItems(ctx context.Context, orderItemIDs []int64) (orderItem.OrderItems, error)
	// Save 批量落库：把这批行项整批写回去（新建或覆盖，由实现按主键判）。
	//
	// ⚠️ 名字里**不带状态**：存储只管把给它的行项写下去，「草稿 → 正式行」是领域的事
	// （OrderItems.ConfirmAll），别让仓储接口去表达状态迁移。
	Save(ctx context.Context, items orderItem.OrderItems) error

	// ListByOrder 查一张订单下的所有正式行（草稿行不属于任何订单，自然查不到）。
	//
	// 超时关单要照着它们退课：订单聚合只认钱，**买了什么、买给谁**都长在行项上，
	// 而「从订单号反查行项」领域里没有别的路。
	ListByOrder(ctx context.Context, orderID int64) (orderItem.OrderItems, error)
}
