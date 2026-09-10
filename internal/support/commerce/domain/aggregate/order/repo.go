package order

// OrderQuery 订单查询接口
type OrderQuery interface {
	// GetByID 根据 ID 获取订单
	GetByID(id string) (*Order, error)
	// List 获取订单列表
	List() ([]*Order, error)
	// ListByStudentID 根据学员 ID 获取订单列表
	ListByStudentID(studentID string) ([]*Order, error)
}

// OrderCommand 订单命令接口
type OrderCommand interface {
	// Save 保存订单（新增或更新）
	Save(o *Order) error
	// Delete 删除订单
	Delete(id string) error
}
