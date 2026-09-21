package event

const topic = "Order"

type OrderPlacedEvent struct {
	EventMeta
	CustomerId string
}

func (O *OrderPlacedEvent) Topic() EventTopic { return topic }

func (O *OrderPlacedEvent) isDomainEvent() {}
