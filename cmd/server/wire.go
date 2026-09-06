//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/support/commerce"
	"github.com/go-kratos/kratos-layout/internal/core/product"
	"github.com/go-kratos/kratos-layout/internal/core/teaching"
	"github.com/go-kratos/kratos-layout/internal/core/venue"
	"github.com/go-kratos/kratos-layout/internal/server"

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
