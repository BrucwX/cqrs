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
}
