package venue

import (
	"github.com/google/wire"

	"cqrs/internal/core/venue/biz/classroom"
	"cqrs/internal/core/venue/data"
	"cqrs/internal/core/venue/server"
	"cqrs/internal/core/venue/service"
)

// Servers is an alias for server.Servers.
type Servers = server.Servers

// ProviderSet is venue providers.
var ProviderSet = wire.NewSet(
	// biz
	classroom.NewClassroomUsecase,
	// data
	data.NewData,
	data.NewClassroomRepo,
	// service
	service.NewClassroomService,
	// server
	server.NewServers,
	server.NewHTTPServer,
	server.NewGRPCServer,
)
