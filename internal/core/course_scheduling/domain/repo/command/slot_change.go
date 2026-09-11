package command

import (
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// SlotChangeRepo 是临时换课做冲突检查所需的数据来源。
//
// checkSlotChange 由命令服务自己实现，它需要「目标讲师的排期 / 目标教室的排期 /
// 其他换课记录」；仓库只负责按需提供这些数据，不参与规则判定。
type SlotChangeRepo interface {
	// GetTeacherSlots 取该讲师现有的全部排期
	GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error)
	// GetClassroomSlots 取该教室现有的全部排期
	GetClassroomSlots(classroomID string) (courseSlot.CourseSlots, error)
	// todo 这里不应该取这个，应该取 和 CourseSlotChange 有交集的换课记录
	// GetOtherSlotChanges 取除该变更单以外的全部换课记录
	GetOtherSlotChanges(id int64) ([]*courseSlotChange.CourseSlotChange, error)
}
