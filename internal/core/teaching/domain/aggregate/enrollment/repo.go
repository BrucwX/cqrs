package enrollment

// CourseEnrollmentQuery 课程注册查询接口
type CourseEnrollmentQuery interface {
	// GetByID 根据 ID 获取课程注册
	GetByID(id int64) (*CourseEnrollment, error)
	// List 获取课程注册列表
	List() ([]*CourseEnrollment, error)
	// ListByStudentID 根据学员 ID 获取课程注册列表
	ListByStudentID(studentID int64) ([]*CourseEnrollment, error)
	// ListByCourseID 根据课程 ID 获取课程注册列表
	ListByCourseID(courseID string) ([]*CourseEnrollment, error)
}

// CourseEnrollmentCommand 课程注册命令接口
type CourseEnrollmentCommand interface {
	// Save 保存课程注册（新增或更新）
	Save(e *CourseEnrollment) error
	// Delete 删除课程注册
	Delete(id int64) error
}
