package query

import "cqrs/internal/support/commerce/domain/aggregate/order"

// OrderQuery 订单查询接口
type OrderQuery interface {
	// GetByID 根据 ID 获取订单
	GetByID(id string) (*order.Order, error)
	// List 获取订单列表
	List() ([]*order.Order, error)
	// ListByStudentID 根据学员 ID 获取订单列表
	ListByStudentID(studentID string) ([]*order.Order, error)
}
