package teaching

import (
	"github.com/google/wire"

	"cqrs/internal/core/teaching/biz"
	"cqrs/internal/core/teaching/biz/course"
	"cqrs/internal/core/teaching/biz/student"
	"cqrs/internal/core/teaching/biz/teacher"
	"cqrs/internal/core/teaching/data"
	"cqrs/internal/core/teaching/server"
	"cqrs/internal/core/teaching/service"
)

// Servers is an alias for server.Servers.
type Servers = server.Servers

// ProviderSet is teaching providers.
var ProviderSet = wire.NewSet(
	// biz
	teacher.NewTeacherUsecase,
	student.NewStudentUsecase,
	course.NewCourseUsecase,
	biz.NewCourseFinder,
	biz.NewScheduleValidator,
	// data
	data.NewData,
	data.NewTeacherRepo,
	data.NewStudentRepo,
	data.NewCourseRepo,
	// service
	service.NewTeacherService,
	service.NewStudentService,
	service.NewCourseService,
	// server
	server.NewServers,
	server.NewHTTPServer,
	server.NewGRPCServer,
)
