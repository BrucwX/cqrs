package course_scheduling

import (
	"github.com/google/wire"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	querymemory "cqrs/internal/core/course_scheduling/adapters/query/memory/implement"
	"cqrs/internal/core/course_scheduling/domain/service/command"
	"cqrs/internal/core/course_scheduling/domain/service/query"
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
	// domain - command handlers
	command.ProviderSet,
	// domain - query handlers
	query.ProviderSet,
	// service
	service.NewTeacherService,
	service.NewStudentService,
	service.NewCourseService,
	// ports
	ports.NewServers,
	ports.NewHTTPServer,
	ports.NewGRPCServer,
)
