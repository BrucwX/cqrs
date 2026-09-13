// Package query 是 course_scheduling 读侧的命令适配器。
//
// QueryImpl 是读侧查询实现的载体：持有 mysql（以后再加 redis）等存储客户端。
// 每个 domain/repo/query 接口的实现在各自的 <资源>.go 里 —— 当前先转发给
// Data，等 Redis 接入后再在对应方法里组合两者的读取。
package query

import (
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql"
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
