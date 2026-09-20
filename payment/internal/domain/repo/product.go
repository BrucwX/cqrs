package repo

import (
	"context"

	"payment/internal/domain/aggregate/orderItem"
	"payment/internal/domain/aggregate/product"
)

// ProductCommand 商品的写侧仓储。
type ProductCommand interface {
	MustGetProduct(ctx context.Context, productID int64) (product.Product, error)

	GetProsByItems(ctx context.Context, oitems orderItem.OrderItems) (product.Products, error)
}
