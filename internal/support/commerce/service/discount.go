package service

import (
	"context"

	"cqrs/internal/support/commerce/app/command"
	"cqrs/internal/support/commerce/domain/aggregate/discount"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// DiscountService is a discount service.
type DiscountService struct {
	uc *command.DiscountUsecase
}

// NewDiscountService creates a new DiscountService.
func NewDiscountService(uc *command.DiscountUsecase) *DiscountService {
	return &DiscountService{uc: uc}
}

// CreateDiscount creates a discount.
func (s *DiscountService) CreateDiscount(ctx context.Context, d *discount.Discount) (*discount.Discount, error) {
	if err := s.uc.Command.Save(d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetDiscount returns a discount by ID.
func (s *DiscountService) GetDiscount(ctx context.Context, id string) (*discount.Discount, error) {
	return s.uc.Query.GetByID(id)
}

// GetByCode returns a discount by code.
func (s *DiscountService) GetByCode(ctx context.Context, code string) (*discount.Discount, error) {
	return s.uc.Query.GetByCode(code)
}

// UpdateDiscount updates a discount with field mask.
func (s *DiscountService) UpdateDiscount(ctx context.Context, d *discount.Discount, mask *fieldmaskpb.FieldMask) (*discount.Discount, error) {
	if err := s.uc.Command.Save(d); err != nil {
		return nil, err
	}
	return d, nil
}

// DeleteDiscount deletes a discount.
func (s *DiscountService) DeleteDiscount(ctx context.Context, id string) error {
	return s.uc.Command.Delete(id)
}
