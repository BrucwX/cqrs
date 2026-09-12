package implement

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/imp"
)

// ProviderSet 是写侧 MySQL 实现的 provider 集合。
//
// 与内存那套（adapters/command/memory/implement.ProviderSet）一一对应：
// 二者实现同一批 domain/repo/command 接口，可以整体互换。
var ProviderSet = wire.NewSet(
	imp.NewAbsenceRecordImp,
	imp.NewClassroomImp,
	imp.NewCourseEnrollmentImp,
	imp.NewCourseImp,
	imp.NewCourseSlotChangeImp,
	imp.NewCourseSlotImp,
	imp.NewCourseTypeImp,
	imp.NewQualificationImp,
	imp.NewStudentImp,
	imp.NewStudentMakeupImp,
	imp.NewTeacherImp,
)
