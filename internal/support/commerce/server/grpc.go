package server

import (
	"cqrs/internal/conf"
	"cqrs/internal/support/commerce/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
)

// GRPCServer is a wrapper for commerce gRPC server.
type GRPCServer struct {
	*grpc.Server
}

// NewGRPCServer creates a new gRPC server for commerce context.
func NewGRPCServer(c *conf.Server, payment *service.PaymentService, discount *service.DiscountService) *GRPCServer {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	// TODO: Register gRPC services
	return &GRPCServer{Server: srv}
}
