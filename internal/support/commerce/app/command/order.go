package command

import (
	"cqrs/internal/support/commerce/domain/repo/command"
	"cqrs/internal/support/commerce/domain/repo/query"
)

// OrderUsecase is the order usecase.
type OrderUsecase struct {
	Query   query.OrderQuery
	Command command.OrderCommand
}

// NewOrderUsecase creates a new OrderUsecase.
func NewOrderUsecase(q query.OrderQuery, c command.OrderCommand) *OrderUsecase {
	return &OrderUsecase{Query: q, Command: c}
}
