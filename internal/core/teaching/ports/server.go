package ports

import "github.com/go-kratos/kratos/v3/transport"

// Servers holds all servers for the teaching bounded context.
type Servers struct {
	HTTP transport.Server
	GRPC transport.Server
}

// NewServers creates all servers for teaching context.
func NewServers(http *HTTPServer, grpc *GRPCServer) *Servers {
	return &Servers{
		HTTP: http,
		GRPC: grpc,
	}
}
