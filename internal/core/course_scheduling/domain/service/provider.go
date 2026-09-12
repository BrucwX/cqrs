package service

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/domain/service/classroomCapacity"
	"cqrs/internal/core/course_scheduling/domain/service/qualificationCheck"
	"cqrs/internal/core/course_scheduling/domain/service/scheduleConflict"
)

// ProviderSet 领域判定服务的构造器。
//
// 一个服务一个子包，包名 = 服务名；每个包里定结构体的文件都叫 service.go，
// 结构体都叫 Service：
//
//	scheduleConflict.Service       要占的这段时间是不是已经有人占了
//	classroomCapacity.Service      教室装不装得下
//	qualificationCheck.Service     讲师资质够不够
//
// 这 3 个服务都跨聚合根、被多个命令处理器共享，所以集中注册在一个 Set 里；
// app/command 的 handler 直接用 wire 注入 *scheduleConflict.Service 这类类型即可。
var ProviderSet = wire.NewSet(
	scheduleConflict.NewService,
	classroomCapacity.NewService,
	qualificationCheck.NewService,
)
