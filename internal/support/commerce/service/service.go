package service

import "github.com/google/wire"

// ProviderSet is commerce service providers.
var ProviderSet = wire.NewSet(
	NewPaymentService,
	NewDiscountService,
)
