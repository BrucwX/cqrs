package service

import (
	"context"

	"cqrs/internal/support/commerce/biz/discount"
	"cqrs/internal/shared/types"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// DiscountService is a discount service.
type DiscountService struct {
	uc *discount.DiscountUsecase
}

// NewDiscountService creates a new DiscountService.
func NewDiscountService(uc *discount.DiscountUsecase) *DiscountService {
	return &DiscountService{uc: uc}
}

// CreateDiscount creates a discount.
func (s *DiscountService) CreateDiscount(ctx context.Context, d *discount.Discount) (*discount.Discount, error) {
	return s.uc.Repo.CreateDiscount(ctx, d)
}

// GetDiscount returns a discount by ID.
func (s *DiscountService) GetDiscount(ctx context.Context, id string) (*discount.Discount, error) {
	return s.uc.Repo.FindByID(ctx, id)
}

// GetByCode returns a discount by code.
func (s *DiscountService) GetByCode(ctx context.Context, code string) (*discount.Discount, error) {
	return s.uc.Repo.FindByCode(ctx, code)
}

// ListDiscounts lists discounts.
func (s *DiscountService) ListDiscounts(ctx context.Context, opts ...types.ListOption) ([]*discount.Discount, error) {
	return s.uc.Repo.ListDiscounts(ctx, opts...)
}

// UpdateDiscount updates a discount with field mask.
func (s *DiscountService) UpdateDiscount(ctx context.Context, d *discount.Discount, mask *fieldmaskpb.FieldMask) (*discount.Discount, error) {
	return s.uc.Repo.UpdateDiscount(ctx, d, mask)
}

// DeleteDiscount deletes a discount.
func (s *DiscountService) DeleteDiscount(ctx context.Context, id string) error {
	return s.uc.Repo.DeleteDiscount(ctx, id)
}
