package order

// OrderUsecase is the order usecase.
type OrderUsecase struct {
	Repo OrderRepo
}

// NewOrderUsecase creates a new OrderUsecase.
func NewOrderUsecase(repo OrderRepo) *OrderUsecase {
	return &OrderUsecase{Repo: repo}
}
