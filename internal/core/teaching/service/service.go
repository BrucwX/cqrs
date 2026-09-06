package service

import "github.com/google/wire"

// ProviderSet is teaching service providers.
var ProviderSet = wire.NewSet(
	NewTeacherService,
	NewStudentService,
	NewCourseService,
)
