package ports

import "github.com/google/wire"

// ProviderSet is product ports providers.
var ProviderSet = wire.NewSet(NewServers)
