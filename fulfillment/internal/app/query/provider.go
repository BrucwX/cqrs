// Package query 是 fulfillment 读侧用例层。
//
// 一个聚合根一个子包，每个子包一个 handler（Handler 持有它用到的
// domain/repo/query 接口，不塞用不到的）。
package query

import "github.com/google/wire"

// ProviderSet 读侧集合。
//
// 加查询时：<聚合>.NewHandler 写进这里。
var ProviderSet = wire.NewSet()
