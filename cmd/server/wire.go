//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"cqrs/internal/conf"
	"cqrs/internal/core/course_scheduling"
	"cqrs/internal/server"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
//
// 目前只 wire course_scheduling 一个上下文。product / commerce 的 ProviderSet
// 先不挂进来：它们自己的 provider 还没补齐，挂进来会让整个应用没法生成。
func wireApp(*conf.Server, *conf.Data, *slog.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		course_scheduling.ProviderSet,
		server.ProviderSet,
		newApp,
	))
}
