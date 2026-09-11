package command

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/domain/service/command/absence"
	"cqrs/internal/core/course_scheduling/domain/service/command/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/service/command/courseSlotChange"
	"cqrs/internal/core/course_scheduling/domain/service/command/enrollment"
	"cqrs/internal/core/course_scheduling/domain/service/command/makeup"
)

// ProviderSet command handler 依赖注入
var ProviderSet = wire.NewSet(
	absence.NewHandler,
	courseSlot.NewHandler,
	courseSlotChange.NewHandler,
	enrollment.NewHandler,
	makeup.NewHandler,
)
