package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.CourseTypeCommand = (*CommandImpl)(nil)

// MustGetCourseType 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetCourseType(ctx context.Context, id string) (courseType.CourseType, error) {
	return d.MysqlData.MustGetCourseType(ctx, id)
}

// GetCourseType 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetCourseType(ctx context.Context, courseID string) (courseType.CourseType, error) {
	return d.MysqlData.GetCourseType(ctx, courseID)
}
