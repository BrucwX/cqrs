package makeup

// StudentMakeupQuery 补课申请查询接口
type StudentMakeupQuery interface {
	// GetByID 根据 ID 获取补课申请
	GetByID(id int64) (*StudentMakeup, error)
	// List 获取补课申请列表
	List() ([]*StudentMakeup, error)
	// ListByStudentID 根据学员 ID 获取补课申请列表
	ListByStudentID(studentID int64) ([]*StudentMakeup, error)
	// ListByCourseID 根据课程 ID 获取补课申请列表
	ListByCourseID(courseID string) ([]*StudentMakeup, error)
}

// StudentMakeupCommand 补课申请命令接口
type StudentMakeupCommand interface {
	// Save 保存补课申请（新增或更新）
	Save(m *StudentMakeup) error
	// Delete 删除补课申请
	Delete(id int64) error
}
