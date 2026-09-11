package makeup

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 补课预约命令处理器
type Handler struct {
	MakeupCmd command.StudentMakeupCommand
}

// NewHandler 创建补课预约命令处理器
func NewHandler(mc command.StudentMakeupCommand) *Handler {
	return &Handler{MakeupCmd: mc}
}
