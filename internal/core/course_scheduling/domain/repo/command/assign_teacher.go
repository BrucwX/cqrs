package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

var (
	// ErrCourseTypeNotFound 课程或它归属的课程类型不存在。
	ErrCourseTypeNotFound = errors.New("course type not found")
)

// AssignTeacherRepo 是排讲师做冲突检查所需的数据来源。
//
// checkTeacher 由命令服务自己实现，它需要「该讲师现有排期 / 该讲师资质 /
// 目标槽位所属课程类型」；仓库只负责按需提供这些数据，不参与规则判定。
type AssignTeacherRepo interface {
	// GetTeacherSlots 取该讲师现有的全部排期
	GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error)
	// GetQualifications 取该讲师持有的全部资质
	GetQualifications(teacherID int64) ([]qualification.Qualification, error)
	// GetCourseType 取某门课程归属的课程类型
	GetCourseType(courseID string) (courseType.CourseType, error)
}
