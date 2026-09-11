package teacher

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 讲师命令处理器
type Handler struct {
	TeacherCmd command.TeacherCommand
}

// NewHandler 创建讲师命令处理器
func NewHandler(tc command.TeacherCommand) *Handler {
	return &Handler{TeacherCmd: tc}
}
