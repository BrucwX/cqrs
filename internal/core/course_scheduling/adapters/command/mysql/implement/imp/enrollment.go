package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type CourseEnrollmentImp struct {
	data *mysql.Data
}

var _ repo.CourseEnrollmentCommand = (*CourseEnrollmentImp)(nil)

func NewCourseEnrollmentImp(d *mysql.Data) repo.CourseEnrollmentCommand {
	return &CourseEnrollmentImp{data: d}
}

func (c *CourseEnrollmentImp) Save(e *enrollment.CourseEnrollment) error {
	panic("implement me")
}

func (c *CourseEnrollmentImp) Delete(id int64) error {
	panic("implement me")
}

func (c *CourseEnrollmentImp) Enroll(ctx context.Context, e *enrollment.CourseEnrollment,
	checkConflictFn func(ctx context.Context, e *enrollment.CourseEnrollment, c course.Course) (bool, error)) error {
	panic("implement me")
}
