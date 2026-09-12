package imp

import (
	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type EnrollImp struct {
	data *mysql.Data
}

var _ repo.EnrollRepo = (*EnrollImp)(nil)

func NewEnrollImp(d *mysql.Data) repo.EnrollRepo {
	return &EnrollImp{data: d}
}

func (c *EnrollImp) GetCourseSlots(courseID string) (courseSlot.CourseSlots, error) {
	panic("implement me")
}

func (c *EnrollImp) GetStudentSlots(studentID int64) (courseSlot.CourseSlots, error) {
	panic("implement me")
}
