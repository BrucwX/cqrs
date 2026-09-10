package qualification

// QualificationQuery 授课资质查询接口
type QualificationQuery interface {
	// GetByID 根据 ID 获取授课资质
	GetByID(id int64) (*Qualification, error)
	// List 获取授课资质列表
	List() ([]*Qualification, error)
	// ListByTeacherID 根据讲师 ID 获取授课资质列表
	ListByTeacherID(teacherID int64) ([]*Qualification, error)
	// ListByCourseID 根据课程 ID 获取授课资质列表
	ListByCourseID(courseID string) ([]*Qualification, error)
}

// QualificationCommand 授课资质命令接口
type QualificationCommand interface {
	// Save 保存授课资质（新增或更新）
	Save(q *Qualification) error
	// Delete 删除授课资质
	Delete(id int64) error
}
