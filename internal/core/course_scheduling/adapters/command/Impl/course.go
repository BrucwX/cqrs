package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.CourseCommand = (*CommandImpl)(nil)

// CreateCourse 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) CreateCourse(ctx context.Context, c *course.Course) error {
	return d.MysqlData.CreateCourse(ctx, c)
}

// UpdateCourse 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) UpdateCourse(ctx context.Context, id string, updateFn func(ctx context.Context, crs *course.Course) (*course.Course, error)) error {
	return d.MysqlData.UpdateCourse(ctx, id, updateFn)
}

// DeleteCourse 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteCourse(ctx context.Context, id string) error {
	return d.MysqlData.DeleteCourse(ctx, id)
}

// GetCourse 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetCourse(ctx context.Context, id string) (*course.Course, error) {
	return d.MysqlData.GetCourse(ctx, id)
}

// MustGetCourse 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetCourse(ctx context.Context, id string) (course.Course, error) {
	return d.MysqlData.MustGetCourse(ctx, id)
}

// GetCourses 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetCourses(ctx context.Context, courseTypeID string) ([]course.Course, error) {
	return d.MysqlData.GetCourses(ctx, courseTypeID)
}
