package qualification

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 授课资质命令处理器
type Handler struct {
	QualificationCmd command.QualificationCommand
}

// NewHandler 创建授课资质命令处理器
func NewHandler(qc command.QualificationCommand) *Handler {
	return &Handler{QualificationCmd: qc}
}
