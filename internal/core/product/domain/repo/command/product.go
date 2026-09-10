package command

import "cqrs/internal/core/product/domain/aggregate/product"

// ProductCommand 商品命令接口
type ProductCommand interface {
	// Save 保存商品（新增或更新）
	Save(p *product.Product) error
	// Delete 删除商品
	Delete(id string) error
}
