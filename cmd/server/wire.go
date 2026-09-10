//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"cqrs/internal/conf"
	"cqrs/internal/core/product"
	"cqrs/internal/core/course_scheduling"
	"cqrs/internal/server"
	"cqrs/internal/support/commerce"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *slog.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		course_scheduling.ProviderSet,
		product.ProviderSet,
		commerce.ProviderSet,
		server.ProviderSet,
		newApp,
	))
}
