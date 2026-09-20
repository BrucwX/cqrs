package event

import (
	"time"

	"payment/pkg/uuid"
)

// EventMeta 提供领域事件的公共元数据实现
type EventMeta struct {
	id          string
	name        string
	occurredAt  time.Time
	aggregateID string
}

// NewBaseEvent 初始化事件基础元数据
func NewEventMeta(name, aggregateID string) EventMeta {
	return EventMeta{
		id:          uuid.NewString(),
		name:        name,
		occurredAt:  time.Now().UTC(),
		aggregateID: aggregateID,
	}
}

func (b EventMeta) EventID() string       { return b.id }
func (b EventMeta) EventName() string     { return b.name }
func (b EventMeta) OccurredAt() time.Time { return b.occurredAt }
func (b EventMeta) AggregateID() string   { return b.aggregateID }
