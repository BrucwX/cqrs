package server

import (
	"cqrs/internal/support/commerce"
	"cqrs/internal/core/product"
	"cqrs/internal/core/teaching"
	"cqrs/internal/core/venue"

	"github.com/go-kratos/kratos/v3/transport"
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServers)

// Servers collects all servers from bounded contexts.
type Servers struct {
	HTTP []transport.Server
	GRPC []transport.Server
}

// NewServers creates a Servers that aggregates all bounded context servers.
func NewServers(
	teaching *teaching.Servers,
	venue *venue.Servers,
	product *product.Servers,
	commerce *commerce.Servers,
) *Servers {
	return &Servers{
		HTTP: []transport.Server{
			teaching.HTTP,
			venue.HTTP,
			product.HTTP,
			commerce.HTTP,
		},
		GRPC: []transport.Server{
			teaching.GRPC,
			venue.GRPC,
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
