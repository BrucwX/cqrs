package product

import "context"

type ProductCommand interface {
	MustGetProduct(ctx context.Context, productID int64) (Product, error)
	// Update 落库商品（扣完库存 / 改完状态之后存）。
	//
	// 参数是接口 `Product`：三种商品都能存，适配器按具体类型取状态与字段。
	Update(ctx context.Context, p Product) error
}
