// Package redis 是 payment 的 Redis 侧适配。
//
// 目前只有定时器：
//
//	redis.go      连库：把 kratos 的 conf.Data_Redis 翻成 go-zero 的客户端
//	scheduler.go  Redis 版 Scheduler / Timer（跨实例的定时器）
//
// 落在这里（而不是 adapters/command/Impl）是因为它跟 mysql 一样是「基础设施实现」：
// Impl 里挂的是各聚合的仓储，这里挂的是别的技术能力。
package redis

import (
	"errors"

	gzredis "github.com/zeromicro/go-zero/core/stores/redis"

	"payment/internal/conf"
)

// NewRedisClient 用 kratos 的 data.redis 配置建一个 go-zero 客户端。
//
// 只做「配置 → 客户端」，不缓存、不复用：谁需要谁持有。go-zero 的 Redis 自带连接池
// 与熔断，没有（也不需要）Close 一个池子的动作。
func NewRedisClient(c *conf.Data_Redis) (*gzredis.Redis, error) {
	if c == nil {
		return nil, errors.New("redis: data.redis is required")
	}

	rc := gzredis.RedisConf{
		Host: c.GetAddr(),
		// kratos 的 network 只管 tcp / unix，跟这里的 node / cluster 不是一回事：
		// 单节点就写死 node，要连集群得另开一个配置项。
		Type: gzredis.NodeType,
	}
	// 读超时借来当 ping 超时：连不上的时候别把启动卡太久。
	if d := c.GetReadTimeout(); d != nil && d.AsDuration() > 0 {
		rc.PingTimeout = d.AsDuration()
	}

	return gzredis.NewRedis(rc)
}
