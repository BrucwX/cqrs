package product

import "errors"

var (
	// ErrProductNameRequired 商品必须有名字。
	ErrProductNameRequired = errors.New("product name is required")
	// ErrInvalidPrice 售价必须大于 0（单位：分）。
	ErrInvalidPrice = errors.New("product price must be greater than zero")

	// ErrProductAlreadyOnSale 商品已经在售了。
	ErrProductAlreadyOnSale = errors.New("product is already on sale")
	// ErrProductAlreadyOffShelf 商品已经下架了。
	ErrProductAlreadyOffShelf = errors.New("product is already off shelf")

	// ErrProductRequired 传入的商品为空。
	ErrProductRequired = errors.New("product is required")
	// ErrProductNotFound 指定的商品不存在。
	ErrProductNotFound = errors.New("product not found")
)
