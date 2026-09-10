package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/qualification"

// QualificationCommand 授课资质命令接口
type QualificationCommand interface {
	// Save 保存授课资质（新增或更新）
	Save(q *qualification.Qualification) error
	// Delete 删除授课资质
	Delete(id int64) error
}
