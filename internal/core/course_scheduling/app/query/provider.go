package query

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/app/query/absence"
	"cqrs/internal/core/course_scheduling/app/query/classroom"
	"cqrs/internal/core/course_scheduling/app/query/course"
	"cqrs/internal/core/course_scheduling/app/query/courseSlot"
	"cqrs/internal/core/course_scheduling/app/query/courseSlotChange"
	"cqrs/internal/core/course_scheduling/app/query/enrollment"
	"cqrs/internal/core/course_scheduling/app/query/makeup"
	"cqrs/internal/core/course_scheduling/app/query/qualification"
	"cqrs/internal/core/course_scheduling/app/query/student"
	"cqrs/internal/core/course_scheduling/app/query/teacher"
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
