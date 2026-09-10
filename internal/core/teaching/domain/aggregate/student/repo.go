package student

// StudentQuery 学员查询接口
type StudentQuery interface {
	// GetByID 根据 ID 获取学员
	GetByID(id int64) (*Student, error)
	// List 获取学员列表
	List() ([]*Student, error)
}

// StudentCommand 学员命令接口
type StudentCommand interface {
	// Save 保存学员（新增或更新）
	Save(s *Student) error
	// Delete 删除学员
	Delete(id int64) error
}
