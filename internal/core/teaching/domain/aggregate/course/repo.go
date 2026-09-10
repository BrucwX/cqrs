package course

// CourseQuery 课程查询接口
type CourseQuery interface {
	// GetByID 根据 ID 获取课程
	GetByID(id string) (*Course, error)
	// List 获取课程列表
	List() ([]*Course, error)
}

// CourseCommand 课程命令接口
type CourseCommand interface {
	// Save 保存课程（新增或更新）
	Save(c *Course) error
	// Delete 删除课程
	Delete(id string) error
}
