package payment

// PaymentQuery 支付查询接口
type PaymentQuery interface {
	// GetByID 根据 ID 获取支付
	GetByID(id string) (*Payment, error)
	// List 获取支付列表
	List() ([]*Payment, error)
	// ListByOrderID 根据订单 ID 获取支付列表
	ListByOrderID(orderID string) ([]*Payment, error)
	// ListByStudentID 根据学员 ID 获取支付列表
	ListByStudentID(studentID string) ([]*Payment, error)
}

// PaymentCommand 支付命令接口
type PaymentCommand interface {
	// Save 保存支付（新增或更新）
	Save(p *Payment) error
	// Delete 删除支付
	Delete(id string) error
}
