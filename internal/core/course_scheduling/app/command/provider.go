// Package command 是命令侧用例层。
//
// 本包只放一个文件 provider.go：汇总 10 个 handler 子包的构造函数，
// 以及命令侧专用的领域判定服务（domain/service）。
// Transaction 端口在子包 transaction/ 里 —— 子包不 import 本包，本包 import 子包，
// 所以不会成环。
package command

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/app/command/absence"
	"cqrs/internal/core/course_scheduling/app/command/classroom"
	"cqrs/internal/core/course_scheduling/app/command/course"
	"cqrs/internal/core/course_scheduling/app/command/courseSlot"
	"cqrs/internal/core/course_scheduling/app/command/courseSlotChange"
	"cqrs/internal/core/course_scheduling/app/command/enrollment"
	"cqrs/internal/core/course_scheduling/app/command/makeup"
	"cqrs/internal/core/course_scheduling/app/command/qualification"
	"cqrs/internal/core/course_scheduling/app/command/student"
	"cqrs/internal/core/course_scheduling/app/command/teacher"
	domainservice "cqrs/internal/core/course_scheduling/domain/service"
)

// ProviderSet 命令侧集合：10 个 handler + 领域判定服务。
var ProviderSet = wire.NewSet(
	absence.NewHandler,
	classroom.NewHandler,
	course.NewHandler,
	courseSlot.NewHandler,
	courseSlotChange.NewHandler,
	makeup.NewHandler,
	enrollment.NewHandler,
	qualification.NewHandler,
	student.NewHandler,
	teacher.NewHandler,
	// domain - 判定服务（跨聚合根的领域规则，只在写路径里用）
	domainservice.ProviderSet,
)
