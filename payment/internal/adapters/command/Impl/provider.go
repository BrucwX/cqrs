package command

import (
	"github.com/google/wire"

	"payment/internal/adapters/command/data/mysql"
	"payment/internal/domain/repo"
	"payment/internal/ports"
)

// ProviderSet 是写侧命令适配器的 provider 集合：存储客户端 + 载体 + 接口绑定。
//
// 加仓库 / 端口的步骤：
//
//  1. 接口放 domain/repo/（聚合的写侧仓储）或 internal/ports/（事务 / 外部能力端口）
//  2. 实现在本包 <资源>.go，方法挂在 *CommandImpl 上
//  3. 在这里加 wire.Bind(new(接口), new(*CommandImpl))
//
// ⚠️ 两个领域服务用的是各自「用得到的那几个方法」的接口（repo.OrderCheckRepo /
// repo.OrderItemsCheckRepo），所以也要各绑一条 —— 载体满足签名就够了，服务不认识「完整的仓储」。
var ProviderSet = wire.NewSet(
	mysql.NewMysqlData,
	NewCommandImpl,
	wire.Bind(new(ports.Transaction), new(*CommandImpl)),
	wire.Bind(new(repo.OrderCheckRepo), new(*CommandImpl)),
	wire.Bind(new(repo.OrderItemsCheckRepo), new(*CommandImpl)),
)
