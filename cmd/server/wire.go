//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"cqrs/internal/conf"
	"cqrs/internal/support/commerce"
	"cqrs/internal/core/product"
	"cqrs/internal/core/teaching"
	"cqrs/internal/core/venue"
	"cqrs/internal/server"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *slog.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		teaching.ProviderSet,
		venue.ProviderSet,
		product.ProviderSet,
		commerce.ProviderSet,
		server.ProviderSet,
		newApp,
	))
}
