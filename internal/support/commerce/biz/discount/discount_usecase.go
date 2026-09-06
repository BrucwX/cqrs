package discount

// DiscountUsecase is the discount usecase.
type DiscountUsecase struct {
	Repo DiscountRepo
}

// NewDiscountUsecase creates a new DiscountUsecase.
func NewDiscountUsecase(repo DiscountRepo) *DiscountUsecase {
	return &DiscountUsecase{Repo: repo}
}
