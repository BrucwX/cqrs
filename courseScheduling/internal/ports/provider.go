package ports

import "github.com/google/wire"

// ProviderSet is course_scheduling ports providers.
var ProviderSet = wire.NewSet(
	NewServers,
	NewHTTPServer,
	NewGRPCServer,
)
