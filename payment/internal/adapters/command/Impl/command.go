// Package command 是 payment 写侧的存储适配器载体。
//
// CommandImpl 是写侧命令实现的载体：持有 mysql（以后再加 redis）等存储客户端。
// 每个 聚合包里的 <聚合根>_repo.go 接口的实现在各自的 <资源>.go 里 —— 当前先转发给
// MysqlData，等 Redis 接入后再在对应方法里组合两者的读写。
package command

import (
	"payment/internal/adapters/command/data/mysql"
	"payment/internal/adapters/command/data/redis"
)

// CommandImpl 写侧命令实现。
type CommandImpl struct {
	MysqlData *mysql.MysqlData
	// RedisData 目前只承载延时队列（ports.DelayScheduler），
	// 各聚合仓储要读写 Redis 时再在对应方法里组合它。
	RedisData *redis.RedisData
}

// NewCommandImpl 组装写侧命令实现。
func NewCommandImpl(mysqlData *mysql.MysqlData, redisData *redis.RedisData) *CommandImpl {
	return &CommandImpl{MysqlData: mysqlData, RedisData: redisData}
}
