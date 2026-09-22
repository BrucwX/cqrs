package service

import (
	"payment/internal/conf"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
)

// HTTPServer 是 payment 上下文的 HTTP 服务器。
type HTTPServer struct {
	*http.Server
}

// NewHTTPServer 创建 payment 的 HTTP 服务器。
//
// 加资源时：把本包的 NewXxxService 加进参数列表，并在下面用
// xxxv1.RegisterXxxServiceHTTPServer(srv, svc) 注册 —— 多个 service 各自注册到
// 同一个 mux，路由不会冲突。
func NewHTTPServer(c *conf.Server) *HTTPServer {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	return &HTTPServer{Server: srv}
}
