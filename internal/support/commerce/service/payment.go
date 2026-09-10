package service

import "cqrs/internal/support/commerce/app/command"

// PaymentService is a payment service.
type PaymentService struct {
	uc *command.PaymentUsecase
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(uc *command.PaymentUsecase) *PaymentService {
	return &PaymentService{uc: uc}
}
