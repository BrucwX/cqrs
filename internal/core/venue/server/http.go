package server

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/core/venue/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
)

// HTTPServer is a wrapper for venue HTTP server.
type HTTPServer struct {
	*http.Server
}

// NewHTTPServer creates a new HTTP server for venue context.
func NewHTTPServer(c *conf.Server, classroom *service.ClassroomService) *HTTPServer {
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
	// TODO: Register HTTP routes
	return &HTTPServer{Server: srv}
}
