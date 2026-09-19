package product

import "errors"

var (
	// ErrProductNameRequired 商品必须有名字。
	ErrProductNameRequired = errors.New("product name is required")
	// ErrInvalidPrice 售价必须大于 0（单位：分）。
	ErrInvalidPrice = errors.New("product price must be greater than zero")
	// ErrInvalidProductType 商品分类不合法（零值 Type，或 name / code 缺项）。
	ErrInvalidProductType = errors.New("product type is invalid")
	// ErrProductUnitRequired 有限类 / 无限类商品必须写明计价单位（件 / 千克 / GB…）。
	ErrProductUnitRequired = errors.New("product unit is required")
	// ErrInvalidStock 库存不合法：剩余不能为负；退单数量不能超过累计销售量。
	ErrInvalidStock = errors.New("product stock quantity is invalid")
	// ErrInsufficientStock 库存不足，卖不出这么多。
	ErrInsufficientStock = errors.New("insufficient product stock")
	// ErrInvalidQuantity 数量必须大于 0（整数件 / 小数量，也不能是 NaN、±Inf）。
	ErrInvalidQuantity = errors.New("product quantity must be greater than zero")

	// --- 充值商品 (TopUpProduct) ---

	// ErrTopUpUnitRequired 充值商品必须写明计量单位（MB / GB / 点数…）。
	ErrTopUpUnitRequired = errors.New("top-up unit is required")
	// ErrInvalidTopUpAmount 起充金额必须大于 0；单笔上限为 0（不限）或不能小于起充。
	ErrInvalidTopUpAmount = errors.New("top-up amount range is invalid")
	// ErrTopUpAmountTooSmall 充值金额低于起充金额。
	ErrTopUpAmountTooSmall = errors.New("top-up amount is below the minimum")
	// ErrTopUpAmountTooLarge 充值金额超过单笔上限。
	ErrTopUpAmountTooLarge = errors.New("top-up amount exceeds the maximum")

	// ErrProductAlreadyOnSale 商品已经在售了。
	ErrProductAlreadyOnSale = errors.New("product is already on sale")
	// ErrProductAlreadyOffShelf 商品已经下架了。
	ErrProductAlreadyOffShelf = errors.New("product is already off shelf")

	// ErrProductRequired 传入的商品为空。
	ErrProductRequired = errors.New("product is required")
	// ErrProductNotFound 指定的商品不存在。
	ErrProductNotFound = errors.New("product not found")
)
