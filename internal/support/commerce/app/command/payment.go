package command

import "cqrs/internal/support/commerce/domain/repo"

// PaymentUsecase is the payment usecase.
type PaymentUsecase struct {
	Repo repo.PaymentRepo
}

// NewPaymentUsecase creates a new PaymentUsecase.
func NewPaymentUsecase(r repo.PaymentRepo) *PaymentUsecase {
	return &PaymentUsecase{Repo: r}
}
