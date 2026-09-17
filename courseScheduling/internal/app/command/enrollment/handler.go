package enrollment

import (
	appcommand "courseScheduling/internal/app/command/appRepo"
	"courseScheduling/internal/domain/repo/command"
	"courseScheduling/internal/domain/repo/query"
	"courseScheduling/internal/domain/service/scheduleConflict"
)

// Handler 课程注册命令处理器
//
// 选课的准入规则不在本地 —— 交给领域服务
// scheduleConflict.Service.CheckEnrollment（窗口 / 容量 / 时间冲突），
// 本处理器只留「不能选就拦下」。
type Handler struct {
	EnrollmentCmd   command.CourseEnrollmentCommand
	EnrollmentQuery query.CourseEnrollmentQuery
	StudentQuery    query.StudentQuery
	CourseQuery     query.CourseQuery
	conflict        *scheduleConflict.Service
	tx              appcommand.Transaction
}

// NewHandler 创建课程注册命令处理器
func NewHandler(
	ec command.CourseEnrollmentCommand,
	eq query.CourseEnrollmentQuery,
	sq query.StudentQuery,
	cq query.CourseQuery,
	conflict *scheduleConflict.Service,
	tx appcommand.Transaction,
) *Handler {
	return &Handler{
		EnrollmentCmd:   ec,
		EnrollmentQuery: eq,
		StudentQuery:    sq,
		CourseQuery:     cq,
		conflict:        conflict,
		tx:              tx,
	}
}
