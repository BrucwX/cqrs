package memory

import "github.com/google/wire"

// ProviderSet is product memory adapter providers.
var ProviderSet = wire.NewSet(NewData)
