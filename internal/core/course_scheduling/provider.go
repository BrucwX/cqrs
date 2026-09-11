package course_scheduling

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/memory/command"
	memoryquery "cqrs/internal/core/course_scheduling/adapters/memory/query"
	"cqrs/internal/core/course_scheduling/domain/service/command"
	"cqrs/internal/core/course_scheduling/domain/service/query"
	"cqrs/internal/core/course_scheduling/ports"
	"cqrs/internal/core/course_scheduling/service"
)

// ProviderSet is course_scheduling providers.
var ProviderSet = wire.NewSet(
	// adapters - data
	memory.NewData,
	// adapters - command (接口实现)
	memorycmd.ProviderSet,
	// adapters - query (接口实现)
	memoryquery.ProviderSet,
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
