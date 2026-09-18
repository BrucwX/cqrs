// Package service 是 fulfillment 的传输层：一个聚合根一个文件，负责
// proto 请求/响应 <-> 用例入参的转换、以及领域错误到 Kratos 错误的映射。
//
// 不在这里写业务规则 —— 规则在 app/（用例）与 domain/（聚合、领域服务）。
package service

import "github.com/google/wire"

// ProviderSet 传输层集合。
//
// 加资源时：NewXxxService 写进这里。
var ProviderSet = wire.NewSet()
