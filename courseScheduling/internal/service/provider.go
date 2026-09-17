package service

import "github.com/google/wire"

// ProviderSet is course_scheduling service providers.
var ProviderSet = wire.NewSet(
	NewTeacherService,
	NewStudentService,
	NewClassroomService,
	NewCourseService,
	NewCourseSlotService,
	NewCourseSlotChangeService,
	NewEnrollmentService,
	NewAbsenceService,
	NewMakeupService,
	NewQualificationService,
)
