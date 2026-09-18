package orderItem

import "errors"

var (
	// ErrOrderRequired 订单行项必须挂在一张订单上。
	ErrOrderRequired = errors.New("order is required")
	// ErrProductRequired 订单行项必须指向一件商品。
	ErrProductRequired = errors.New("product is required")
	// ErrProductNameRequired 商品名快照不能为空。
	ErrProductNameRequired = errors.New("product name snapshot is required")
	// ErrInvalidUnitPrice 单价快照必须大于 0（单位：分）。
	ErrInvalidUnitPrice = errors.New("unit price must be greater than zero")
	// ErrInvalidQuantity 数量必须大于 0。
	ErrInvalidQuantity = errors.New("quantity must be positive")

	// ErrQuantityBelowRefunded 新数量不能小于已退款数量。
	ErrQuantityBelowRefunded = errors.New("quantity can not be less than refunded quantity")
	// ErrRefundQuantityExceeded 退款数量超过了可退数量（已购 - 已退）。
	ErrRefundQuantityExceeded = errors.New("refund quantity exceeds refundable quantity")
	// ErrRestoreQuantityExceeded 还原的退款数量超过了已登记退款数量。
	ErrRestoreQuantityExceeded = errors.New("restore quantity exceeds refunded quantity")

	// ErrOrderItemRequired 传入的订单行项为空。
	ErrOrderItemRequired = errors.New("order item is required")
	// ErrOrderItemNotFound 指定的订单行项不存在。
	ErrOrderItemNotFound = errors.New("order item not found")
)
