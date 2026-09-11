package courseSlot

import (
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// Handler 课表槽位命令处理器
type Handler struct {
	SlotCmd command.CourseSlotCommand
}

// NewHandler 创建课表槽位命令处理器
func NewHandler(sc command.CourseSlotCommand) *Handler {
	return &Handler{SlotCmd: sc}
}
