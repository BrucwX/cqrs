package command

import (
	"github.com/google/wire"

	"payment/internal/adapters/command/data/mysql"
	"payment/internal/adapters/command/data/redis"
	"payment/internal/domain/repo"
	"payment/internal/ports"
)

var ProviderSet = wire.NewSet(
	redis.NewConfig,
	redis.NewRedisData,
	mysql.NewMysqlData,
	NewCommandImpl,
	wire.Bind(new(ports.Transaction), new(*CommandImpl)),
	wire.Bind(new(repo.OrderCheckRepo), new(*CommandImpl)),
	wire.Bind(new(repo.OrderItemsCheckRepo), new(*CommandImpl)),
	wire.Bind(new(ports.DelayScheduler), new(*CommandImpl)),
)
