package course

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 课程命令处理器
type Handler struct {
	CourseCmd command.CourseCommand
}

// NewHandler 创建课程命令处理器
func NewHandler(cc command.CourseCommand) *Handler {
	return &Handler{CourseCmd: cc}
}
