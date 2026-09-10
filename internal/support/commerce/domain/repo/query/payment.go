package query

import "cqrs/internal/support/commerce/domain/aggregate/payment"

// PaymentQuery 支付查询接口
type PaymentQuery interface {
	// GetByID 根据 ID 获取支付
	GetByID(id string) (*payment.Payment, error)
	// List 获取支付列表
	List() ([]*payment.Payment, error)
	// ListByOrderID 根据订单 ID 获取支付列表
	ListByOrderID(orderID string) ([]*payment.Payment, error)
	// ListByStudentID 根据学员 ID 获取支付列表
	ListByStudentID(studentID string) ([]*payment.Payment, error)
}
