package courseSlotChange

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 课表变更命令处理器
type Handler struct {
	ChangeCmd command.CourseSlotChangeCommand
}

// NewHandler 创建课表变更命令处理器
func NewHandler(cc command.CourseSlotChangeCommand) *Handler {
	return &Handler{ChangeCmd: cc}
}
