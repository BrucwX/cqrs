package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// CourseSlotCommand 课表槽位命令接口
type CourseSlotCommand interface {
	// Save 保存课表槽位（新增或更新）
	Save(cs *courseSlot.CourseSlot) error

	// Delete 删除课表槽位
	Delete(id int64) error

	Update(ctx context.Context, slotIDs []int64, teacherID int64) error

	// AssignTeacher 给指定课表槽位们配置老师
	AssignTeacher(ctx context.Context, slotIDs []int64, teacherID int64, checkConflictFn func(ctx context.Context, cs *courseSlot.CourseSlot) (bool, error)) error

	// AssignClassroom 给指定课表槽位们配置教室
	AssignClassroom(ctx context.Context, slotIDs []int64, classroomID string, checkConflictFn func(ctx context.Context, cs *courseSlot.CourseSlot) (bool, error)) error
}
