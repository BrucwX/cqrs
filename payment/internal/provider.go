// Package payment 是 payment 界限上下文的组装根。
//
// 各层的构造函数都在自己的包里声明 ProviderSet，这里只把它们串起来：
//
//	adapters/command/Impl  写侧存储：MySQL 连接池 + 事务 + 仓库实现载体
//	adapters/query/Impl    读侧存储：读库连接池 + 查询实现载体
//	app/command            写侧用例（handler）+ 命令侧领域判定服务
//	app/query              读侧用例（handler）
//	domain/service         跨聚合根的领域规则（写路径判定）
//	service                传输层：HTTP / gRPC 服务器 + proto <-> 用例入参的转换
//
// ports（端口 / 接口）不在这张表里 —— 接口不需要谁“提供”，绑定写在 adapters 里。
package payment

import (
	"github.com/google/wire"

	commandimpl "payment/internal/adapters/command/Impl"
	queryimpl "payment/internal/adapters/query/Impl"
	appcommand "payment/internal/app/command"
	appquery "payment/internal/app/query"
	domainservice "payment/internal/domain/service"
	"payment/internal/service"
)

// ProviderSet 是 payment 上下文的全部 provider：只做组装，细节在各自包里。
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
	// service - 传输层（HTTP / gRPC 服务器 + 各 service）
	service.ProviderSet,
)
