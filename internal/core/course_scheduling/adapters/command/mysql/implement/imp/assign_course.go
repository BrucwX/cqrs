package imp

import (
	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type AssignCourseImp struct {
	data *mysql.Data
}

var _ repo.AssignCourseRepo = (*AssignCourseImp)(nil)

func NewAssignCourseImp(d *mysql.Data) repo.AssignCourseRepo {
	return &AssignCourseImp{data: d}
}

func (c *AssignCourseImp) GetCourseSlots(courseID string) (courseSlot.CourseSlots, error) {
	panic("implement me")
}
