// Package command 是命令侧用例层。
//
// 本包只放一个文件 provider.go：汇总 10 个 handler 子包的构造函数，
// 以及命令侧专用的领域判定服务（domain/service）。
// Transaction 端口在子包 transaction/ 里 —— 子包不 import 本包，本包 import 子包，
// 所以不会成环。
package command

import (
	"github.com/google/wire"

	"courseScheduling/internal/app/command/absence"
	"courseScheduling/internal/app/command/classroom"
	"courseScheduling/internal/app/command/course"
	"courseScheduling/internal/app/command/courseSlot"
	"courseScheduling/internal/app/command/courseSlotChange"
	"courseScheduling/internal/app/command/enrollment"
	"courseScheduling/internal/app/command/makeup"
	"courseScheduling/internal/app/command/qualification"
	"courseScheduling/internal/app/command/student"
	"courseScheduling/internal/app/command/teacher"
	domainservice "courseScheduling/internal/domain/service"
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
