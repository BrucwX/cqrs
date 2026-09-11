package command

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/domain/service/command/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/service/command/enrollment"
)

// ProviderSet command handler 依赖注入
var ProviderSet = wire.NewSet(
	enrollment.NewHandler,
	courseSlot.NewHandler,
)
