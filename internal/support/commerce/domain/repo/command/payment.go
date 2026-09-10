package command

import "cqrs/internal/support/commerce/domain/aggregate/payment"

// PaymentCommand 支付命令接口
type PaymentCommand interface {
	// Save 保存支付（新增或更新）
	Save(p *payment.Payment) error
	// Delete 删除支付
	Delete(id string) error
}
