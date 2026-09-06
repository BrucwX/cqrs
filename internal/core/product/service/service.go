package service

import "github.com/google/wire"

// ProviderSet is product service providers.
var ProviderSet = wire.NewSet(NewProductService)
