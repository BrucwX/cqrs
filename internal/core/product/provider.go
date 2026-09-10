package product

import (
	"github.com/google/wire"

	"cqrs/internal/core/product/adapters/memory"
	memorycmd "cqrs/internal/core/product/adapters/memory/command"
	memoryquery "cqrs/internal/core/product/adapters/memory/query"
	"cqrs/internal/core/product/app/command"
	"cqrs/internal/core/product/ports"
	"cqrs/internal/core/product/service"
)

// ProviderSet is product providers.
var ProviderSet = wire.NewSet(
	// adapter
	memory.NewData,
	// query adapters
	memoryquery.NewProductQuery,
	// command adapters
	memorycmd.NewProductCommand,
	// app
	command.NewProductUsecase,
	// service
	service.NewProductService,
	// ports
	ports.NewServers,
	ports.NewHTTPServer,
	ports.NewGRPCServer,
)
