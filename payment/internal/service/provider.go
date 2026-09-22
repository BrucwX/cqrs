// Package service 是 payment 的传输层：HTTP / gRPC 服务器本身，加上
// proto 请求/响应 <-> 用例入参的转换、领域错误到 Kratos 错误的映射。
//
// 不在这里写业务规则 —— 规则在 app/（用例）与 domain/（聚合、领域服务）。
package service

import "github.com/google/wire"

// ProviderSet 传输层集合：服务器 + 各 service。
//
// 加资源时：NewXxxService 写进这里，并在下面 http.go / grpc.go 里挂到服务器上。
var ProviderSet = wire.NewSet(
	NewServers,
	NewHTTPServer,
	NewGRPCServer,
)
