package command

import (
	"context"
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

var (
	// ErrCourseSlotRequired 传入的课表槽位为空。
	ErrCourseSlotRequired = errors.New("course slot is required")
	// ErrCourseSlotNotFound 指定的课表槽位不存在。
	ErrCourseSlotNotFound = errors.New("course slot not found")
	// ErrCourseSlotConflict 目标讲师/教室/课程在该槽位的时间上已有其他安排。
	ErrCourseSlotConflict = errors.New("course slot conflicts with an existing schedule")
)

// CourseSlotCommand 课表槽位命令接口
//
// 槽位 ID 是聚合生成的 UUID（string），所以 slotIDs 用 []string。
// AssignTeacher 只把「本次要排的槽位 ID + 讲师 ID」交给调用方注入的
// checkConflictFn 判定：仓库既不认识业务规则，也不替调用方查数据 ——
// 判定要用的聚合由调用方自己按 ID 取（见 AssignTeacherRepo）。
// AssignCourse / AssignClassroom 暂时还是把「槽位 + 目标聚合」传进去。
type CourseSlotCommand interface {
	// Save 保存课表槽位（新增或更新）
	Save(cs *courseSlot.CourseSlot) error

	// Delete 删除课表槽位
	Delete(id string) error

	// AssignTeacher 给指定课表槽位们配置老师
	//
	// 传 nil 表示不做检查。冲突时整批中止，不写入任何槽位。
	AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64,
		checkConflictFn func(ctx context.Context, slotIDs []string, teacherID int64) (bool, error)) error

	// AssignCourse 给指定课表槽位们配置课程
	//
	// 传 nil 表示不做检查。冲突时整批中止，不写入任何槽位。
	AssignCourse(ctx context.Context, slotIDs []string, courseID string,
		checkConflictFn func(ctx context.Context, slots []courseSlot.CourseSlot, c course.Course) (bool, error)) error

	// AssignClassroom 给指定课表槽位们配置教室
	//
	// 传 nil 表示不做检查。冲突时整批中止，不写入任何槽位。
	AssignClassroom(ctx context.Context, slotIDs []string, classroomID string,
		checkConflictFn func(ctx context.Context, slots []courseSlot.CourseSlot, c classroom.Classroom) (bool, error)) error
}
