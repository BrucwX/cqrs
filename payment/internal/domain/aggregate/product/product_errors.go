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

	// ErrProductTargetRequired 具体商品必须指向一个目标（如课程 / 服务包的 id）。
	ErrProductTargetRequired = errors.New("product target is required")
	// ErrUnsupportedProductType 还不支持的**商品类型**：这类商品怎么校验、买完要做什么还没定义。
	ErrUnsupportedProductType = errors.New("unsupported product type")

	// --- 仓库层用（等 MustGetProduct 在适配层落地就会用到） ---

	// ErrProductRequired 传入的商品为空。
	ErrProductRequired = errors.New("product is required")
	// ErrProductNotFound 指定的商品不存在。
	ErrProductNotFound = errors.New("product not found")
)
