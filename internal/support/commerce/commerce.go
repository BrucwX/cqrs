package commerce

import (
	"github.com/google/wire"

	"github.com/go-kratos/kratos-layout/internal/support/commerce/biz/discount"
	"github.com/go-kratos/kratos-layout/internal/support/commerce/biz/payment"
	"github.com/go-kratos/kratos-layout/internal/support/commerce/data"
	"github.com/go-kratos/kratos-layout/internal/support/commerce/server"
	"github.com/go-kratos/kratos-layout/internal/support/commerce/service"
)

// Servers is an alias for server.Servers.
type Servers = server.Servers

// ProviderSet is commerce providers.
var ProviderSet = wire.NewSet(
	// biz
	payment.NewPaymentUsecase,
	discount.NewDiscountUsecase,
	// data
	data.NewData,
	data.NewPaymentRepo,
	data.NewDiscountRepo,
	// service
	service.NewPaymentService,
	service.NewDiscountService,
	// server
	server.NewServers,
	server.NewHTTPServer,
	server.NewGRPCServer,
)
