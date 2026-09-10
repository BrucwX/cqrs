package server

import (
	"cqrs/internal/core/product"
	"cqrs/internal/core/teaching"
	"cqrs/internal/support/commerce"

	"github.com/go-kratos/kratos/v3/transport"
)

// Servers collects all servers from bounded contexts.
type Servers struct {
	HTTP []transport.Server
	GRPC []transport.Server
}

// NewServers creates a Servers that aggregates all bounded context servers.
func NewServers(
	teaching *teaching.Servers,
	product *product.Servers,
	commerce *commerce.Servers,
) *Servers {
	return &Servers{
		HTTP: []transport.Server{
			teaching.HTTP,
			product.HTTP,
			commerce.HTTP,
		},
		GRPC: []transport.Server{
			teaching.GRPC,
			product.GRPC,
			commerce.GRPC,
		},
	}
}

// All returns all servers (HTTP + gRPC) for starting.
func (s *Servers) All() []transport.Server {
	var all []transport.Server
	all = append(all, s.HTTP...)
	all = append(all, s.GRPC...)
	return all
}
