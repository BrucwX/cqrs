package refund

import "time"

// --- 审批拒绝（终态） ---

// RejectedRefund 审批拒绝状态：不同意退，终态。
//
// 客户想再退只能重新发起一张新的退款单 —— 这张单已经定性了，不再改。
type RejectedRefund struct {
	RefundMeta
	reviewedAt   time.Time // 审批时间
	rejectReason string    // 拒绝原因（必填）
}

// ReconstituteRejected 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteRejected(
	meta RefundMeta,
	reviewedAt time.Time,
	rejectReason string,
) RejectedRefund {
	return RejectedRefund{
		RefundMeta:   meta,
		reviewedAt:   reviewedAt,
		rejectReason: rejectReason,
	}
}

func (r RejectedRefund) ReviewedAt() time.Time { return r.reviewedAt }
func (r RejectedRefund) RejectReason() string  { return r.rejectReason }
