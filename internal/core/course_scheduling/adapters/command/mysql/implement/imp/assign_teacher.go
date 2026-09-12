package imp

import (
	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

/*
type AssignTeacherRepo interface {
	// GetTeacherSlots 取该讲师现有的全部排期
	GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error)
	// GetQualifications 取该讲师持有的全部资质
	GetQualifications(teacherID int64) ([]qualification.Qualification, error)
	// GetCourseType 取某门课程归属的课程类型
	GetCourseType(courseID string) (courseType.CourseType, error)
}
*/

type AssignTeacherImp struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repo.AssignTeacherRepo = (*AssignTeacherImp)(nil)

// NewClassroomQuery 创建 MySQL 版教室查询。
func NewAssignTeacherImp(d *mysql.Data) repo.AssignTeacherRepo {
	return &AssignTeacherImp{data: d}
}

func (c *AssignTeacherImp) GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error) {

}

func (c *AssignTeacherImp) GetQualifications(teacherID int64) ([]qualification.Qualification, error) {

}

func (c *AssignTeacherImp) GetCourseType(courseID string) (courseType.CourseType, error) {

}
