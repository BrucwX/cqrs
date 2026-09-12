package implement

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/repoImp"
)

var ProviderSet = wire.NewSet(
	repoImp.NewAbsenceRecordQuery,
	repoImp.NewClassroomQuery,
	repoImp.NewCourseEnrollmentQuery,
	repoImp.NewCourseQuery,
	repoImp.NewCourseSlotChangeQuery,
	repoImp.NewCourseSlotQuery,
	repoImp.NewQualificationQuery,
	repoImp.NewStudentMakeupQuery,
	repoImp.NewStudentQuery,
	repoImp.NewTeacherQuery,
)
