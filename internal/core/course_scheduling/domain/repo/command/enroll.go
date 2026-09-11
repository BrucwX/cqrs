package command

import (
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// EnrollRepo 是学员选课做冲突检查所需的数据来源。
//
// checkEnrollment 由命令服务自己实现，它需要「目标课程的排期 / 该学员现有在学课程的
// 排期」；仓库只负责按需提供这些数据，不参与规则判定。
type EnrollRepo interface {
	// GetCourseSlots 取某门课程现有的全部排期
	GetCourseSlots(courseID string) (courseSlot.CourseSlots, error)
	// GetStudentSlots 取该学员现有在学课程的全部排期
	GetStudentSlots(studentID int64) (courseSlot.CourseSlots, error)
}
