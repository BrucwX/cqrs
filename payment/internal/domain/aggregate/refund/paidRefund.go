package refund

import "time"

// --- 已打款（终态） ---

// PaidRefund 已打款状态：退款额已经回到客户手里，终态。
//
// 名字里的 Paid 指「钱已经从商家出去」—— 与订单那边的 PaidOrder（钱进来）对称，
// 这样两个聚合的终态都叫 Paid，读代码时不用猜「Refunded 之后还能干嘛」。
type PaidRefund struct {
	RefundMeta
	reviewedAt time.Time // 审批时间
	paidAt     time.Time // 打款完成时间
}

// ReconstitutePaid 仓储恢复：字段全部由外部给定，不做校验。
func ReconstitutePaid(meta RefundMeta, reviewedAt time.Time, paidAt time.Time) PaidRefund {
	return PaidRefund{
		RefundMeta: meta,
		reviewedAt: reviewedAt,
		paidAt:     paidAt,
	}
}

func (r PaidRefund) ReviewedAt() time.Time { return r.reviewedAt }
func (r PaidRefund) PaidAt() time.Time     { return r.paidAt }
