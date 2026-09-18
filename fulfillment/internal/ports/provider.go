package ports

import "github.com/google/wire"

// ProviderSet is fulfillment ports providers.
var ProviderSet = wire.NewSet(
	NewServers,
	NewHTTPServer,
	NewGRPCServer,
)
