package repo

import (
	"context"

	"cqrs/internal/core/product/domain/aggregate/product"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// ProductRepo is the product repository interface.
type ProductRepo interface {
	FindByID(ctx context.Context, id string) (*product.Product, error)
	CreateProduct(ctx context.Context, p *product.Product) (*product.Product, error)
	UpdateProduct(ctx context.Context, p *product.Product, mask *fieldmaskpb.FieldMask) (*product.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, quantity int32) error
}
