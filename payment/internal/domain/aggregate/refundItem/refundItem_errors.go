package refundItem

import "errors"

var (
	// ErrRefundRequired 退款行项必须挂在一张退款单上。
	ErrRefundRequired = errors.New("refund is required")
	// ErrOrderItemRequired 退款行项必须指向一个订单行项。
	ErrOrderItemRequired = errors.New("order item is required")
	// ErrInvalidQuantity 退款件数必须大于 0。
	ErrInvalidQuantity = errors.New("refund quantity must be positive")
	// ErrInvalidAmount 退款金额必须大于 0（单位：分）。
	ErrInvalidAmount = errors.New("refund item amount must be greater than zero")

	// ErrRefundItemRequired 传入的退款行项为空。
	ErrRefundItemRequired = errors.New("refund item is required")
	// ErrRefundItemNotFound 指定的退款行项不存在。
	ErrRefundItemNotFound = errors.New("refund item not found")
)
