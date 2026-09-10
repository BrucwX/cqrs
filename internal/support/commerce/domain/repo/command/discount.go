package command

import "cqrs/internal/support/commerce/domain/aggregate/discount"

// DiscountCommand 折扣命令接口
type DiscountCommand interface {
	// Save 保存折扣（新增或更新）
	Save(d *discount.Discount) error
	// Delete 删除折扣
	Delete(id string) error
}
