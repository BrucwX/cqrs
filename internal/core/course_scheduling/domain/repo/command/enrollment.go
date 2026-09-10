package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"

// CourseEnrollmentCommand 课程注册命令接口
type CourseEnrollmentCommand interface {
	// Save 保存课程注册（新增或更新）
	Save(e *enrollment.CourseEnrollment) error
	// Delete 删除课程注册
	Delete(id int64) error
}
