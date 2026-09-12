package imp

import (
	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type SlotChangeImp struct {
	data *mysql.Data
}

var _ repo.SlotChangeRepo = (*SlotChangeImp)(nil)

func NewSlotChangeImp(d *mysql.Data) repo.SlotChangeRepo {
	return &SlotChangeImp{data: d}
}

func (c *SlotChangeImp) GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error) {
	panic("implement me")
}

func (c *SlotChangeImp) GetClassroomSlots(classroomID string) (courseSlot.CourseSlots, error) {
	panic("implement me")
}

func (c *SlotChangeImp) GetOtherSlotChanges(id int64) ([]*courseSlotChange.CourseSlotChange, error) {
	panic("implement me")
}
