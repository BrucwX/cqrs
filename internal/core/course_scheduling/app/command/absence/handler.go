package absence

import (
	appcommand "cqrs/internal/core/course_scheduling/app/command"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// Handler 缺勤记录命令处理器
type Handler struct {
	AbsenceCmd command.AbsenceRecordCommand
	tx         appcommand.Transaction
}

// NewHandler 创建缺勤记录命令处理器
func NewHandler(ac command.AbsenceRecordCommand, tx appcommand.Transaction) *Handler {
	return &Handler{
		AbsenceCmd: ac,
		tx:         tx,
	}
}
