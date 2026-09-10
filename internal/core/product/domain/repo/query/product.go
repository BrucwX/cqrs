package query

import "cqrs/internal/core/product/domain/aggregate/product"

// ProductQuery 商品查询接口
type ProductQuery interface {
	// GetByID 根据 ID 获取商品
	GetByID(id string) (*product.Product, error)
	// List 获取商品列表
	List() ([]*product.Product, error)
	// ListByCategory 根据分类获取商品列表
	ListByCategory(category product.ProductCategory) ([]*product.Product, error)
}
