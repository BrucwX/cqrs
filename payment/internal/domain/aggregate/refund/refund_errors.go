package refund

import "errors"

var (
	// --- 退款单 (Refund) 校验 ---

	// ErrRefundOrderRequired 退款必须挂在一张订单上。
	ErrRefundOrderRequired = errors.New("refund order is required")
	// ErrRefundCustomerRequired 退款必须有申请人。
	ErrRefundCustomerRequired = errors.New("refund customer is required")
	// ErrInvalidAmount 退款金额必须大于 0（单位：分）。
	ErrInvalidAmount = errors.New("refund amount must be greater than zero")
	// ErrRefundReasonRequired 申请退款必须写明原因。
	ErrRefundReasonRequired = errors.New("refund reason is required")
	// ErrRejectReasonRequired 拒绝退款必须写明原因。
	ErrRejectReasonRequired = errors.New("refund reject reason is required")

	// 注：「只有已申请的退款单能审批」「只有审批通过的才能打款」这类状态约束
	// 已经由状态类型本身在**编译期**保证（ApprovedRefund 才有 MarkPaid），
	// 所以这里不再有 ErrRefundNotApplied / ErrRefundNotApproved。

	// --- 仓库层用 ---

	// ErrRefundRequired 传入的退款单为空。
	ErrRefundRequired = errors.New("refund is required")
	// ErrRefundNotFound 指定的退款单不存在。
	ErrRefundNotFound = errors.New("refund not found")
)
