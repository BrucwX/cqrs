package command

import (
	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/payment"
)

type paymentCommand struct {
	data *memory.Data
}

// NewPaymentCommand creates a new PaymentCommand instance.
func NewPaymentCommand(d *memory.Data) *paymentCommand {
	return &paymentCommand{data: d}
}

func (r *paymentCommand) Save(p *payment.Payment) error {
	// TODO: implement database insert/update
	return nil
}

func (r *paymentCommand) Delete(id string) error {
	// TODO: implement database soft delete
	return nil
}
