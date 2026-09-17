package ports

import "github.com/google/wire"

// ProviderSet is payment ports providers.
var ProviderSet = wire.NewSet(
	NewServers,
	NewHTTPServer,
	NewGRPCServer,
)
