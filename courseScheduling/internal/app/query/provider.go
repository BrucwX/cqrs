package query

import (
	"github.com/google/wire"

	"cqrs/internal/app/query/absence"
	"cqrs/internal/app/query/classroom"
	"cqrs/internal/app/query/course"
	"cqrs/internal/app/query/courseSlot"
	"cqrs/internal/app/query/courseSlotChange"
	"cqrs/internal/app/query/enrollment"
	"cqrs/internal/app/query/makeup"
	"cqrs/internal/app/query/qualification"
	"cqrs/internal/app/query/student"
	"cqrs/internal/app/query/teacher"
)

// ProviderSet query handler 依赖注入
var ProviderSet = wire.NewSet(
	absence.NewHandler,
	classroom.NewHandler,
	course.NewHandler,
	courseSlot.NewHandler,
	courseSlotChange.NewHandler,
	enrollment.NewHandler,
	makeup.NewHandler,
	qualification.NewHandler,
	student.NewHandler,
	teacher.NewHandler,
)
