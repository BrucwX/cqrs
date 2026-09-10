package memory

import "github.com/google/wire"

// ProviderSet is commerce memory adapter providers.
var ProviderSet = wire.NewSet(NewData)
