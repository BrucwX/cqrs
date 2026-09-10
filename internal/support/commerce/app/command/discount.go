package command

import "cqrs/internal/support/commerce/domain/repo"

// DiscountUsecase is the discount usecase.
type DiscountUsecase struct {
	Repo repo.DiscountRepo
}

// NewDiscountUsecase creates a new DiscountUsecase.
func NewDiscountUsecase(r repo.DiscountRepo) *DiscountUsecase {
	return &DiscountUsecase{Repo: r}
}
