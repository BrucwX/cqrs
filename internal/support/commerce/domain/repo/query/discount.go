package query

import "cqrs/internal/support/commerce/domain/aggregate/discount"

// DiscountQuery 折扣查询接口
type DiscountQuery interface {
	// GetByID 根据 ID 获取折扣
	GetByID(id string) (*discount.Discount, error)
	// List 获取折扣列表
	List() ([]*discount.Discount, error)
	// GetByCode 根据优惠码获取折扣
	GetByCode(code string) (*discount.Discount, error)
}
