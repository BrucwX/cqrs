package query

import "github.com/google/wire"

// ProviderSet is course_scheduling memory query providers.
var ProviderSet = wire.NewSet(
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
