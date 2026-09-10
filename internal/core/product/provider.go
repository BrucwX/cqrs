package product

import (
	"github.com/google/wire"

	"cqrs/internal/core/product/adapters/memory"
	memorycmd "cqrs/internal/core/product/adapters/memory/command"
	"cqrs/internal/core/product/app/command"
	"cqrs/internal/core/product/ports"
	"cqrs/internal/core/product/service"
)

// ProviderSet is product providers.
var ProviderSet = wire.NewSet(
	// adapter
	memory.NewData,
	memorycmd.NewProductRepo,
	// app
	command.NewProductUsecase,
	// service
	service.NewProductService,
	// ports
	ports.NewServers,
	ports.NewHTTPServer,
	ports.NewGRPCServer,
)
