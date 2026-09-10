package classroom

// ClassroomQuery 教室查询接口
type ClassroomQuery interface {
	// GetByID 根据 ID 获取教室
	GetByID(id string) (*Classroom, error)
	// List 获取教室列表
	List() ([]*Classroom, error)
}

// ClassroomCommand 教室命令接口
type ClassroomCommand interface {
	// Save 保存教室（新增或更新）
	Save(c *Classroom) error
	// Delete 删除教室
	Delete(id string) error
}
