package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

type CourseSlotCommand interface {
	// Save 保存课表槽位（新增或更新）
	Save(cs *courseSlot.CourseSlot) error

	// Delete 删除课表槽位
	Delete(id string) error

	// AssignTeacher 给指定课表槽位们配置老师
	//
	// 传 nil 表示不做检查。冲突时整批中止，不写入任何槽位。
	AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64,
		checkConflictFn func(ctx context.Context, slots []courseSlot.CourseSlot, t teacher.Teacher) (bool, error)) error

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
