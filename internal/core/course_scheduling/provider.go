package course_scheduling

import (
	"github.com/google/wire"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	querymemory "cqrs/internal/core/course_scheduling/adapters/query/memory/implement"
	"cqrs/internal/core/course_scheduling/app/command/absence"
	"cqrs/internal/core/course_scheduling/app/command/classroom"
	"cqrs/internal/core/course_scheduling/app/command/course"
	"cqrs/internal/core/course_scheduling/app/command/courseSlot"
	"cqrs/internal/core/course_scheduling/app/command/courseSlotChange"
	"cqrs/internal/core/course_scheduling/app/command/enrollment"
	"cqrs/internal/core/course_scheduling/app/command/makeup"
	"cqrs/internal/core/course_scheduling/app/command/qualification"
	"cqrs/internal/core/course_scheduling/app/command/student"
	"cqrs/internal/core/course_scheduling/app/command/teacher"
	"cqrs/internal/core/course_scheduling/app/query"
	domainservice "cqrs/internal/core/course_scheduling/domain/service"
	"cqrs/internal/core/course_scheduling/ports"
	"cqrs/internal/core/course_scheduling/service"
)

// ProviderSet is course_scheduling providers.
//
// 存储按读写两侧分别接线：写侧（command）与读侧（query）各自挑一个实现，
// 两边可以不是同一套。内存实现多一个共享的 store（memorystore），
// 它既不属于读侧也不属于写侧，两侧各自把它收窄成自己需要的那一面。
var ProviderSet = wire.NewSet(
	// adapters - 内存共享存储（读写两侧共用一个实例）
	memorystore.ProviderSet,
	// adapters - command (接口实现)
	commandmemory.ProviderSet,
	// adapters - query (接口实现)
	querymemory.ProviderSet,
	// app - command handlers
	//
	// 逐个列在这里，而不是 app/command.ProviderSet：那个包现在只放跨 handler 的
	// 接口（transaction.go），一旦它 import 子包，子包就没法反过来 import 它了。
	absence.NewHandler,
	classroom.NewHandler,
	course.NewHandler,
	courseSlot.NewHandler,
	courseSlotChange.NewHandler,
	makeup.NewHandler,
	enrollment.NewHandler,
	qualification.NewHandler,
	student.NewHandler,
	teacher.NewHandler,
	// app - query handlers
	query.ProviderSet,
	// domain - 判定服务
	domainservice.ProviderSet,
	// service
	service.NewTeacherService,
	service.NewStudentService,
	service.NewCourseService,
	// ports
	ports.NewServers,
	ports.NewHTTPServer,
	ports.NewGRPCServer,
)
