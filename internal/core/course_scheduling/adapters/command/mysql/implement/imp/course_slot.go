package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
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

func (c *CourseSlotImp) AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64) error {
	panic("implement me")
}

func (c *CourseSlotImp) AssignCourse(ctx context.Context, slotIDs []string, courseID string) error {
	panic("implement me")
}

func (c *CourseSlotImp) AssignClassroom(ctx context.Context, slotIDs []string, classroomID string) error {
	panic("implement me")
}

// --- 排期读取（按返回值的聚合根归到本接口）---

func (c *CourseSlotImp) GetSlots(slotIDs []string) (courseSlot.CourseSlots, error) {
	panic("implement me")
}

func (c *CourseSlotImp) GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error) {
	panic("implement me")
}

func (c *CourseSlotImp) GetClassroomSlots(classroomID string) (courseSlot.CourseSlots, error) {
	panic("implement me")
}

func (c *CourseSlotImp) GetCourseSlots(courseID string) (courseSlot.CourseSlots, error) {
	panic("implement me")
}

func (c *CourseSlotImp) GetStudentSlots(studentID int64) (courseSlot.CourseSlots, error) {
	panic("implement me")
}
