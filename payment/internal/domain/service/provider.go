// Package service 放跨聚合根的领域规则（写路径的准入判定）。
//
// 组织原则：一个服务对「一类问题」，不对「一个用例」；同一类问题的新场景
// 加一个入口方法，不新开服务。规则是纯的：入参是聚合根 + 基本类型，
// 不碰存储（取数由调用方 app handler 通过 repo 接口做好后传进来）。
//
// 包名 = 服务名（lowerCamelCase），每个包里定结构体的文件都叫 service.go，
// 结构体叫 Service，构造函数叫 NewService。
package service

import (
	"github.com/google/wire"

	ordercheck "payment/internal/domain/service/orderCheck"
	orderitemscheck "payment/internal/domain/service/orderItemsCheck"
)

// ProviderSet 领域服务集合。
//
// 加服务时：NewXxxService 写进这里。
var ProviderSet = wire.NewSet(
	ordercheck.NewService,
	orderitemscheck.NewService,
)
