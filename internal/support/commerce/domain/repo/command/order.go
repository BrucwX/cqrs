package command

import "cqrs/internal/support/commerce/domain/aggregate/order"

// OrderCommand 订单命令接口
type OrderCommand interface {
	// Save 保存订单（新增或更新）
	Save(o *order.Order) error
	// Delete 删除订单
	Delete(id string) error
}
