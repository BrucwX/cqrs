package commerce

import (
	"github.com/google/wire"

	"cqrs/internal/support/commerce/adapters/memory"
	memorycmd "cqrs/internal/support/commerce/adapters/memory/command"
	"cqrs/internal/support/commerce/app/command"
	"cqrs/internal/support/commerce/ports"
	"cqrs/internal/support/commerce/service"
)

// ProviderSet is commerce providers.
var ProviderSet = wire.NewSet(
	// adapter
	memory.NewData,
	memorycmd.NewPaymentRepo,
	memorycmd.NewDiscountRepo,
	// app
	command.NewPaymentUsecase,
	command.NewDiscountUsecase,
	// service
	service.NewPaymentService,
	service.NewDiscountService,
	// ports
	ports.NewServers,
	ports.NewHTTPServer,
	ports.NewGRPCServer,
)
