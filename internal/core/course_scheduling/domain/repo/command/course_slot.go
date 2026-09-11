package command

import (
	"context"
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

var (
	// ErrCourseSlotRequired 传入的课表槽位为空。
	ErrCourseSlotRequired = errors.New("course slot is required")
	// ErrCourseSlotNotFound 指定的课表槽位不存在。
	ErrCourseSlotNotFound = errors.New("course slot not found")
	// ErrCourseSlotConflict 目标讲师/教室/课程在该槽位的时间上已有其他安排。
	ErrCourseSlotConflict = errors.New("course slot conflicts with an existing schedule")
)

// ClassroomAssignContext 是仓库侧在排教室时一并交给冲突检查的上下文。
type ClassroomAssignContext struct {
	// TargetSlots 本次要排的槽位（值拷贝，改动不会回写）
	TargetSlots courseSlot.CourseSlots
	// ClassroomSlots 该教室现有的槽位
	ClassroomSlots courseSlot.CourseSlots
	// Course 目标槽位所属课程（用于比对容量）；取不到时为零值
	Course course.Course
	// Classroom 目标教室（用于比对容量）；取不到时为零值
	Classroom classroom.Classroom
}

// CourseAssignContext 是仓库侧在给槽位配课程时一并交给冲突检查的上下文。
type CourseAssignContext struct {
	// TargetSlots 本次要配的槽位（值拷贝，改动不会回写）
	TargetSlots courseSlot.CourseSlots
	// CourseSlots 该课程现有的其他槽位
	CourseSlots courseSlot.CourseSlots
}

// CourseSlotCommand 课表槽位命令接口
//
// 槽位 ID 是聚合生成的 UUID（string），所以 slotIDs 用 []string。
type CourseSlotCommand interface {
	// Save 保存课表槽位（新增或更新）
	Save(cs *courseSlot.CourseSlot) error

	// Delete 删除课表槽位
	Delete(id string) error

	// checkConflictFn 由调用方注入，仓库会把 TeacherAssignContext 装好传进去；
	// 传 nil 表示不做检查。冲突时整批中止，不写入任何槽位。
	AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64,
		checkConflictFn func(ctx context.Context, slots []courseSlot.CourseSlot, t teacher.Teacher) (bool, error)) error

	// checkConflictFn 由调用方注入，仓库会把 CourseAssignContext 装好传进去；
	// 传 nil 表示不做检查。
	AssignCourse(ctx context.Context, slotIDs []string, courseID string, checkConflictFn func(ctx context.Context, ac CourseAssignContext) (bool, error)) error

	// checkConflictFn 由调用方注入，仓库会把 ClassroomAssignContext 装好传进去；
	// 传 nil 表示不做检查。
	AssignClassroom(ctx context.Context, slotIDs []string, classroomID string, checkConflictFn func(ctx context.Context, ac ClassroomAssignContext) (bool, error)) error
}
