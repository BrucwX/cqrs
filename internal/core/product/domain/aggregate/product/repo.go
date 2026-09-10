package product

import (
	"context"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// ProductRepo is the product repository interface.
type ProductRepo interface {
	FindByID(ctx context.Context, id string) (*Product, error)
	CreateProduct(ctx context.Context, p *Product) (*Product, error)
	UpdateProduct(ctx context.Context, p *Product, mask *fieldmaskpb.FieldMask) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, quantity int32) error
}
