package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// CourseSlotCommand 课表槽位命令接口
//
// 槽位 ID 是聚合生成的 UUID（string），所以 slotIDs 用 []string。
//
// 三个 AssignXxx 只管写：把目标 ID 写进这批槽位。冲突判定不在仓库里 ——
// 那是调用方的事（见 domain/service 的 schedule.Conflict / qualification.Check），
// 判定通过才调过来。所以这里没有回调，也没有「传 nil 表示不检查」这类分支。
type CourseSlotCommand interface {
	// SaveCourseSlot 保存课表槽位（新增或更新）
	SaveCourseSlot(ctx context.Context, cs *courseSlot.CourseSlot) error

	// DeleteCourseSlot 删除课表槽位
	DeleteCourseSlot(ctx context.Context, id string) error

	// MustGetCourseSlot 取课表槽位聚合；不存在时报 ErrCourseSlotNotFound
	MustGetCourseSlot(ctx context.Context, id string) (courseSlot.CourseSlot, error)

	// AssignTeacher 给指定课表槽位们配置老师
	AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64) error

	// AssignCourse 给指定课表槽位们配置课程
	AssignCourse(ctx context.Context, slotIDs []string, courseID string) error

	// AssignClassroom 给指定课表槽位们配置教室
	AssignClassroom(ctx context.Context, slotIDs []string, classroomID string) error

	// --- 下面都是「取排期」，按返回值的聚合根归到本接口 ---

	// GetSlots 按 ID 取本次要排的槽位，顺序与入参一致；少一个就报 not found
	GetSlots(ctx context.Context, slotIDs []string) (courseSlot.CourseSlots, error)
	// GetTeacherSlots 取该讲师现有的全部排期
	GetTeacherSlots(ctx context.Context, teacherID int64) (courseSlot.CourseSlots, error)
	// GetClassroomSlots 取该教室现有的全部排期
	GetClassroomSlots(ctx context.Context, classroomID string) (courseSlot.CourseSlots, error)
	// GetCourseSlots 取该课程现有的全部排期
	GetCourseSlots(ctx context.Context, courseID string) (courseSlot.CourseSlots, error)
	// GetStudentSlots 取该学员在学课程的全部排期
	GetStudentSlots(ctx context.Context, studentID int64) (courseSlot.CourseSlots, error)
}
