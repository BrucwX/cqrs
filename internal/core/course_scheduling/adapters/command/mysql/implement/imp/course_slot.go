package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type CourseSlotImp struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repo.CourseSlotCommand = (*CourseSlotImp)(nil)

// NewClassroomQuery 创建 MySQL 版教室查询。
func NewCourseSlotImp(d *mysql.Data) repo.CourseSlotCommand {
	return &CourseSlotImp{data: d}
}

func (c *CourseSlotImp) Save(cs *courseSlot.CourseSlot) error {
	panic("implement me")
}

func (c *CourseSlotImp) Delete(id string) error {
	panic("implement me")
}

func (c *CourseSlotImp) AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64,
	checkConflictFn func(ctx context.Context, slotIDs []string, teacherID int64) (bool, error)) error {
	panic("implement me")
}

func (c *CourseSlotImp) AssignCourse(ctx context.Context, slotIDs []string, courseID string,
	checkConflictFn func(ctx context.Context, slots []courseSlot.CourseSlot, c course.Course) (bool, error)) error {
	panic("implement me")
}

func (c *CourseSlotImp) AssignClassroom(ctx context.Context, slotIDs []string, classroomID string,
	checkConflictFn func(ctx context.Context, slots []courseSlot.CourseSlot, c classroom.Classroom) (bool, error)) error {
	panic("implement me")
}
