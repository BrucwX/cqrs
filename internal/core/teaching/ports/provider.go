package ports

import "github.com/google/wire"

// ProviderSet is teaching ports providers.
var ProviderSet = wire.NewSet(NewServers)
