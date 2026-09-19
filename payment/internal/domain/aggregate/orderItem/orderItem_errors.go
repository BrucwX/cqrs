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

	// ErrOrderItemRequired 传入的订单行项为空。
	ErrOrderItemRequired = errors.New("order item is required")
	// ErrOrderItemNotFound 指定的订单行项不存在。
	ErrOrderItemNotFound = errors.New("order item not found")
)
