//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	fulfillment "fulfillment/internal"
	"fulfillment/internal/conf"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
//
// fulfillment 上下文只有一个，它的 ProviderSet（internal 根包）把
// adapters / app / domain / service / ports 一次装齐；HTTP/gRPC 由 ports 提供。
func wireApp(*conf.Server, *conf.Data, *slog.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		fulfillment.ProviderSet,
		newApp,
	))
}
