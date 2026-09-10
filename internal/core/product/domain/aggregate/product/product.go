package product

import (
	"strings"
	"time"

	"github.com/google/uuid"
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
