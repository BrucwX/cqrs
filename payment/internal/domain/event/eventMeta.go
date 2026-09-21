package event

import (
	"time"

	"payment/pkg/uuid"
)

// EventMeta 提供领域事件的公共元数据实现
//
// 事件的「名字」不在这里 —— 那是各事件自己的 Topic（见 Event.Topic）。
type EventMeta struct {
	id          string
	occurredAt  time.Time
	aggregateID string
}

// NewEventMeta 初始化事件基础元数据：事件 id 与发生时间都在这生成。
func NewEventMeta(aggregateID string) EventMeta {
	return EventMeta{
		id:          uuid.NewString(),
		occurredAt:  time.Now().UTC(),
		aggregateID: aggregateID,
	}
}

func (b EventMeta) EventID() string       { return b.id }
func (b EventMeta) OccurredAt() time.Time { return b.occurredAt }
func (b EventMeta) AggregateID() string   { return b.aggregateID }
