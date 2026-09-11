package courseSlot

import (
	"cqrs/internal/core/course_scheduling/domain/repo/command"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// Handler 课表槽位命令处理器
type Handler struct {
	SlotCmd  command.CourseSlotCommand
	SlotQuery query.CourseSlotQuery
}

// NewHandler 创建课表槽位命令处理器
func NewHandler(sc command.CourseSlotCommand, sq query.CourseSlotQuery) *Handler {
	return &Handler{
		SlotCmd:  sc,
		SlotQuery: sq,
	}
}
