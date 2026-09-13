package command

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

var _ repo.CourseSlotCommand = (*CommandImpl)(nil)

// SaveCourseSlot 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) SaveCourseSlot(ctx context.Context, cs *courseSlot.CourseSlot) error {
	return d.MysqlData.SaveCourseSlot(ctx, cs)
}

// DeleteCourseSlot 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) DeleteCourseSlot(ctx context.Context, id string) error {
	return d.MysqlData.DeleteCourseSlot(ctx, id)
}

// MustGetCourseSlot 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) MustGetCourseSlot(ctx context.Context, id string) (courseSlot.CourseSlot, error) {
	return d.MysqlData.MustGetCourseSlot(ctx, id)
}

// AssignTeacher 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64) error {
	return d.MysqlData.AssignTeacher(ctx, slotIDs, teacherID)
}

// AssignCourse 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) AssignCourse(ctx context.Context, slotIDs []string, courseID string) error {
	return d.MysqlData.AssignCourse(ctx, slotIDs, courseID)
}

// AssignClassroom 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) AssignClassroom(ctx context.Context, slotIDs []string, classroomID string) error {
	return d.MysqlData.AssignClassroom(ctx, slotIDs, classroomID)
}

// GetSlots 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetSlots(ctx context.Context, slotIDs []string) (courseSlot.CourseSlots, error) {
	return d.MysqlData.GetSlots(ctx, slotIDs)
}

// GetTeacherSlots 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetTeacherSlots(ctx context.Context, teacherID int64) (courseSlot.CourseSlots, error) {
	return d.MysqlData.GetTeacherSlots(ctx, teacherID)
}

// GetClassroomSlots 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetClassroomSlots(ctx context.Context, classroomID string) (courseSlot.CourseSlots, error) {
	return d.MysqlData.GetClassroomSlots(ctx, classroomID)
}

// GetCourseSlots 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetCourseSlots(ctx context.Context, courseID string) (courseSlot.CourseSlots, error) {
	return d.MysqlData.GetCourseSlots(ctx, courseID)
}

// GetStudentSlots 转发给 MysqlData，等 Redis 接入后在这里组合两者。
func (d *CommandImpl) GetStudentSlots(ctx context.Context, studentID int64) (courseSlot.CourseSlots, error) {
	return d.MysqlData.GetStudentSlots(ctx, studentID)
}
