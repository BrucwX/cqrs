package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/student"

// StudentCommand 学员命令接口
type StudentCommand interface {
	// Save 保存学员（新增或更新）
	Save(s *student.Student) error
	// Delete 删除学员
	Delete(id int64) error
}
