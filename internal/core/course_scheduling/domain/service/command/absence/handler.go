package absence

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 缺勤记录命令处理器
type Handler struct {
	AbsenceCmd command.AbsenceRecordCommand
}

// NewHandler 创建缺勤记录命令处理器
func NewHandler(ac command.AbsenceRecordCommand) *Handler {
	return &Handler{AbsenceCmd: ac}
}
