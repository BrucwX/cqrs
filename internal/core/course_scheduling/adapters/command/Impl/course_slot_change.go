package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.CourseSlotChangeCommand = (*CommandImpl)(nil)

// SaveCourseSlotChange 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) SaveCourseSlotChange(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error {
	return d.MysqlData.SaveCourseSlotChange(ctx, csc)
}

// DeleteCourseSlotChange 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteCourseSlotChange(ctx context.Context, id int64) error {
	return d.MysqlData.DeleteCourseSlotChange(ctx, id)
}

// MustGetCourseSlotChange 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetCourseSlotChange(ctx context.Context, id int64) (courseSlotChange.CourseSlotChange, error) {
	return d.MysqlData.MustGetCourseSlotChange(ctx, id)
}

// Change 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error {
	return d.MysqlData.Change(ctx, csc)
}

// GetOtherSlotChanges 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetOtherSlotChanges(ctx context.Context, id int64) ([]*courseSlotChange.CourseSlotChange, error) {
	return d.MysqlData.GetOtherSlotChanges(ctx, id)
}
