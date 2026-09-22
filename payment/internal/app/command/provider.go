// Package command 是 payment 写侧用例层。
//
// 本包只放 provider.go：汇总各 handler 子包的构造函数，以及命令侧专用的
// 领域判定服务（domain/service）。
//
// 端口（Transaction / CourseService / PayService / IdempotencyService）在
// internal/ports/ 里 —— 用例只认接口，实现由 adapters/ 提供。
package command

import "github.com/google/wire"

// ProviderSet 命令侧集合。
//
// 加用例时：<聚合>.NewHandler 写进这里（子包名用 lowerCamelCase，如 courseSlot）。
var ProviderSet = wire.NewSet()
