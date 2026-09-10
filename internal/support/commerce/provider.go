package commerce

import (
	"github.com/google/wire"

	"cqrs/internal/support/commerce/adapters/memory"
	memorycmd "cqrs/internal/support/commerce/adapters/memory/command"
	memoryquery "cqrs/internal/support/commerce/adapters/memory/query"
	"cqrs/internal/support/commerce/app/command"
	"cqrs/internal/support/commerce/ports"
	"cqrs/internal/support/commerce/service"
)

// ProviderSet is commerce providers.
var ProviderSet = wire.NewSet(
	// adapter
	memory.NewData,
	// query adapters
	memoryquery.NewDiscountQuery,
	memoryquery.NewPaymentQuery,
	memoryquery.NewEnrollmentQuery,
	memoryquery.NewOrderQuery,
	// command adapters
	memorycmd.NewDiscountCommand,
	memorycmd.NewPaymentCommand,
	memorycmd.NewEnrollmentCommand,
	memorycmd.NewOrderCommand,
	// app
	command.NewPaymentUsecase,
	command.NewDiscountUsecase,
	command.NewEnrollmentUsecase,
	command.NewOrderUsecase,
	// service
	service.NewPaymentService,
	service.NewDiscountService,
	// ports
	ports.NewServers,
	ports.NewHTTPServer,
	ports.NewGRPCServer,
)
