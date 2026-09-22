// Package closePayScheduler 是「超时关单」用例（写侧）：定时器到点后调它，
// 把过了支付窗口还没付的订单关掉 —— 先本地关单，再让渠道侧也关一下。
package course

import (
	"context"
	"time"

	orderagg "payment/internal/domain/aggregate/order"
	"payment/internal/domain/repo"
	"payment/internal/ports"
)

type closePayScheduler struct {
	orderCommand repo.OrderCommand
	pay          ports.PayService
	scheduler    ports.ClosePayProducer
}

func (h *closePayScheduler) closeIfUnpaid(ctx context.Context, orderID int64) error {
	o, err := h.orderCommand.MustGetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	pending, ok := o.(orderagg.PendingOrder)
	if !ok {
		return nil // 付了 / 取消 / 已关单，没得可关
	}

	now := time.Now()
	if !pending.IsExpired(now) {
		return nil // 定时器提前触发，还没到点
	}

	closed, err := pending.CloseIfExpired(now)
	if err != nil {
		return err
	}
	if err := h.orderCommand.Save(ctx, closed); err != nil {
		return err
	}

	// 渠道侧也关一下（微信 closeorder；支付宝自动关 / 无副作用）。
	return h.pay.Close(ctx, orderID)
}

func (h *closePayScheduler) Start(ctx context.Context) error {

	return nil
}
