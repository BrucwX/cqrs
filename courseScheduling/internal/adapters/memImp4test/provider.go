package memImp4test

import "github.com/google/wire"

// ProviderSet is course_scheduling memory adapter providers.
var ProviderSet = wire.NewSet(NewData)
