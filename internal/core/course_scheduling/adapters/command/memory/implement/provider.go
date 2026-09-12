package implement

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
)

// ProviderSet 是写侧内存实现的 provider 集合。
//
// 与 MySQL 那套（adapters/command/mysql）一一对应：二者实现同一批
// domain/repo/command 接口，可以整体互换。
//
// memory.NewData 把共享的内存存储收窄成写侧面（能读也能写），下面这些仓库
// 拿到的就是这个窄面。
var ProviderSet = wire.NewSet(
	memory.NewData,
	NewAbsenceCommand,
	NewAssignClassroomRepo,
	NewAssignCourseRepo,
	NewAssignTeacherRepo,
	NewClassroomCommand,
	NewCourseCommand,
	NewCourseSlotChangeCommand,
	NewCourseSlotCommand,
	NewEnrollRepo,
	NewEnrollmentCommand,
	NewMakeupCommand,
	NewQualificationCommand,
	NewQualifyRepo,
	NewSlotChangeRepo,
	NewStudentCommand,
	NewTeacherCommand,
)
