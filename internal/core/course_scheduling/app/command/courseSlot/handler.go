package courseSlot

import (
	appcommand "cqrs/internal/core/course_scheduling/app/command"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
	"cqrs/internal/core/course_scheduling/domain/service/classroomCapacity"
	"cqrs/internal/core/course_scheduling/domain/service/qualificationCheck"
	"cqrs/internal/core/course_scheduling/domain/service/scheduleConflict"
)

// Handler 课表槽位命令处理器
//
// 三个 AssignXxx 的规则：判定部分全交给领域服务 —— 时间冲突与课程存在性走
// scheduleConflict，资质走 qualificationCheck，教室容量走 classroomCapacity；
// 本处理器只留「有一条不过就拦下」。
type Handler struct {
	SlotCmd  command.CourseSlotCommand
	conflict *scheduleConflict.Service
	qualify  *qualificationCheck.Service
	capacity *classroomCapacity.Service
	tx       appcommand.Transaction
}

// NewHandler 创建课表槽位命令处理器
func NewHandler(
	sc command.CourseSlotCommand,
	conflict *scheduleConflict.Service,
	qualify *qualificationCheck.Service,
	capacity *classroomCapacity.Service,
	tx appcommand.Transaction,
) *Handler {
	return &Handler{
		SlotCmd:  sc,
		conflict: conflict,
		qualify:  qualify,
		capacity: capacity,
		tx:       tx,
	}
}
