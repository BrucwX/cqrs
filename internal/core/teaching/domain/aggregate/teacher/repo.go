package teacher

// TeacherQuery 讲师查询接口
type TeacherQuery interface {
	// GetByID 根据 ID 获取讲师
	GetByID(id int64) (*Teacher, error)
	// List 获取讲师列表
	List() ([]*Teacher, error)
}

// TeacherCommand 讲师命令接口
type TeacherCommand interface {
	// Save 保存讲师（新增或更新）
	Save(t *Teacher) error
	// Delete 删除讲师
	Delete(id int64) error
}
