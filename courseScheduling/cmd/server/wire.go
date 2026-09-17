//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	course_scheduling "cqrs/internal"
	"cqrs/internal/conf"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
//
// 全项目只有一个界限上下文 course_scheduling，它的 ProviderSet（internal 根包）
// 把 adapters / app / service / ports 一次装齐；HTTP/gRPC 服务器由 ports 提供。
func wireApp(*conf.Server, *conf.Data, *slog.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		course_scheduling.ProviderSet,
		newApp,
	))
}
