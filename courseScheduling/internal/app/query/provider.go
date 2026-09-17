package query

import (
	"github.com/google/wire"

	"courseScheduling/internal/app/query/absence"
	"courseScheduling/internal/app/query/classroom"
	"courseScheduling/internal/app/query/course"
	"courseScheduling/internal/app/query/courseSlot"
	"courseScheduling/internal/app/query/courseSlotChange"
	"courseScheduling/internal/app/query/enrollment"
	"courseScheduling/internal/app/query/makeup"
	"courseScheduling/internal/app/query/qualification"
	"courseScheduling/internal/app/query/student"
	"courseScheduling/internal/app/query/teacher"
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
