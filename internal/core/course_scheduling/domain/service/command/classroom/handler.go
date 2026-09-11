package classroom

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 教室命令处理器
type Handler struct {
	ClassroomCmd command.ClassroomCommand
}

// NewHandler 创建教室命令处理器
func NewHandler(cc command.ClassroomCommand) *Handler {
	return &Handler{ClassroomCmd: cc}
}
