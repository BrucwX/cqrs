package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
	"time"
)

var _ repo.StudentMakeupCommand = (*CommandImpl)(nil)

// SaveMakeup 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) SaveMakeup(ctx context.Context, m *makeup.StudentMakeup) error {
	return d.MysqlData.SaveMakeup(ctx, m)
}

// DeleteMakeup 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteMakeup(ctx context.Context, id int64) error {
	return d.MysqlData.DeleteMakeup(ctx, id)
}

// MustGetMakeup 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetMakeup(ctx context.Context, id int64) (makeup.StudentMakeup, error) {
	return d.MysqlData.MustGetMakeup(ctx, id)
}

// GetMakeupsForTarget 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetMakeupsForTarget(ctx context.Context, targetSlotID string, targetDate time.Time) ([]makeup.StudentMakeup, error) {
	return d.MysqlData.GetMakeupsForTarget(ctx, targetSlotID, targetDate)
}
