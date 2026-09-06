package payment

// PaymentUsecase is the payment usecase.
type PaymentUsecase struct {
	Repo PaymentRepo
}

// NewPaymentUsecase creates a new PaymentUsecase.
func NewPaymentUsecase(repo PaymentRepo) *PaymentUsecase {
	return &PaymentUsecase{Repo: repo}
}
