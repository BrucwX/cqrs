package order

import "errors"

var (
	// ErrOrderBuyerRequired 订单必须属于某个买家。
	ErrOrderBuyerRequired = errors.New("order buyer is required")
	// ErrInvalidAmount 订单原价必须大于 0。
	ErrInvalidAmount = errors.New("order amount must be greater than zero")
	// ErrInvalidPaymentWindow 支付时限必须大于 0。
	ErrInvalidPaymentWindow = errors.New("order payment window must be positive")

	// ErrInvalidDiscount 单笔优惠额必须大于 0，且累计优惠不能超过原价。
	ErrInvalidDiscount = errors.New("order discount amount is invalid")

	// ErrOrderNotPending 订单不在待支付状态，不能做这次变更。
	ErrOrderNotPending = errors.New("order is not pending")
	// ErrOrderExpired 订单已过支付时限，需先关单。
	ErrOrderExpired = errors.New("order has expired")
	// ErrOrderNotExpired 订单还没到期，不能关单。
	ErrOrderNotExpired = errors.New("order has not expired yet")
	// ErrPaymentAmountMismatch 实付金额与应付金额不一致。
	ErrPaymentAmountMismatch = errors.New("paid amount does not match payable amount")
	// ErrOrderNotPaid 只有已支付订单才能退款。
	ErrOrderNotPaid = errors.New("order is not paid")

	// ErrOrderRequired 传入的订单为空。
	ErrOrderRequired = errors.New("order is required")
	// ErrOrderNotFound 指定的订单不存在。
	ErrOrderNotFound = errors.New("order not found")
)
