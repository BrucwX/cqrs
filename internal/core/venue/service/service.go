package service

import "github.com/google/wire"

// ProviderSet is venue service providers.
var ProviderSet = wire.NewSet(NewClassroomService)
