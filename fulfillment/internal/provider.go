// Package fulfillment 是 fulfillment 界限上下文的组装根。
//
// 各层的构造函数都在自己的包里声明 ProviderSet，这里只把它们串起来：
//
//	adapters/command/Impl  写侧存储：MySQL 连接池 + 事务 + 仓库实现载体
//	adapters/query/Impl    读侧存储：读库连接池 + 查询实现载体
//	app/command            写侧用例（handler）+ 命令侧领域判定服务
//	app/query              读侧用例（handler）
//	domain/service         跨聚合根的领域规则（写路径判定）
//	service                传输层：proto <-> 用例入参的转换
//	ports                  HTTP / gRPC 服务器
package fulfillment

import (
	"github.com/google/wire"

	commandimpl "fulfillment/internal/adapters/command/Impl"
	queryimpl "fulfillment/internal/adapters/query/Impl"
	appcommand "fulfillment/internal/app/command"
	appquery "fulfillment/internal/app/query"
	domainservice "fulfillment/internal/domain/service"
	"fulfillment/internal/ports"
	"fulfillment/internal/service"
)

// ProviderSet 是 fulfillment 上下文的全部 provider：只做组装，细节在各自包里。
var ProviderSet = wire.NewSet(
	// adapters - 写侧
	commandimpl.ProviderSet,
	// adapters - 读侧
	queryimpl.ProviderSet,
	// app - command handlers（含 domain 判定服务，收在 app/command）
	appcommand.ProviderSet,
	// app - query handlers
	appquery.ProviderSet,
	// domain - 跨聚合规则
	domainservice.ProviderSet,
	// service
	service.ProviderSet,
	// ports
	ports.ProviderSet,
)
