package order

import "errors"

var (
	// ErrOrderCustomerRequired 订单必须属于某个客户。
	ErrOrderCustomerRequired = errors.New("order customer is required")
	// ErrInvalidAmount 订单原价必须大于 0。
	ErrInvalidAmount = errors.New("order amount must be greater than zero")
	// ErrInvalidPaymentWindow 支付时限必须大于 0。
	ErrInvalidPaymentWindow = errors.New("order payment window must be positive")

	// ErrInvalidDiscount 单笔优惠额必须大于 0，且累计优惠不能超过原价。
	ErrInvalidDiscount = errors.New("order discount amount is invalid")

	// --- 状态迁移（能走到这里的都是「状态对了、条件不满足」） ---

	// ErrOrderNotPending 订单不是待支付状态（已经付了 / 取消 / 关单），没有可付的了。
	ErrOrderNotPending = errors.New("order is not pending")
	// ErrOrderExpired 订单已过支付时限，不能再支付（要先关单）。
	ErrOrderExpired = errors.New("order has expired")
	// ErrOrderNotExpired 订单还没到期，不能关单。
	ErrOrderNotExpired = errors.New("order has not expired yet")
	// ErrPaymentAmountMismatch 实付金额与应付金额不一致。
	ErrPaymentAmountMismatch = errors.New("paid amount does not match payable amount")
	// ErrInvalidPayChannel 支付途径不合法（零值 PayChannel，或 code / name 缺项）。
	ErrInvalidPayChannel = errors.New("pay channel is invalid")

	// --- 仓库层用 ---

	// ErrOrderRequired 传入的订单为空。
	ErrOrderRequired = errors.New("order is required")
	// ErrOrderNotFound 指定的订单不存在。
	ErrOrderNotFound = errors.New("order not found")

	// --- 下单用例用 ---

	// ErrOrderItemsRequired 下单必须带上至少一行行项。
	ErrOrderItemsRequired = errors.New("order items are required")
	// ErrOrderNotCreatable 这批行项建不出订单：有的不是草稿 / 库存不够 / 受赠人不能收。
	//
	// 想知道具体是哪一条，看 domain/service/orderItemsCheck 里各 CheckXxx 的返回值。
	ErrOrderNotCreatable = errors.New("order items cannot be used to create an order")
)
