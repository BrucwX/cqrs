package query

import (
	"github.com/google/wire"

	"payment/internal/adapters/query/data/mysql"
)

// ProviderSet 是读侧查询适配器的 provider 集合：读库客户端 + 载体 + 接口绑定。
//
// 加查询时的三步：
//
//  1. 接口放 domain/repo/query/<资源>.go
//  2. 实现在本包 <资源>.go，方法挂在 *QueryImpl 上
//  3. 在这里加一条 wire.Bind(new(repoquery.XxxQuery), new(*QueryImpl))
var ProviderSet = wire.NewSet(
	mysql.NewData,
	NewQueryImpl,
)
