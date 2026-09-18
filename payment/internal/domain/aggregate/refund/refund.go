package refund

import (
	"strings"
	"time"

	"payment/pkg/uuid"
)

// --- 聚合根 (Aggregate Root) ---

// Refund 是「一笔退款申请」的聚合根。
//
// 状态机（只允许下面这几条边）：
//
//	Applied ──Approve──► Approved ──MarkRefunded──► Refunded
//	   ├──Reject──► Rejected
//	   └──Cancel──► Canceled
//
// 金额一律用「分」的 int64。退款额本身对不对（不能超过订单实付、退几件要跟
// 订单行项对得上）是跨聚合的判定，本聚合只认「金额为正」，其余交给 domain/service
// 或 app 层先判完再来写。
type Refund struct {
	id      int64
	orderID int64
	buyerID int64
	amount  int64  // 申请退款金额（分）
	reason  string // 申请原因

	status Status

	appliedAt    time.Time
	reviewedAt   time.Time // 审批（同意 / 拒绝）时间
	refundedAt   time.Time // 打款完成时间
	rejectReason string
}

// NewRefund 买家发起退款申请，落在 StatusApplied。
func NewRefund(
	orderID int64,
	buyerID int64,
	amount int64,
	reason string,
	now time.Time,
) (*Refund, error) {
	if orderID <= 0 {
		return nil, ErrRefundOrderRequired
	}
	if buyerID <= 0 {
		return nil, ErrRefundBuyerRequired
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if strings.TrimSpace(reason) == "" {
		return nil, ErrRefundReasonRequired
	}

	return &Refund{
		id:        uuid.NewInt64(),
		orderID:   orderID,
		buyerID:   buyerID,
		amount:    amount,
		reason:    reason,
		status:    StatusApplied,
		appliedAt: now,
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(
	id int64,
	orderID int64,
	buyerID int64,
	amount int64,
	reason string,
	status Status,
	appliedAt time.Time,
	reviewedAt time.Time,
	refundedAt time.Time,
	rejectReason string,
) *Refund {
	return &Refund{
		id:           id,
		orderID:      orderID,
		buyerID:      buyerID,
		amount:       amount,
		reason:       reason,
		status:       status,
		appliedAt:    appliedAt,
		reviewedAt:   reviewedAt,
		refundedAt:   refundedAt,
		rejectReason: rejectReason,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// Approve 审批通过：进入待打款。
func (r *Refund) Approve(now time.Time) error {
	if r.status != StatusApplied {
		return ErrRefundNotApplied
	}
	r.status = StatusApproved
	r.reviewedAt = now
	return nil
}

// Reject 审批拒绝，必须写明原因。
func (r *Refund) Reject(reason string, now time.Time) error {
	if r.status != StatusApplied {
		return ErrRefundNotApplied
	}
	if strings.TrimSpace(reason) == "" {
		return ErrRejectReasonRequired
	}
	r.status = StatusRejected
	r.rejectReason = reason
	r.reviewedAt = now
	return nil
}

// MarkRefunded 打款完成（终态）。
func (r *Refund) MarkRefunded(now time.Time) error {
	if r.status != StatusApproved {
		return ErrRefundNotApproved
	}
	r.status = StatusRefunded
	r.refundedAt = now
	return nil
}

// Cancel 买家在审批前自己撤销。
func (r *Refund) Cancel(now time.Time) error {
	if r.status != StatusApplied {
		return ErrRefundNotApplied
	}
	r.status = StatusCanceled
	r.reviewedAt = now
	return nil
}

// ChangeAmount 审批时调整退款额（部分同意），只在待审批阶段可改。
func (r *Refund) ChangeAmount(amount int64) error {
	if r.status != StatusApplied {
		return ErrRefundNotApplied
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	r.amount = amount
	return nil
}

// IsOpen 是否还在处理中（钱还没退出去，也没被拒）。
func (r *Refund) IsOpen() bool {
	return r.status == StatusApplied || r.status == StatusApproved
}

// --- 属性只读访问器 (Getters) ---

func (r *Refund) ID() int64             { return r.id }
func (r *Refund) OrderID() int64        { return r.orderID }
func (r *Refund) BuyerID() int64        { return r.buyerID }
func (r *Refund) Amount() int64         { return r.amount }
func (r *Refund) Reason() string        { return r.reason }
func (r *Refund) Status() Status        { return r.status }
func (r *Refund) AppliedAt() time.Time  { return r.appliedAt }
func (r *Refund) ReviewedAt() time.Time { return r.reviewedAt }
func (r *Refund) RefundedAt() time.Time { return r.refundedAt }
func (r *Refund) RejectReason() string  { return r.rejectReason }
