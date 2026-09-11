package courseSlot

import (
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// Handler 课表槽位命令处理器
//
// 三个 AssignXxx 的规则都在本处理器里实现（checkTeacher / checkCourse /
// checkClassroom），各自按需从对应的 repo 取数据。
type Handler struct {
	SlotCmd      command.CourseSlotCommand
	assignTea    command.AssignTeacherRepo
	assignCourse command.AssignCourseRepo
	assignClass  command.AssignClassroomRepo
}

// NewHandler 创建课表槽位命令处理器
func NewHandler(
	sc command.CourseSlotCommand,
	at command.AssignTeacherRepo,
	ac command.AssignCourseRepo,
	acl command.AssignClassroomRepo,
) *Handler {
	return &Handler{
		SlotCmd:      sc,
		assignTea:    at,
		assignCourse: ac,
		assignClass:  acl,
	}
}
