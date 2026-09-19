package order

import "time"

// --- 已取消 ---

// CanceledOrder 取消状态（客户或后台主动取消）。
type CanceledOrder struct {
	OrderMeta
	canceledAt time.Time
	reason     string
}

// ReconstituteCanceled 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteCanceled(meta OrderMeta, canceledAt time.Time, reason string) CanceledOrder {
	return CanceledOrder{
		OrderMeta:  meta,
		canceledAt: canceledAt,
		reason:     reason,
	}
}

func (o CanceledOrder) CanceledAt() time.Time { return o.canceledAt }
func (o CanceledOrder) Reason() string        { return o.reason }
