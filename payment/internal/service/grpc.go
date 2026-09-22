package service

import (
	"payment/internal/conf"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
)

// GRPCServer 是 payment 上下文的 gRPC 服务器。
type GRPCServer struct {
	*grpc.Server
}

// NewGRPCServer 创建 payment 的 gRPC 服务器。
//
// 加资源时：把本包的 NewXxxService 加进参数列表，并在下面用
// xxxv1.RegisterXxxServiceServer(srv, svc) 注册。
func NewGRPCServer(c *conf.Server) *GRPCServer {
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
	return &GRPCServer{Server: srv}
}
