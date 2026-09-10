package command

import (
	"cqrs/internal/support/commerce/domain/repo/command"
	"cqrs/internal/support/commerce/domain/repo/query"
)

// PaymentUsecase is the payment usecase.
type PaymentUsecase struct {
	Query   query.PaymentQuery
	Command command.PaymentCommand
}

// NewPaymentUsecase creates a new PaymentUsecase.
func NewPaymentUsecase(q query.PaymentQuery, c command.PaymentCommand) *PaymentUsecase {
	return &PaymentUsecase{Query: q, Command: c}
}
