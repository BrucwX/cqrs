// Package redis 是 course_scheduling 写侧的 Redis 实现。
//
// 写侧的 Redis 用于缓存/发布等辅助存储，与 MySQL 写库配合使用。
//
//	redis.go  本文件：Redis 连接池（RedisData / NewRedisData）
//	model/    Redis PO 与 DO 转换（json 序列化）
//
// 与 MySQL 侧的差异：
//   - Redis model 不含 LockVersion（乐观锁由 MySQL 侧承担）
//   - Redis model 使用 json tag，MySQL PO 无 struct tag
//   - 转换函数命名为 XxxToRedis / XxxFromRedis
package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"cqrs/internal/conf"
)

// RedisData 持有写侧 Redis 连接池，写侧的所有仓库共享它。
type RedisData struct {
	client *redis.Client
	debug  bool
}

// NewRedisData 打开写侧 Redis 连接池。
func NewRedisData(c *conf.Data) (*RedisData, func(), error) {
	cfg := c.GetRedis()
	if cfg == nil {
		return nil, nil, errors.New("redis/command: data.redis is required")
	}

	network := cfg.GetNetwork()
	if network == "" {
		network = "tcp"
	}

	opts := &redis.Options{
		Network:      network,
		Addr:         cfg.GetAddr(),
		ReadTimeout:  cfg.GetReadTimeout().AsDuration(),
		WriteTimeout: cfg.GetWriteTimeout().AsDuration(),
	}

	client := redis.NewClient(opts)

	// 验证连通性
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, nil, err
	}

	cleanup := func() { client.Close() }
	return &RedisData{client: client}, cleanup, nil
}

// Client 返回底层 *redis.Client，供需要直接操作的场景使用。
func (r *RedisData) Client() *redis.Client {
	return r.client
}
