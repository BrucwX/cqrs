package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

var (
	// ErrCourseNotFound 指定的课程不存在。
	ErrCourseNotFound = errors.New("course not found")
)

// AssignCourseRepo 是给槽位配课程做冲突检查所需的数据来源。
//
// checkCourse 由命令服务自己实现，它需要「该课程现有的排期」；
// 仓库只负责按需提供这些数据，不参与规则判定。
type AssignCourseRepo interface {
	// GetCourseSlots 取该课程现有的全部排期
	GetCourseSlots(courseID string) (courseSlot.CourseSlots, error)
}
