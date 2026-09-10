package command

import (
	"context"

	"cqrs/internal/core/product/adapters/memory"
	"cqrs/internal/core/product/domain/aggregate/product"
	"cqrs/internal/core/product/domain/repo"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type productRepo struct {
	data *memory.Data
}

// NewProductRepo creates a new ProductRepo instance.
func NewProductRepo(d *memory.Data) repo.ProductRepo {
	return &productRepo{data: d}
}

func (r *productRepo) FindByID(ctx context.Context, id string) (*product.Product, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *productRepo) CreateProduct(ctx context.Context, p *product.Product) (*product.Product, error) {
	// TODO: implement database insert
	return p, nil
}

func (r *productRepo) UpdateProduct(ctx context.Context, p *product.Product, mask *fieldmaskpb.FieldMask) (*product.Product, error) {
	// TODO: implement database update
	// 根据 mask 决定更新哪些字段
	return p, nil
}

func (r *productRepo) DeleteProduct(ctx context.Context, id string) error {
	// TODO: implement database soft delete
	return nil
}

func (r *productRepo) UpdateStock(ctx context.Context, id string, quantity int32) error {
	// TODO: implement stock update
	return nil
}
