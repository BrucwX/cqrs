package autoClosePay

import (
	"context"

	"github.com/go-kratos/kratos/v3/log"
)

// ListenBatch 每次唤醒最多从队列取几条 —— 单协程一条条办，取多了在通道里排队。
// 要调就改这个常量（暂不做成配置项）。
const ListenBatch = 10

// Start 等信号：队列唤醒一条就办一条，直到 ctx 结束。
//
// 唤醒模式（不再自己掐表轮询）：等待是队列实现的事；这边拿到信号就办，不空转。
//
// 单条失败**不**让循环停：错误记日志、也不 Ack —— 队列是「至少一次」，之后会重新投递，
// 重复消费交给 CloseIfUnpaid 的幂等兜住。
func (h *closePayScheduler) Start(ctx context.Context) error {
	signals, err := h.queue.Listen(ctx, ListenBatch)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case orderID, ok := <-signals:
			if !ok {
				return nil // 通道关了（Listen 的 ctx 结束）
			}
			h.handle(ctx, orderID)
		}
	}
}

// handle 办一条被唤醒的信号：办成了才 Ack。
func (h *closePayScheduler) handle(ctx context.Context, orderID int64) {
	if err := h.CloseIfUnpaid(ctx, orderID); err != nil {
		log.Error("payment autoClosePay: close unpaid order", "orderID", orderID, "err", err)
		return // 不 Ack：等队列重新投递
	}
	if err := h.queue.Ack(ctx, orderID); err != nil {
		// Ack 失败只会导致重复消费，有幂等兜着，不是致命错误。
		log.Error("payment autoClosePay: ack signal", "orderID", orderID, "err", err)
	}
}
