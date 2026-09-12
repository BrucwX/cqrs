package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

var (
	// ErrCourseTypeNotFound 课程或它归属的课程类型不存在。
	ErrCourseTypeNotFound = errors.New("course type not found")
)

// AssignTeacherRepo 是排讲师做冲突检查所需的数据来源。
//
// checkTeacher 由命令服务自己实现，判定要用到「本次要排的槽位 / 该讲师本人 /
// 该讲师现有排期 / 该讲师资质 / 目标槽位所属课程类型」。回调只收得到 ID，
// 所以这些聚合都由本接口按 ID 取出来；仓库只负责提供数据，不参与规则判定。
type AssignTeacherRepo interface {
	// GetSlots 按 ID 取本次要排的槽位，顺序与入参一致；少一个就报 not found
	GetSlots(slotIDs []string) (courseSlot.CourseSlots, error)
	// GetTeacher 按 ID 取讲师本人，取不到报 not found
	GetTeacher(teacherID int64) (teacher.Teacher, error)
	// GetTeacherSlots 取该讲师现有的全部排期
	GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error)
	// GetQualifications 取该讲师持有的全部资质
	GetQualifications(teacherID int64) ([]qualification.Qualification, error)
	// GetCourseType 取某门课程归属的课程类型
	GetCourseType(courseID string) (courseType.CourseType, error)
}
