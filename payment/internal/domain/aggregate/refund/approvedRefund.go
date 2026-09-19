package refund

import "time"

// --- 审批通过（待打款） ---

// ApprovedRefund 审批通过状态：同意退了，钱还没打出去。
type ApprovedRefund struct {
	RefundMeta
	reviewedAt time.Time // 状态专属字段：审批时间
}

// ReconstituteApproved 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteApproved(meta RefundMeta, reviewedAt time.Time) ApprovedRefund {
	return ApprovedRefund{
		RefundMeta: meta,
		reviewedAt: reviewedAt,
	}
}

func (r ApprovedRefund) ReviewedAt() time.Time { return r.reviewedAt }

// MarkPaid 打款完成 → PaidRefund，终态。
//
// 钱具体怎么退回去（原路退回 / 人工转账 / 退到余额）是支付渠道的事，
// 本聚合只记「钱已经出去了、什么时候出去的」。
func (r ApprovedRefund) MarkPaid(now time.Time) PaidRefund {
	return PaidRefund{
		RefundMeta: r.RefundMeta,
		reviewedAt: r.reviewedAt,
		paidAt:     now,
	}
}
