package redis

import (
	"github.com/redis/go-redis/v9"

	"payment/internal/ports"
)

// 本文件不实现队列逻辑（逻辑在 impl_queue.go / help_*.go），只做「接线」：
// 类型别名（让那几份代码里的写法对上本包的类型）、Notify 用的 Lua 脚本。

type (
	// redisData 是 RedisData 的别名：help_*.go 里用小写写法，跟 RedisData 是同一个类型。
	redisData = RedisData

	// CourseHandler 是 ports.CourseHandler 的别名：签名完全一致，
	// 所以 StartCloseCourse 直接满足 ports.DelayScheduler，不用改签名。
	CourseHandler = ports.CourseHandler
)

// pushAndNotifyLua 一条命令里原子完成「写 ZSET + 发唤醒信号」：
//
//	KEYS[1] = zset        KEYS[2] = stream
//	ARGV[1] = 到期 Unix 秒  ARGV[2] = orderID
//
// 分两条命令发的话，中间崩溃会留下「ZSET 里有计划、但没人被叫醒」的窗口；
// 消费端靠启动时的全量扫描能自愈，但既然能原子就别留这个窗口。
// XADD 带 MAXLEN ~：唤醒信号是一次性的，留一千条给还没来得及重算睡眠时间的消费者就够。
var pushAndNotifyLua = redis.NewScript(`
redis.call('ZADD', KEYS[1], ARGV[1], ARGV[2])
redis.call('XADD', KEYS[2], 'MAXLEN', '~', 1000, '*', 'action', 'wakeup')
return 1
`)
