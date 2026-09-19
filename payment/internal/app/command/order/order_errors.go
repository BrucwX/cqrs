package order

import "errors"

var (
	// ErrOrderItemsRequired 下单必须带上至少一行行项。
	ErrOrderItemsRequired = errors.New("order items are required")
	// ErrOrderNotCreatable 这批行项的预检没过：有的不是草稿 / 库存不够 / 受赠人不能收。
	//
	// 想知道具体是哪一条，看 domain/service/orderItemsCheck 里各 CheckXxx 的返回值。
	ErrOrderNotCreatable = errors.New("order items cannot be used to create an order")
)
