package classroom

import (
	appcommand "cqrs/internal/core/course_scheduling/app/command"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// Handler 教室命令处理器
type Handler struct {
	ClassroomCmd command.ClassroomCommand
	tx           appcommand.Transaction
}

// NewHandler 创建教室命令处理器
func NewHandler(cc command.ClassroomCommand, tx appcommand.Transaction) *Handler {
	return &Handler{
		ClassroomCmd: cc,
		tx:           tx,
	}
}
