package product

// ProductCategory is the category of a product.
type ProductCategory int32

const (
	ProductCategoryUnspecified ProductCategory = 0
	ProductCategoryEquipment   ProductCategory = 1 // 器材
	ProductCategoryClothing    ProductCategory = 2 // 服装
	ProductCategoryAccessory   ProductCategory = 3 // 配件
	ProductCategoryConsumable  ProductCategory = 4 // 消耗品
)
