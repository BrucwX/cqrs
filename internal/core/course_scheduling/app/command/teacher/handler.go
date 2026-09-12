package teacher

import (
	appcommand "cqrs/internal/core/course_scheduling/app/command"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// Handler 讲师命令处理器
type Handler struct {
	TeacherCmd command.TeacherCommand
	tx         appcommand.Transaction
}

// NewHandler 创建讲师命令处理器
func NewHandler(tc command.TeacherCommand, tx appcommand.Transaction) *Handler {
	return &Handler{
		TeacherCmd: tc,
		tx:         tx,
	}
}
