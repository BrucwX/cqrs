package event

import "time"

// Event 是所有领域事件必须实现的核心接口
type DomainEvent interface {
	// EventID 事件的全局唯一标识（防重、幂等排重、追踪）
	EventID() string

	// EventName 显式的事件名称（如 "order.created"、"order.paid"）
	// 避免使用反射获取类型名，方便跨服务消费和序列化契约
	EventName() string

	// OccurredAt 事件发生的真实业务时间点（UTC）
	OccurredAt() time.Time

	// AggregateID 产生该事件的聚合根 ID（用于分区路由、因果时序追踪）
	AggregateID() string
}
