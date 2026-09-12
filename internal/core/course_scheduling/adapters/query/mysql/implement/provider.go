package implement

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/imp"
)

var ProviderSet = wire.NewSet(
	imp.NewAbsenceRecordQuery,
	imp.NewClassroomQuery,
	imp.NewCourseEnrollmentQuery,
	imp.NewCourseQuery,
	imp.NewCourseSlotChangeQuery,
	imp.NewCourseSlotQuery,
	imp.NewQualificationQuery,
	imp.NewStudentMakeupQuery,
	imp.NewStudentQuery,
	imp.NewTeacherQuery,
)
