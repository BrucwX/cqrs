package courseSlotChange

import (
	appcommand "cqrs/internal/core/course_scheduling/app/command"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
	"cqrs/internal/core/course_scheduling/domain/service/scheduleConflict"
)

// Handler 课表变更命令处理器
//
// 换课的冲突判定不在本地 —— 交给领域服务
// scheduleConflict.Service.CheckSlotChange（目标讲师 / 教室在目标时段有没有被
// 占用），本处理器只留「有冲突就拦下」。
type Handler struct {
	ChangeCmd command.CourseSlotChangeCommand
	conflict  *scheduleConflict.Service
	tx        appcommand.Transaction
}

// NewHandler 创建课表变更命令处理器
func NewHandler(cc command.CourseSlotChangeCommand, conflict *scheduleConflict.Service, tx appcommand.Transaction) *Handler {
	return &Handler{
		ChangeCmd: cc,
		conflict:  conflict,
		tx:        tx,
	}
}
