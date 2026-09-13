package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.AbsenceRecordCommand = (*CommandImpl)(nil)

// SaveAbsence 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) SaveAbsence(ctx context.Context, a *absence.AbsenceRecord) error {
	return d.MysqlData.SaveAbsence(ctx, a)
}

// DeleteAbsence 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteAbsence(ctx context.Context, id int64) error {
	return d.MysqlData.DeleteAbsence(ctx, id)
}

// MustGetAbsence 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetAbsence(ctx context.Context, id int64) (absence.AbsenceRecord, error) {
	return d.MysqlData.MustGetAbsence(ctx, id)
}

// GetAbsences 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetAbsences(ctx context.Context, studentID int64) ([]absence.AbsenceRecord, error) {
	return d.MysqlData.GetAbsences(ctx, studentID)
}
