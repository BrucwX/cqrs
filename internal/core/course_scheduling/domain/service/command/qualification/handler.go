package qualification

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Handler 授课资质命令处理器
//
// 授证的准入规则在 qualify.go 的 checkFinished 里实现，按需从 qualify 取数据。
type Handler struct {
	QualificationCmd command.QualificationCommand
	qualify          command.QualifyRepo
}

// NewHandler 创建授课资质命令处理器
func NewHandler(qc command.QualificationCommand, qr command.QualifyRepo) *Handler {
	return &Handler{QualificationCmd: qc, qualify: qr}
}
