package service

import (
	"cqrs/internal/support/commerce/biz/payment"
)

// PaymentService is a payment service.
type PaymentService struct {
	uc *payment.PaymentUsecase
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(uc *payment.PaymentUsecase) *PaymentService {
	return &PaymentService{uc: uc}
}
