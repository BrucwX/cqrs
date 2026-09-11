package courseSlot

import (
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// Handler 课表槽位命令处理器
type Handler struct {
	SlotCmd   command.CourseSlotCommand
	assignTea command.AssignTeacherRepo
}

// NewHandler 创建课表槽位命令处理器
func NewHandler(sc command.CourseSlotCommand, at command.AssignTeacherRepo) *Handler {
	return &Handler{
		SlotCmd:   sc,
		assignTea: at,
	}
}
