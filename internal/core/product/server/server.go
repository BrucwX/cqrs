package server

import (
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/google/wire"
)

// Servers holds all servers for the product bounded context.
type Servers struct {
	HTTP transport.Server
	GRPC transport.Server
}

// ProviderSet is product server providers.
var ProviderSet = wire.NewSet(NewServers)

// NewServers creates all servers for product context.
func NewServers(http *HTTPServer, grpc *GRPCServer) *Servers {
	return &Servers{
		HTTP: http,
		GRPC: grpc,
	}
}
