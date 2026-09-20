package pay

import "payment/internal/domain/repo"

// Handler 是「支付」用例的处理器：下单之后的支付确认都走这里。
//
// 依赖只有订单写侧仓储：支付渠道的差异（验签、报文、回执）在 adapters / 渠道包里，
// 本用例只做领域动作 —— 把「这笔钱到账了」落到订单上（MarkPaid + 保存）。
type Handler struct {
	orderCommand repo.OrderCommand
}

// NewHandler 构造：依赖由上层（wire）注入。
func NewHandler(orderCommand repo.OrderCommand) *Handler {
	return &Handler{orderCommand: orderCommand}
}
