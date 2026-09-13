// Package command 是 course_scheduling 写侧的命令适配器。
//
// CommandImpl 是写侧命令实现的载体：持有 mysql（以后再加 redis）等存储客户端。
// 每个 domain/repo/command 接口的实现在各自的 <资源>.go 里 —— 当前先转发给
// MysqlData，等 Redis 接入后再在对应方法里组合两者的读写。
package command

import (
	"cqrs/internal/core/course_scheduling/adapters/command/data/mysql"
)

// CommandImpl 写侧命令实现。
type CommandImpl struct {
	MysqlData *mysql.MysqlData
	// RedisData *redis.RedisData // 待接入
}

// NewCommandImpl 组装写侧命令实现。
func NewCommandImpl(mysqlData *mysql.MysqlData) *CommandImpl {
	return &CommandImpl{MysqlData: mysqlData}
}
