package order

import "time"

// --- 超时关闭 ---

// ClosedOrder 超时未支付被关单的状态（由定时任务调 PendingOrder.CloseIfExpired 得到）。
type ClosedOrder struct {
	OrderMeta
	closedAt time.Time
}

// ReconstituteClosed 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteClosed(meta OrderMeta, closedAt time.Time) ClosedOrder {
	return ClosedOrder{
		OrderMeta: meta,
		closedAt:  closedAt,
	}
}

func (o ClosedOrder) ClosedAt() time.Time { return o.closedAt }
