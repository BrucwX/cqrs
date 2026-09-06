package product

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// ProductStatus is the lifecycle state of a product.
type ProductStatus int32

const (
	ProductStatusUnspecified  ProductStatus = 0
	ProductStatusActive       ProductStatus = 1 // 在售
	ProductStatusOutOfStock   ProductStatus = 2 // 缺货
	ProductStatusDiscontinued ProductStatus = 3 // 下架
	ProductStatusDeleted      ProductStatus = 4
)

// ProductCategory is the category of a product.
type ProductCategory int32

const (
	ProductCategoryUnspecified ProductCategory = 0
	ProductCategoryEquipment   ProductCategory = 1 // 器材
	ProductCategoryClothing    ProductCategory = 2 // 服装
	ProductCategoryAccessory   ProductCategory = 3 // 配件
	ProductCategoryConsumable  ProductCategory = 4 // 消耗品
)

// Product is a product domain object (teaching supplies).
type Product struct {
	ID          string
	Name        string
	Description string
	Category    ProductCategory
	Price       int64  // 价格，单位：分
	Stock       int32  // 库存数量
	ImageURL    string // 商品图片
	Status      ProductStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProduct creates a new product.
func NewProduct(name string, category ProductCategory, price int64) *Product {
	return &Product{
		ID:        uuid.New().String(),
		Name:      name,
		Category:  category,
		Price:     price,
		Status:    ProductStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate checks if the product domain object is valid.
func (p *Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrProductInvalidArgument
	}
	if p.Price < 0 {
		return ErrProductInvalidArgument
	}
	if p.Stock < 0 {
		return ErrProductInvalidArgument
	}
	return nil
}
