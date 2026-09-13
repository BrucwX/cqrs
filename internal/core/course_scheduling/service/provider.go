package service

import "github.com/google/wire"

// ProviderSet is teaching service providers.
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
