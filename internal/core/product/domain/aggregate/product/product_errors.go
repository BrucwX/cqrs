package product

import "github.com/go-kratos/kratos/v3/errors"

var (
	// ErrProductNotFound is returned when a product does not exist.
	ErrProductNotFound = errors.NotFound("PRODUCT_NOT_FOUND", "product not found")
	// ErrProductInvalidArgument is returned when a product request is invalid.
	ErrProductInvalidArgument = errors.BadRequest("PRODUCT_INVALID_ARGUMENT", "invalid product argument")
	// ErrProductOutOfStock is returned when product stock is insufficient.
	ErrProductOutOfStock = errors.BadRequest("PRODUCT_OUT_OF_STOCK", "product out of stock")
)
