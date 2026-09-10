package courseSlot

// CourseSlotQuery 课表槽位查询接口
type CourseSlotQuery interface {
	// GetByID 根据 ID 获取课表槽位
	GetByID(id int64) (*CourseSlot, error)
	// List 获取课表槽位列表
	List() ([]*CourseSlot, error)
	// ListByCourseID 根据课程 ID 获取课表槽位列表
	ListByCourseID(courseID string) ([]*CourseSlot, error)
}

// CourseSlotCommand 课表槽位命令接口
type CourseSlotCommand interface {
	// Save 保存课表槽位（新增或更新）
	Save(cs *CourseSlot) error
	// Delete 删除课表槽位
	Delete(id int64) error
}
