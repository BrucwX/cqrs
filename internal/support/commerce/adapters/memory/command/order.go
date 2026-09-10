package command

import (
	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/order"
)

type orderCommand struct {
	data *memory.Data
}

// NewOrderCommand creates a new OrderCommand instance.
func NewOrderCommand(d *memory.Data) *orderCommand {
	return &orderCommand{data: d}
}

func (r *orderCommand) Save(o *order.Order) error {
	// TODO: implement database insert/update
	return nil
}

func (r *orderCommand) Delete(id string) error {
	// TODO: implement database soft delete
	return nil
}
