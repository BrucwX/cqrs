package autoClosePay

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v3/log"
)

const (
	// defaultPollInterval 轮询间隔：信号到点后最多晚这么久被处理。
	defaultPollInterval = 30 * time.Second
	// defaultBatch 一轮最多处理多少条信号 —— 欠账多时分批还，别把一轮拖太久。
	defaultBatch = 100
)

// Start 定时轮询：到点就把信号取出来消费，直到 ctx 结束。
//
// 单条失败**不**让循环停：错误记日志、也不 Ack，下一轮还会取到它 —— 队列是「至少一次」，
// 重复消费交给 CloseIfUnpaid 的幂等兜住。
func (h *closePayScheduler) Start(ctx context.Context) error {
	tk := time.NewTicker(h.interval)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tk.C:
			h.pollOnce(ctx)
		}
	}
}

// pollOnce 跑一轮：取一批到点的信号，逐条消费；**只有处理成功的才 Ack**。
func (h *closePayScheduler) pollOnce(ctx context.Context) {
	signals, err := h.queue.Due(ctx, h.batch)
	if err != nil {
		log.Error("payment autoClosePay: read due signals", "err", err)
		return
	}

	for _, orderID := range signals {
		if err := h.CloseIfUnpaid(ctx, orderID); err != nil {
			log.Error("payment autoClosePay: close unpaid order", "orderID", orderID, "err", err)
			continue // 不 Ack：下一轮还会取到，再来一次
		}
		if err := h.queue.Ack(ctx, orderID); err != nil {
			// Ack 失败只会导致重复消费，有幂等兜着，不是致命错误。
			log.Error("payment autoClosePay: ack signal", "orderID", orderID, "err", err)
		}
	}
}
