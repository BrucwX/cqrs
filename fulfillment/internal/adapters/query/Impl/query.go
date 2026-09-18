// Package query 是 fulfillment 读侧的存储适配器载体。
//
// QueryImpl 持有读库客户端；每个 domain/repo/query 接口的实现在各自的
// <资源>.go 里，当前先转发给 mysql.Data。
package query

import (
	"fulfillment/internal/adapters/query/data/mysql"
)

// QueryImpl 读侧查询实现。
type QueryImpl struct {
	Data *mysql.Data
	// RedisData *redis.RedisData // 待接入
}

// NewQueryImpl 组装读侧查询实现。
func NewQueryImpl(data *mysql.Data) *QueryImpl {
	return &QueryImpl{Data: data}
}
