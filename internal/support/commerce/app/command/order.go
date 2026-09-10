package command

import "cqrs/internal/support/commerce/domain/repo"

// OrderUsecase is the order usecase.
type OrderUsecase struct {
	Repo repo.OrderRepo
}

// NewOrderUsecase creates a new OrderUsecase.
func NewOrderUsecase(r repo.OrderRepo) *OrderUsecase {
	return &OrderUsecase{Repo: r}
}
