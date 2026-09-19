package product

import "errors"

var (
	// ErrProductNameRequired 商品必须有名字。
	ErrProductNameRequired = errors.New("product name is required")
	// ErrInvalidPrice 售价必须大于 0（单位：分）。
	ErrInvalidPrice = errors.New("product price must be greater than zero")
	// ErrInvalidProductType 商品分类不合法（零值 Type，或 name / code 缺项）。
	ErrInvalidProductType = errors.New("product type is invalid")
	// ErrInvalidStock 库存数量不合法：有限库存的剩余不能为负；出库 / 退单 / 进货的件数
	// 必须大于 0，退单不能超过累计销售量。
	ErrInvalidStock = errors.New("product stock quantity is invalid")
	// ErrInsufficientStock 库存不足，卖不出这么多。
	ErrInsufficientStock = errors.New("insufficient product stock")

	// ErrProductAlreadyOnSale 商品已经在售了。
	ErrProductAlreadyOnSale = errors.New("product is already on sale")
	// ErrProductAlreadyOffShelf 商品已经下架了。
	ErrProductAlreadyOffShelf = errors.New("product is already off shelf")

	// ErrProductRequired 传入的商品为空。
	ErrProductRequired = errors.New("product is required")
	// ErrProductNotFound 指定的商品不存在。
	ErrProductNotFound = errors.New("product not found")
)
