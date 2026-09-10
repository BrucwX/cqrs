package course_scheduling

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/ports"
	"cqrs/internal/core/course_scheduling/service"
)

// ProviderSet is teaching providers.
var ProviderSet = wire.NewSet(
	// service
	service.NewTeacherService,
	service.NewStudentService,
	service.NewCourseService,
	// ports
	ports.NewServers,
	ports.NewHTTPServer,
	ports.NewGRPCServer,
)
