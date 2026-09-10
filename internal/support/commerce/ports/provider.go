package ports

import "github.com/google/wire"

// ProviderSet is commerce ports providers.
var ProviderSet = wire.NewSet(NewServers)
