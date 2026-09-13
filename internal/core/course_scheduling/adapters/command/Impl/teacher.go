package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.TeacherCommand = (*CommandImpl)(nil)

// CreateTeacher 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) CreateTeacher(ctx context.Context, t *teacher.Teacher) error {
	return d.MysqlData.CreateTeacher(ctx, t)
}

// UpdateTeacher 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) UpdateTeacher(ctx context.Context, id int64, updateFn func(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error)) error {
	return d.MysqlData.UpdateTeacher(ctx, id, updateFn)
}

// DeleteTeacher 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteTeacher(ctx context.Context, id int64) error {
	return d.MysqlData.DeleteTeacher(ctx, id)
}

// GetTeacher 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetTeacher(ctx context.Context, id int64) (*teacher.Teacher, error) {
	return d.MysqlData.GetTeacher(ctx, id)
}

// MustGetTeacher 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetTeacher(ctx context.Context, teacherID int64) (teacher.Teacher, error) {
	return d.MysqlData.MustGetTeacher(ctx, teacherID)
}
