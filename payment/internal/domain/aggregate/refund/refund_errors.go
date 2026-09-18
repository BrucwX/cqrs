package refund

import "errors"

var (
	// --- 退款单 (Refund) 校验 ---

	// ErrRefundOrderRequired 退款必须挂在一张订单上。
	ErrRefundOrderRequired = errors.New("refund order is required")
	// ErrRefundBuyerRequired 退款必须有申请人。
	ErrRefundBuyerRequired = errors.New("refund buyer is required")
	// ErrInvalidAmount 退款金额必须大于 0（单位：分）。
	ErrInvalidAmount = errors.New("refund amount must be greater than zero")
	// ErrRefundReasonRequired 申请退款必须写明原因。
	ErrRefundReasonRequired = errors.New("refund reason is required")
	// ErrRejectReasonRequired 拒绝退款必须写明原因。
	ErrRejectReasonRequired = errors.New("refund reject reason is required")

	// --- 状态流转 ---

	// ErrRefundNotApplied 只有「已申请」的退款单能做这次变更。
	ErrRefundNotApplied = errors.New("refund is not in applied status")
	// ErrRefundNotApproved 只有「审批通过」的退款单能打款。
	ErrRefundNotApproved = errors.New("refund is not approved")

	// --- 仓库层用 ---

	// ErrRefundRequired 传入的退款单为空。
	ErrRefundRequired = errors.New("refund is required")
	// ErrRefundNotFound 指定的退款单不存在。
	ErrRefundNotFound = errors.New("refund not found")
)
