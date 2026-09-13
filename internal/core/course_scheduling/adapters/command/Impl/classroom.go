package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.ClassroomCommand = (*CommandImpl)(nil)

// CreateClassroom 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) CreateClassroom(ctx context.Context, c *classroom.Classroom) error {
	return d.MysqlData.CreateClassroom(ctx, c)
}

// UpdateClassroom 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) UpdateClassroom(ctx context.Context, id string, updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error)) error {
	return d.MysqlData.UpdateClassroom(ctx, id, updateFn)
}

// DeleteClassroom 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteClassroom(ctx context.Context, id string) error {
	return d.MysqlData.DeleteClassroom(ctx, id)
}

// MustGetClassroom 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetClassroom(ctx context.Context, id string) (classroom.Classroom, error) {
	return d.MysqlData.MustGetClassroom(ctx, id)
}
