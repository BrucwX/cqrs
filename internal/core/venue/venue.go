package venue

import (
	"github.com/google/wire"

	"github.com/go-kratos/kratos-layout/internal/core/venue/biz/classroom"
	"github.com/go-kratos/kratos-layout/internal/core/venue/data"
	"github.com/go-kratos/kratos-layout/internal/core/venue/server"
	"github.com/go-kratos/kratos-layout/internal/core/venue/service"
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
