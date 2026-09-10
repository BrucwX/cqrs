package repo

import (
	"context"

	"cqrs/internal/support/commerce/domain/aggregate/discount"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// DiscountRepo is the discount repository interface.
type DiscountRepo interface {
	FindByID(ctx context.Context, id string) (*discount.Discount, error)
	FindByCode(ctx context.Context, code string) (*discount.Discount, error)
	CreateDiscount(ctx context.Context, d *discount.Discount) (*discount.Discount, error)
	UpdateDiscount(ctx context.Context, d *discount.Discount, mask *fieldmaskpb.FieldMask) (*discount.Discount, error)
	DeleteDiscount(ctx context.Context, id string) error
	IncrementUsage(ctx context.Context, id string) error
}
