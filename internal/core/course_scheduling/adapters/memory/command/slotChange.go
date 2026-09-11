package command

import (
	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// SlotChangeRepo 是 repo.SlotChangeRepo 的内存实现。
//
// 它只做「按需取数据」：换课的规则判定在命令服务里（checkSlotChange），
// 这里不参与任何规则。
type SlotChangeRepo struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足接口。
var _ repo.SlotChangeRepo = (*SlotChangeRepo)(nil)

// NewSlotChangeRepo 创建换课冲突检查的数据来源。
func NewSlotChangeRepo(d *memory.Data) repo.SlotChangeRepo {
	return &SlotChangeRepo{data: d}
}

// GetTeacherSlots 取该讲师现有的全部排期。
func (r *SlotChangeRepo) GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range r.data.CourseSlots() {
		if cs.TeacherID() == teacherID {
			out = append(out, *cs)
		}
	}
	return out, nil
}

// GetClassroomSlots 取该教室现有的全部排期。
func (r *SlotChangeRepo) GetClassroomSlots(classroomID string) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range r.data.CourseSlots() {
		if cs.ClassroomID() == classroomID {
			out = append(out, *cs)
		}
	}
	return out, nil
}

// GetOtherSlotChanges 取除该变更单以外的全部换课记录。
func (r *SlotChangeRepo) GetOtherSlotChanges(id int64) ([]*courseSlotChange.CourseSlotChange, error) {
	all := r.data.CourseSlotChanges()

	out := make([]*courseSlotChange.CourseSlotChange, 0, len(all))
	for _, item := range all {
		if item.ID() != id {
			out = append(out, item)
		}
	}
	return out, nil
}
