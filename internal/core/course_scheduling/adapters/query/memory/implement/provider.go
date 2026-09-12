package implement

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/adapters/query/memory"
)

// ProviderSet 是读侧内存实现的 provider 集合。
//
// 与 MySQL 那套（adapters/query/mysql/implement.ProviderSet）一一对应：
// 二者实现同一批 domain/repo/query 接口，可以整体互换 —— 换实现不动上层代码。
//
// memory.NewData 把共享的内存存储收窄成只读面，下面这些仓库拿到的就是这个窄面，
// 所以它们编译期就调不到 Save* / Delete*。
var ProviderSet = wire.NewSet(
	memory.NewData,
	NewAbsenceRecordQuery,
	NewClassroomQuery,
	NewCourseQuery,
	NewCourseSlotQuery,
	NewCourseSlotChangeQuery,
	NewCourseEnrollmentQuery,
	NewStudentMakeupQuery,
	NewQualificationQuery,
	NewStudentQuery,
	NewTeacherQuery,
)
