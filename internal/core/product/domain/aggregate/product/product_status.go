package product

// ProductStatus is the lifecycle state of a product.
type ProductStatus int32

const (
	ProductStatusUnspecified  ProductStatus = 0
	ProductStatusActive       ProductStatus = 1 // 在售
	ProductStatusOutOfStock   ProductStatus = 2 // 缺货
	ProductStatusDiscontinued ProductStatus = 3 // 下架
	ProductStatusDeleted      ProductStatus = 4
)
