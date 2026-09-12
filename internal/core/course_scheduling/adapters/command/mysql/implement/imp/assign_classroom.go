package imp

import (
	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type AssignClassroomImp struct {
	data *mysql.Data
}

var _ repo.AssignClassroomRepo = (*AssignClassroomImp)(nil)

func NewAssignClassroomImp(d *mysql.Data) repo.AssignClassroomRepo {
	return &AssignClassroomImp{data: d}
}

func (c *AssignClassroomImp) GetClassroomSlots(classroomID string) (courseSlot.CourseSlots, error) {
	panic("implement me")
}

func (c *AssignClassroomImp) GetCourse(courseID string) (course.Course, error) {
	panic("implement me")
}
