package redis

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"

	goredis "github.com/redis/go-redis/v9"

	"payment/internal/conf"
)

// RedisData 持有 Redis 客户端，延时队列的方法都挂在它上面。
type RedisData struct {
	rdb   *goredis.Client
	podID string
	debug bool
}

// NewConfig 从 conf.Data 里取 Redis 段 —— wire 只能串联构造函数，取不了 struct 字段。
func NewConfig(c *conf.Data) *conf.Data_Redis {
	if c == nil {
		return nil
	}
	return c.GetRedis()
}

// NewRedisData 用 kratos 的 data.redis 配置建客户端。
//
// kratos 的 network 只管 tcp / unix，跟 go-redis 的 node / cluster 不是一回事：
// 单节点固定 tcp，要连集群得另开配置项并换成 *goredis.ClusterClient。
//
// 只做「配置 → 客户端」，不在这 ping：go-redis 自带连接池、按需建连，
// 启动时 ping 不过就起不来，反而把 Redis 抖动放大成服务不可用。
func NewRedisData(c *conf.Data_Redis) (*RedisData, error) {
	if c == nil {
		return nil, errors.New("redis: data.redis is required")
	}

	opts := &goredis.Options{
		Addr: c.GetAddr(),
	}
	// 超时都要给：连不上或对端半死不活时，别把消费 / 写入卡太久。
	if d := c.GetReadTimeout(); d != nil && d.AsDuration() > 0 {
		opts.ReadTimeout = d.AsDuration()
	}
	if d := c.GetWriteTimeout(); d != nil && d.AsDuration() > 0 {
		opts.WriteTimeout = d.AsDuration()
	}

	return NewredisData(goredis.NewClient(opts), podID()), nil
}

// NewredisData 用现成的客户端组装 —— 给集成测试注入自己连的实例用。
func NewredisData(rdb *goredis.Client, podID string) *RedisData {
	return &RedisData{rdb: rdb, podID: podID}
}

// Close 关闭连接池；进程退出前调用。
func (s *RedisData) Close() error {
	if s == nil || s.rdb == nil {
		return nil
	}
	return s.rdb.Close()
}

// podID 消费者名里区分实例的部分：优先 HOSTNAME（k8s / Docker 里就是 Pod 名），
// 取不到就随机一个 —— 两个实例撞名会让 Stream 的 pending 互相抢，
// 所以宁可不稳定也别重名。
func podID() string {
	if h := os.Getenv("HOSTNAME"); h != "" {
		return h
	}

	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b[:])
}
