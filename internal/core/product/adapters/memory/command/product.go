package command

import (
	"cqrs/internal/core/product/adapters/memory"
	"cqrs/internal/core/product/domain/aggregate/product"
)

type productCommand struct {
	data *memory.Data
}

// NewProductCommand creates a new ProductCommand instance.
func NewProductCommand(d *memory.Data) *productCommand {
	return &productCommand{data: d}
}

func (r *productCommand) Save(p *product.Product) error {
	// TODO: implement database insert/update
	return nil
}

func (r *productCommand) Delete(id string) error {
	// TODO: implement database soft delete
	return nil
}
