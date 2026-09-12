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

func (c *CourseImp) Create(ctx context.Context, crs *course.Course) error {
	panic("implement me")
}

func (c *CourseImp) Update(
	ctx context.Context,
	id string,
	updateFn func(ctx context.Context, crs *course.Course) (*course.Course, error),
) error {
	panic("implement me")
}

func (c *CourseImp) Delete(ctx context.Context, id string) error {
	panic("implement me")
}

func (c *CourseImp) Get(ctx context.Context, id string) (*course.Course, error) {
	panic("implement me")
}

func (c *CourseImp) MustGet(ctx context.Context, id string) (course.Course, error) {
	panic("implement me")
}

func (c *CourseImp) GetCourses(ctx context.Context, courseTypeID string) ([]course.Course, error) {
	panic("implement me")
}
