package course_scheduling

import (
	"github.com/google/wire"

	commandimpl "cqrs/internal/core/course_scheduling/adapters/command/Impl"
	queryimpl "cqrs/internal/core/course_scheduling/adapters/query/Impl"
	commandproviders "cqrs/internal/core/course_scheduling/app/command"
	"cqrs/internal/core/course_scheduling/app/query"
	"cqrs/internal/core/course_scheduling/ports"
	"cqrs/internal/core/course_scheduling/service"
)

// ProviderSet is course_scheduling providers：只做组装，细节在各自包里。
//
//   - 写侧：commandimpl.ProviderSet（MySQL 存储 + CommandImpl 载体 + 接口绑定）
//   - 读侧：queryimpl.ProviderSet（MySQL 存储 + QueryImpl 载体 + 接口绑定）
//   - app/query、app/command（含领域判定服务）、service、ports 各带自己的 ProviderSet
var ProviderSet = wire.NewSet(
	// adapters - 写侧
	commandimpl.ProviderSet,
	// adapters - 读侧
	queryimpl.ProviderSet,
	// app - command handlers（含 domain 判定服务，收在 app/command）
	commandproviders.ProviderSet,
	// app - query handlers
	query.ProviderSet,
	// service
	service.ProviderSet,
	// ports
	ports.ProviderSet,
)
