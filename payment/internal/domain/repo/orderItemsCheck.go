package repo

import (
	"context"

	"payment/internal/domain/aggregate/customer"
	"payment/internal/domain/aggregate/orderItem"
	"payment/internal/domain/aggregate/product"
)

// OrderItemsCheckRepo 是一批行项预检服务（domain/service/orderItemsCheck）**自己要用的那几个方法**。
type OrderItemsCheckRepo interface {
	MustGetOrderItems(ctx context.Context, orderItemIDs []int64) (orderItem.OrderItems, error)
	MustGetProduct(ctx context.Context, productID int64) (product.Product, error)
	MustGetCustomer(ctx context.Context, customerID int64) (customer.Customer, error)
}
