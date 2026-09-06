package discount

import (
	"context"

	"cqrs/internal/shared/types"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// DiscountRepo is the discount repository interface.
type DiscountRepo interface {
	FindByID(ctx context.Context, id string) (*Discount, error)
	FindByCode(ctx context.Context, code string) (*Discount, error)
	ListDiscounts(ctx context.Context, opts ...types.ListOption) ([]*Discount, error)
	CreateDiscount(ctx context.Context, d *Discount) (*Discount, error)
	UpdateDiscount(ctx context.Context, d *Discount, mask *fieldmaskpb.FieldMask) (*Discount, error)
	DeleteDiscount(ctx context.Context, id string) error
	IncrementUsage(ctx context.Context, id string) error
}
