package courseSlotChange

// CourseSlotChangeQuery 课表变更查询接口
type CourseSlotChangeQuery interface {
	// GetByID 根据 ID 获取课表变更
	GetByID(id int64) (*CourseSlotChange, error)
	// List 获取课表变更列表
	List() ([]*CourseSlotChange, error)
	// ListByCourseID 根据课程 ID 获取课表变更列表
	ListByCourseID(courseID string) ([]*CourseSlotChange, error)
}

// CourseSlotChangeCommand 课表变更命令接口
type CourseSlotChangeCommand interface {
	// Save 保存课表变更（新增或更新）
	Save(csc *CourseSlotChange) error
	// Delete 删除课表变更
	Delete(id int64) error
}
