package teaching

import (
	"github.com/google/wire"

	"github.com/go-kratos/kratos-layout/internal/core/teaching/biz"
	"github.com/go-kratos/kratos-layout/internal/core/teaching/biz/course"
	"github.com/go-kratos/kratos-layout/internal/core/teaching/biz/student"
	"github.com/go-kratos/kratos-layout/internal/core/teaching/biz/teacher"
	"github.com/go-kratos/kratos-layout/internal/core/teaching/data"
	"github.com/go-kratos/kratos-layout/internal/core/teaching/server"
	"github.com/go-kratos/kratos-layout/internal/core/teaching/service"
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
