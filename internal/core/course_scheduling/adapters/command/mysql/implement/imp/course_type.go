package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type CourseTypeImp struct {
	data *mysql.Data
}

var _ repo.CourseTypeCommand = (*CourseTypeImp)(nil)

func NewCourseTypeImp(d *mysql.Data) repo.CourseTypeCommand {
	return &CourseTypeImp{data: d}
}

func (c *CourseTypeImp) GetCourseType(ctx context.Context, courseID string) (courseType.CourseType, error) {
	panic("implement me")
}
