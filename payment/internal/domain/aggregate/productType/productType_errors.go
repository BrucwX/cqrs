package productType

import "errors"

var (
	// --- 商品分类 (ProductType) 校验 ---

	// ErrProductTypeNameRequired 分类必须有名字。
	ErrProductTypeNameRequired = errors.New("product type name is required")
	// ErrProductTypeCodeRequired 分类必须有业务编码。
	ErrProductTypeCodeRequired = errors.New("product type code is required")

	// --- 启用 / 停用 ---

	// ErrProductTypeAlreadyEnabled 分类已经启用了。
	ErrProductTypeAlreadyEnabled = errors.New("product type is already enabled")
	// ErrProductTypeAlreadyDisabled 分类已经停用了。
	ErrProductTypeAlreadyDisabled = errors.New("product type is already disabled")

	// --- 仓库层用 ---

	// ErrProductTypeRequired 传入的商品分类为空。
	ErrProductTypeRequired = errors.New("product type is required")
	// ErrProductTypeNotFound 指定的商品分类不存在。
	ErrProductTypeNotFound = errors.New("product type not found")
)
