package refund

import (
	"strings"
	"time"

	"payment/pkg/uuid"
)

// --- 已申请（待审批） ---

// AppliedRefund 已申请状态：客户提了退款，还等着审批。
//
// 这是唯一「还能改金额」的阶段（部分同意 —— 审批时把金额砍掉一部分），
// 所以 ChangeAmount 返回更新后的退款单；一旦 Approve / Reject / Cancel 出去，就回不来了。
type AppliedRefund struct {
	RefundMeta // 自动获得 ID()、OrderID()、CustomerID()、Meta() 与 isRefundState()
}

// NewRefund 客户发起退款申请，落在「已申请」。
//
// 退款额本身对不对（不能超过订单实付、退几件要跟订单行项对得上）是跨聚合的判定，
// 本聚合只认「金额为正、必须写明原因」，其余交给 domain/service 或 app 层先判完再来写。
func NewRefund(
	orderID int64,
	customerID int64,
	amount int64,
	reason string,
	now time.Time,
) (AppliedRefund, error) {
	if orderID <= 0 {
		return AppliedRefund{}, ErrRefundOrderRequired
	}
	if customerID <= 0 {
		return AppliedRefund{}, ErrRefundCustomerRequired
	}
	if amount <= 0 {
		return AppliedRefund{}, ErrInvalidAmount
	}
	if strings.TrimSpace(reason) == "" {
		return AppliedRefund{}, ErrRefundReasonRequired
	}

	return AppliedRefund{
		RefundMeta: RefundMeta{
			id:         uuid.NewInt64(),
			orderID:    orderID,
			customerID: customerID,
			amount:     amount,
			reason:     reason,
			appliedAt:  now,
		},
	}, nil
}

// ReconstituteApplied 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteApplied(meta RefundMeta) AppliedRefund {
	return AppliedRefund{RefundMeta: meta}
}

// ChangeAmount 审批时调整退款额（部分同意），返回更新后的退款单。
func (r AppliedRefund) ChangeAmount(amount int64) (AppliedRefund, error) {
	if amount <= 0 {
		return r, ErrInvalidAmount
	}
	r.amount = amount
	return r, nil
}

// Approve 审批通过 → ApprovedRefund（等着打款）。
func (r AppliedRefund) Approve(now time.Time) ApprovedRefund {
	return ApprovedRefund{
		RefundMeta: r.RefundMeta,
		reviewedAt: now,
	}
}

// Reject 审批拒绝 → RejectedRefund，终态，必须写明原因。
func (r AppliedRefund) Reject(reason string, now time.Time) (RejectedRefund, error) {
	if strings.TrimSpace(reason) == "" {
		return RejectedRefund{}, ErrRejectReasonRequired
	}
	return RejectedRefund{
		RefundMeta:   r.RefundMeta,
		reviewedAt:   now,
		rejectReason: reason,
	}, nil
}

// Cancel 客户在审批前自己撤销 → CanceledRefund，终态。
func (r AppliedRefund) Cancel(now time.Time) CanceledRefund {
	return CanceledRefund{
		RefundMeta: r.RefundMeta,
		canceledAt: now,
	}
}
