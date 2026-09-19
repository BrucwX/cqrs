package command

import (
	"github.com/google/wire"

	"payment/internal/adapters/command/data/mysql"
	appcmd "payment/internal/app/command/appRepo"
	ordercheck "payment/internal/domain/service/orderCheck"
	orderitemscheck "payment/internal/domain/service/orderItemsCheck"
)

// ProviderSet 是写侧命令适配器的 provider 集合：存储客户端 + 载体 + 接口绑定。
//
// 加仓库 / 端口的步骤：
//
//  1. 接口放聚合旁边（domain/aggregate/<聚合>/<聚合>_repo.go）或 app 层（app/command/appRepo）
//  2. 实现在本包 <资源>.go，方法挂在 *CommandImpl 上
//  3. 在这里加 wire.Bind(new(接口), new(*CommandImpl))
//
// ⚠️ 两个领域服务各自声明了自己的 Repo（只用得到的那几个方法），所以也要各绑一条 ——
// 载体满足签名就够了，服务不认识「完整的仓储」。
var ProviderSet = wire.NewSet(
	mysql.NewMysqlData,
	NewCommandImpl,
	wire.Bind(new(appcmd.Transaction), new(*CommandImpl)),
	wire.Bind(new(ordercheck.Repo), new(*CommandImpl)),
	wire.Bind(new(orderitemscheck.Repo), new(*CommandImpl)),
)
