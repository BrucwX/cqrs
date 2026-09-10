package absence

// AbsenceRecordQuery 缺勤记录查询接口
type AbsenceRecordQuery interface {
	// GetByID 根据 ID 获取缺勤记录
	GetByID(id int64) (*AbsenceRecord, error)
	// List 获取缺勤记录列表
	List() ([]*AbsenceRecord, error)
	// ListByStudentID 根据学员 ID 获取缺勤记录列表
	ListByStudentID(studentID int64) ([]*AbsenceRecord, error)
	// ListByCourseID 根据课程 ID 获取缺勤记录列表
	ListByCourseID(courseID string) ([]*AbsenceRecord, error)
}

// AbsenceRecordCommand 缺勤记录命令接口
type AbsenceRecordCommand interface {
	// Save 保存缺勤记录（新增或更新）
	Save(a *AbsenceRecord) error
	// Delete 删除缺勤记录
	Delete(id int64) error
}
