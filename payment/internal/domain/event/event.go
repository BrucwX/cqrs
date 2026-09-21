package event

import (
	"context"
	"time"
)

type EventTopic string

// Event 是所有领域事件必须实现的核心接口
type Event interface {
	// EventID 事件的全局唯一标识（防重、幂等排重、追踪）
	EventID() string

	// OccurredAt 事件发生的真实业务时间点（UTC）
	OccurredAt() time.Time

	// AggregateID 产生该事件的聚合根 ID（用于分区路由、因果时序追踪）
	AggregateID() string

	// Topic 事件的**唯一名字**：订阅路由、跨服务消费与序列化契约都靠它。
	Topic() EventTopic

	isDomainEvent()
}

// 2. EventHandler 类型：订阅者回调函数
type EventHandler func(ctx context.Context, event Event) error

// 3. EventBus 接口：暴露总线操作
type EventBus interface {
	Subscribe(topic EventTopic, handler EventHandler)
	Publish(ctx context.Context, event Event) error // 同步广播（阻塞等待全部执行完）
	PublishAsync(ctx context.Context, event Event)  // 异步广播（非阻塞，goroutine 中分发）
}
