package command

import (
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// AssignClassroomRepo 是排教室做冲突检查所需的数据来源。
//
// checkClassroom 由命令服务自己实现，它需要「该教室现有的排期 / 目标槽位所属课程
// （用来比对容量）」；仓库只负责按需提供这些数据，不参与规则判定。
type AssignClassroomRepo interface {
	// GetClassroomSlots 取该教室现有的全部排期
	GetClassroomSlots(classroomID string) (courseSlot.CourseSlots, error)
	// GetCourse 取课程本身
	GetCourse(courseID string) (course.Course, error)
}
