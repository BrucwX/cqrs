package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type CourseImp struct {
	data *mysql.Data
}

var _ repo.CourseCommand = (*CourseImp)(nil)

func NewCourseImp(d *mysql.Data) repo.CourseCommand {
	return &CourseImp{data: d}
}

func (c *CourseImp) Create(crs *course.Course) error {
	panic("implement me")
}

func (c *CourseImp) Update(
	ctx context.Context,
	id string,
	updateFn func(ctx context.Context, crs *course.Course) (*course.Course, error),
) error {
	panic("implement me")
}

func (c *CourseImp) Delete(id string) error {
	panic("implement me")
}

func (c *CourseImp) Get(id string) (*course.Course, error) {
	panic("implement me")
}
