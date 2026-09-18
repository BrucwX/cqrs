package command

import (
	"github.com/google/wire"

	"fulfillment/internal/adapters/command/data/mysql"
	appcmd "fulfillment/internal/app/command/appRepo"
)

// ProviderSet 是写侧命令适配器的 provider 集合：存储客户端 + 载体 + 接口绑定。
//
// 加仓库时的三步：
//
//  1. 接口放 domain/repo/command/<资源>.go
//  2. 实现在本包 <资源>.go，方法挂在 *CommandImpl 上
//  3. 在这里加一条 wire.Bind(new(repo.XxxCommand), new(*CommandImpl))
var ProviderSet = wire.NewSet(
	mysql.NewMysqlData,
	NewCommandImpl,
	wire.Bind(new(appcmd.Transaction), new(*CommandImpl)),
)
