package makeup

import (
	appcommand "cqrs/internal/core/course_scheduling/app/command"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
	"cqrs/internal/core/course_scheduling/domain/service/classroomCapacity"
)

// Handler 补课预约命令处理器
//
// 预约的准入规则不在本地 —— 交给领域服务
// classroomCapacity.Service.CheckMakeup（目标那节课的教室装不装得下），
// 本处理器只留「装不下就拦下」。
type Handler struct {
	MakeupCmd command.StudentMakeupCommand
	capacity  *classroomCapacity.Service
	tx        appcommand.Transaction
}

// NewHandler 创建补课预约命令处理器
func NewHandler(mc command.StudentMakeupCommand, capacity *classroomCapacity.Service, tx appcommand.Transaction) *Handler {
	return &Handler{
		MakeupCmd: mc,
		capacity:  capacity,
		tx:        tx,
	}
}
