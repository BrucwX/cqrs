package product

// ProductUsecase is the product usecase.
type ProductUsecase struct {
	Repo ProductRepo
}

// NewProductUsecase creates a new ProductUsecase.
func NewProductUsecase(repo ProductRepo) *ProductUsecase {
	return &ProductUsecase{Repo: repo}
}
