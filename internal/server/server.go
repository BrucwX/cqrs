package server

import (
	"cqrs/internal/core/course_scheduling"

	"github.com/go-kratos/kratos/v3/transport"
)

// Servers collects all servers from bounded contexts.
type Servers struct {
	HTTP []transport.Server
	GRPC []transport.Server
}

// NewServers creates a Servers that aggregates all bounded context servers.
//
// 目前只有 course_scheduling 一个上下文在跑；product / commerce 的代码保留着，
// 只是没进依赖注入。等它们的 provider 补齐、重新挂回 wire.Build 之后，
// 再把它们的 *Servers 参数加回来。
func NewServers(courseScheduling *course_scheduling.Servers) *Servers {
	return &Servers{
		HTTP: []transport.Server{courseScheduling.HTTP},
		GRPC: []transport.Server{courseScheduling.GRPC},
	}
}

// All returns all servers (HTTP + gRPC) for starting.
func (s *Servers) All() []transport.Server {
	var all []transport.Server
	all = append(all, s.HTTP...)
	all = append(all, s.GRPC...)
	return all
}
