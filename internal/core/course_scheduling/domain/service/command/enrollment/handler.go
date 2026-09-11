package enrollment

import (
	"cqrs/internal/core/course_scheduling/domain/repo/command"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// Handler 课程注册命令处理器
//
// 选课的规则在 checkEnrollment 里实现，按需从 enroll 取数据。
type Handler struct {
	EnrollmentCmd   command.CourseEnrollmentCommand
	EnrollmentQuery query.CourseEnrollmentQuery
	StudentQuery    query.StudentQuery
	CourseQuery     query.CourseQuery
	enroll          command.EnrollRepo
}

// NewHandler 创建课程注册命令处理器
func NewHandler(
	ec command.CourseEnrollmentCommand,
	eq query.CourseEnrollmentQuery,
	sq query.StudentQuery,
	cq query.CourseQuery,
	er command.EnrollRepo,
) *Handler {
	return &Handler{
		EnrollmentCmd:   ec,
		EnrollmentQuery: eq,
		StudentQuery:    sq,
		CourseQuery:     cq,
		enroll:          er,
	}
}
