package query

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/domain/service/query/absence"
	"cqrs/internal/core/course_scheduling/domain/service/query/classroom"
	"cqrs/internal/core/course_scheduling/domain/service/query/course"
	"cqrs/internal/core/course_scheduling/domain/service/query/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/service/query/courseSlotChange"
	"cqrs/internal/core/course_scheduling/domain/service/query/enrollment"
	"cqrs/internal/core/course_scheduling/domain/service/query/makeup"
	"cqrs/internal/core/course_scheduling/domain/service/query/qualification"
	"cqrs/internal/core/course_scheduling/domain/service/query/student"
	"cqrs/internal/core/course_scheduling/domain/service/query/teacher"
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
