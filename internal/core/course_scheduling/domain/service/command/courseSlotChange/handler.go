package courseSlotChange

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 课表变更命令处理器
//
// 换课的规则在 checkSlotChange 里实现，按需从 slotChange 取数据。
type Handler struct {
	ChangeCmd  command.CourseSlotChangeCommand
	slotChange command.SlotChangeRepo
}

// NewHandler 创建课表变更命令处理器
func NewHandler(cc command.CourseSlotChangeCommand, sc command.SlotChangeRepo) *Handler {
	return &Handler{
		ChangeCmd:  cc,
		slotChange: sc,
	}
}
