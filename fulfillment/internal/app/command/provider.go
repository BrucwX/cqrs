// Package command 是 fulfillment 写侧用例层。
//
// 本包只放 provider.go：汇总各 handler 子包的构造函数，以及命令侧专用的
// 领域判定服务（domain/service）。
//
// Transaction 端口在子包 appRepo/ 里 —— handler 子包 import 端口，而本包
// import handler 子包，两端都不碰对方，所以不会成环。
package command

import "github.com/google/wire"

// ProviderSet 命令侧集合。
//
// 加用例时：<聚合>.NewHandler 写进这里（子包名用 lowerCamelCase，如 courseSlot）。
var ProviderSet = wire.NewSet()
