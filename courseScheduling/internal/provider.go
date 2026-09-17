package course_scheduling

import (
	"github.com/google/wire"

	commandimpl "cqrs/internal/adapters/command/Impl"
	queryimpl "cqrs/internal/adapters/query/Impl"
	commandproviders "cqrs/internal/app/command"
	"cqrs/internal/app/query"
	"cqrs/internal/ports"
	"cqrs/internal/service"
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
