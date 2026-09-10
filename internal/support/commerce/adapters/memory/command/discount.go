package command

import (
	"context"

	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/discount"
	"cqrs/internal/support/commerce/domain/repo"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type discountRepo struct {
	data *memory.Data
}

// NewDiscountRepo creates a new DiscountRepo instance.
func NewDiscountRepo(d *memory.Data) repo.DiscountRepo {
	return &discountRepo{data: d}
}

func (r *discountRepo) FindByID(ctx context.Context, id string) (*discount.Discount, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *discountRepo) FindByCode(ctx context.Context, code string) (*discount.Discount, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *discountRepo) CreateDiscount(ctx context.Context, d *discount.Discount) (*discount.Discount, error) {
	// TODO: implement database insert
	return d, nil
}

func (r *discountRepo) UpdateDiscount(ctx context.Context, d *discount.Discount, mask *fieldmaskpb.FieldMask) (*discount.Discount, error) {
	// TODO: implement database update
	// 根据 mask 决定更新哪些字段
	return d, nil
}

func (r *discountRepo) DeleteDiscount(ctx context.Context, id string) error {
	// TODO: implement database soft delete
	return nil
}

func (r *discountRepo) IncrementUsage(ctx context.Context, id string) error {
	// TODO: implement usage increment
	return nil
}
