package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.QualificationCommand = (*CommandImpl)(nil)

// GrantQualification 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GrantQualification(ctx context.Context, q *qualification.Qualification) error {
	return d.MysqlData.GrantQualification(ctx, q)
}

// DeleteQualification 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteQualification(ctx context.Context, id int64) error {
	return d.MysqlData.DeleteQualification(ctx, id)
}

// MustGetQualification 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetQualification(ctx context.Context, id int64) (qualification.Qualification, error) {
	return d.MysqlData.MustGetQualification(ctx, id)
}

// GetQualifications 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetQualifications(ctx context.Context, teacherID int64) ([]qualification.Qualification, error) {
	return d.MysqlData.GetQualifications(ctx, teacherID)
}
