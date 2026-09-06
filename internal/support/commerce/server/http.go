package server

import (
	"cqrs/internal/conf"
	"cqrs/internal/support/commerce/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
)

// HTTPServer is a wrapper for commerce HTTP server.
type HTTPServer struct {
	*http.Server
}

// NewHTTPServer creates a new HTTP server for commerce context.
func NewHTTPServer(c *conf.Server, payment *service.PaymentService, discount *service.DiscountService) *HTTPServer {
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
