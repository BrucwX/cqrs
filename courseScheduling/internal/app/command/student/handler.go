package student

import (
	appcommand "courseScheduling/internal/app/command/appRepo"
	"courseScheduling/internal/domain/repo/command"
)

// Handler 学员命令处理器
type Handler struct {
	StudentCmd command.StudentCommand
	tx         appcommand.Transaction
}

// NewHandler 创建学员命令处理器
func NewHandler(sc command.StudentCommand, tx appcommand.Transaction) *Handler {
	return &Handler{
		StudentCmd: sc,
		tx:         tx,
	}
}
