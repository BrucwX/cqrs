package refund

import "time"

// --- 已撤销（终态） ---

// CanceledRefund 客户自己撤销的退款单，终态。
//
// 注意撤销时间**不是**审批时间：这条单根本没走到审批，所以单独记 canceledAt，
// 不要塞进类似状态的 reviewedAt 里凑数。
type CanceledRefund struct {
	RefundMeta
	canceledAt time.Time
}

// ReconstituteCanceled 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteCanceled(meta RefundMeta, canceledAt time.Time) CanceledRefund {
	return CanceledRefund{
		RefundMeta: meta,
		canceledAt: canceledAt,
	}
}

func (r CanceledRefund) CanceledAt() time.Time { return r.canceledAt }
