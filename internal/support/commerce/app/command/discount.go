package command

import (
	"cqrs/internal/support/commerce/domain/repo/command"
	"cqrs/internal/support/commerce/domain/repo/query"
)

// DiscountUsecase is the discount usecase.
type DiscountUsecase struct {
	Query   query.DiscountQuery
	Command command.DiscountCommand
}

// NewDiscountUsecase creates a new DiscountUsecase.
func NewDiscountUsecase(q query.DiscountQuery, c command.DiscountCommand) *DiscountUsecase {
	return &DiscountUsecase{Query: q, Command: c}
}
