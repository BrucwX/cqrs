package discount

// DiscountQuery 折扣查询接口
type DiscountQuery interface {
	// GetByID 根据 ID 获取折扣
	GetByID(id string) (*Discount, error)
	// List 获取折扣列表
	List() ([]*Discount, error)
	// GetByCode 根据优惠码获取折扣
	GetByCode(code string) (*Discount, error)
}

// DiscountCommand 折扣命令接口
type DiscountCommand interface {
	// Save 保存折扣（新增或更新）
	Save(d *Discount) error
	// Delete 删除折扣
	Delete(id string) error
}
