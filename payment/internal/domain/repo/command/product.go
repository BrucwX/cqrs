package command

import (
	"context"
	"payment/internal/domain/aggregate/product"
)

type ProductCommand interface {
	MustGetProduct(ctx context.Context, productID int64) (product.Product, error)
}
