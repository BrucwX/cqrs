package teaching

import (
	"github.com/google/wire"

	"cqrs/internal/core/teaching/ports"
	"cqrs/internal/core/teaching/service"
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
