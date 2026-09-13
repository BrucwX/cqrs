package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.StudentCommand = (*CommandImpl)(nil)

// CreateStudent 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) CreateStudent(ctx context.Context, s *student.Student) error {
	return d.MysqlData.CreateStudent(ctx, s)
}

// UpdateStudent 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) UpdateStudent(ctx context.Context, id int64, updateFn func(ctx context.Context, s *student.Student) (*student.Student, error)) error {
	return d.MysqlData.UpdateStudent(ctx, id, updateFn)
}

// DeleteStudent 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteStudent(ctx context.Context, id int64) error {
	return d.MysqlData.DeleteStudent(ctx, id)
}

// GetStudent 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetStudent(ctx context.Context, id int64) (*student.Student, error) {
	return d.MysqlData.GetStudent(ctx, id)
}

// MustGetStudent 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetStudent(ctx context.Context, id int64) (student.Student, error) {
	return d.MysqlData.MustGetStudent(ctx, id)
}
