package pay

import (
	"context"
	"time"

	"payment/internal/domain/aggregate/order"
)

// ConfirmIn 渠道回调解析后的**统一入参**：adapter 已经验过签、把渠道报文翻译完，
// 领域 / 用例不碰原始报文（XML / 表单 / 签名）。
type ConfirmIn struct {
	OrderID    int64
	Channel    order.PayChannel
	TradeNo    string    // 渠道流水号，对账 / 幂等靠它
	PaidAmount int64     // 实付金额（分）
	PaidAt     time.Time // 渠道侧的支付时间
}

// Confirm 异步回调：确认一笔支付成功（幂等）。
//
// 渠道会重试回调、我方也可能重复消费，三种情形都要能扛：
//   - 订单还是待支付 → MarkPaid 落库；
//   - 已经支付过 → 直接返回 nil（幂等 ACK，让渠道别再重试）；
//   - 已关单 / 已取消却收到「支付成功」→ 钱收了、单没了，属于资金差异，
//     返回 ErrUnexpectedNotify 交给人工 / 对账，不能默默吞掉。
//
// ⚠️ 实付金额对不对由 MarkPaid 把关（必须等于应付），不在这里单独比对。
func (h *Handler) Confirm(ctx context.Context, in ConfirmIn) error {
	o, err := h.orderCommand.MustGetOrder(ctx, in.OrderID)
	if err != nil {
		return err
	}

	switch o := o.(type) {
	case order.PendingOrder:
		paid, err := o.MarkPaid(in.PaidAmount, in.Channel, in.TradeNo, in.PaidAt)
		if err != nil {
			return err
		}
		return h.orderCommand.Save(ctx, paid)

	case order.PaidOrder:
		// 重复回调：已经付过，幂等成功。
		return nil

	default:
		// 已关单 / 已取消却收到「支付成功」。
		return ErrUnexpectedNotify
	}
}
