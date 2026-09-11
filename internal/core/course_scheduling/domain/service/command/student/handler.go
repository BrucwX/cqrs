package student

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 学员命令处理器
type Handler struct {
	StudentCmd command.StudentCommand
}

// NewHandler 创建学员命令处理器
func NewHandler(sc command.StudentCommand) *Handler {
	return &Handler{StudentCmd: sc}
}
