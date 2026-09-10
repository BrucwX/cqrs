package command

import (
	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/discount"
)

type discountCommand struct {
	data *memory.Data
}

// NewDiscountCommand creates a new DiscountCommand instance.
func NewDiscountCommand(d *memory.Data) *discountCommand {
	return &discountCommand{data: d}
}

func (r *discountCommand) Save(d *discount.Discount) error {
	// TODO: implement database insert/update
	return nil
}

func (r *discountCommand) Delete(id string) error {
	// TODO: implement database soft delete
	return nil
}
