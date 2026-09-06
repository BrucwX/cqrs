package product

import (
	"github.com/google/wire"

	"cqrs/internal/core/product/biz/product"
	"cqrs/internal/core/product/data"
	"cqrs/internal/core/product/server"
	"cqrs/internal/core/product/service"
)

// Servers is an alias for server.Servers.
type Servers = server.Servers

// ProviderSet is product providers.
var ProviderSet = wire.NewSet(
	// biz
	product.NewProductUsecase,
	// data
	data.NewData,
	data.NewProductRepo,
	// service
	service.NewProductService,
	// server
	server.NewServers,
	server.NewHTTPServer,
	server.NewGRPCServer,
)
