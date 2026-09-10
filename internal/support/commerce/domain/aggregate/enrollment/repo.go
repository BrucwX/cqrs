package enrollment

// EnrollmentQuery 注册查询接口
type EnrollmentQuery interface {
	// GetByID 根据 ID 获取注册
	GetByID(id string) (*Enrollment, error)
	// List 获取注册列表
	List() ([]*Enrollment, error)
	// ListByStudentID 根据学员 ID 获取注册列表
	ListByStudentID(studentID string) ([]*Enrollment, error)
	// ListByCourseID 根据课程 ID 获取注册列表
	ListByCourseID(courseID string) ([]*Enrollment, error)
}

// EnrollmentCommand 注册命令接口
type EnrollmentCommand interface {
	// Save 保存注册（新增或更新）
	Save(e *Enrollment) error
	// Delete 删除注册
	Delete(id string) error
}
