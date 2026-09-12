package course

import (
	appcommand "cqrs/internal/core/course_scheduling/app/command"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// Handler 课程命令处理器
type Handler struct {
	CourseCmd command.CourseCommand
	tx        appcommand.Transaction
}

// NewHandler 创建课程命令处理器
func NewHandler(cc command.CourseCommand, tx appcommand.Transaction) *Handler {
	return &Handler{
		CourseCmd: cc,
		tx:        tx,
	}
}
