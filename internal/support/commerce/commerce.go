package commerce

import (
	"github.com/google/wire"

	"cqrs/internal/support/commerce/biz/discount"
	"cqrs/internal/support/commerce/biz/payment"
	"cqrs/internal/support/commerce/data"
	"cqrs/internal/support/commerce/server"
	"cqrs/internal/support/commerce/service"
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
