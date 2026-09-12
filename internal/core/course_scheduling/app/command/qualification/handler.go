package qualification

import (
	appcommand "cqrs/internal/core/course_scheduling/app/command"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
	"cqrs/internal/core/course_scheduling/domain/service/qualificationCheck"
)

// Handler 授课资质命令处理器
//
// 发证的准入规则不在本地 —— 交给领域服务
// qualificationCheck.Service.CheckTeacherGrantable（讲师得先以学员身份修完该类型下
// 的一门课），本处理器只留「不够格就拦下」。
type Handler struct {
	QualificationCmd command.QualificationCommand
	qualify          *qualificationCheck.Service
	tx               appcommand.Transaction
}

// NewHandler 创建授课资质命令处理器
func NewHandler(qc command.QualificationCommand, qualify *qualificationCheck.Service, tx appcommand.Transaction) *Handler {
	return &Handler{
		QualificationCmd: qc,
		qualify:          qualify,
		tx:               tx,
	}
}
