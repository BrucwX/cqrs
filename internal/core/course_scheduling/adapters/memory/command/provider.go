package command

import "github.com/google/wire"

// ProviderSet is course_scheduling memory command providers.
var ProviderSet = wire.NewSet(
	NewAbsenceCommand,
	NewAssignClassroomRepo,
	NewAssignCourseRepo,
	NewAssignTeacherRepo,
	NewCourseSlotChangeCommand,
	NewCourseSlotCommand,
	NewEnrollRepo,
	NewEnrollmentCommand,
	NewMakeupCommand,
	NewSlotChangeRepo,
)
