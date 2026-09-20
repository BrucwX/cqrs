package repo

import (
	"context"

	"payment/internal/domain/aggregate/product"
)

// ProductCommand 商品的写侧仓储。
type ProductCommand interface {
	MustGetProduct(ctx context.Context, productID int64) (product.Product, error)
	// Update 落库商品（扣完库存 / 改完状态之后存）。
	//
	// 参数是接口 `product.Product`：三种商品都能存，适配器按具体类型取状态与字段。
	Update(ctx context.Context, p product.Product) error
}
