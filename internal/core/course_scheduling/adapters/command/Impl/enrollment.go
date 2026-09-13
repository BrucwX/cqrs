package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.CourseEnrollmentCommand = (*CommandImpl)(nil)

// SaveEnrollment 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) SaveEnrollment(ctx context.Context, e *enrollment.CourseEnrollment) error {
	return d.MysqlData.SaveEnrollment(ctx, e)
}

// DeleteEnrollment 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteEnrollment(ctx context.Context, id int64) error {
	return d.MysqlData.DeleteEnrollment(ctx, id)
}

// MustGetEnrollment 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetEnrollment(ctx context.Context, id int64) (enrollment.CourseEnrollment, error) {
	return d.MysqlData.MustGetEnrollment(ctx, id)
}

// Enroll 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) Enroll(ctx context.Context, e *enrollment.CourseEnrollment) error {
	return d.MysqlData.Enroll(ctx, e)
}

// GetEnrollments 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetEnrollments(ctx context.Context, studentID int64) ([]enrollment.CourseEnrollment, error) {
	return d.MysqlData.GetEnrollments(ctx, studentID)
}
