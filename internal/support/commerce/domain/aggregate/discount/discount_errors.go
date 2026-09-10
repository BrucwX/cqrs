package discount

import "github.com/go-kratos/kratos/v3/errors"

var (
	// ErrDiscountNotFound is returned when a discount does not exist.
	ErrDiscountNotFound = errors.NotFound("DISCOUNT_NOT_FOUND", "discount not found")
	// ErrDiscountInvalidArgument is returned when a discount request is invalid.
	ErrDiscountInvalidArgument = errors.BadRequest("DISCOUNT_INVALID_ARGUMENT", "invalid discount argument")
	// ErrDiscountInactive is returned when a discount is not active.
	ErrDiscountInactive = errors.BadRequest("DISCOUNT_INACTIVE", "discount is not active")
	// ErrDiscountExpired is returned when a discount has expired.
	ErrDiscountExpired = errors.BadRequest("DISCOUNT_EXPIRED", "discount has expired")
	// ErrDiscountUsageLimit is returned when a discount usage limit is reached.
	ErrDiscountUsageLimit = errors.BadRequest("DISCOUNT_USAGE_LIMIT", "discount usage limit reached")
	// ErrDiscountMinAmount is returned when the order amount is below the minimum.
	ErrDiscountMinAmount = errors.BadRequest("DISCOUNT_MIN_AMOUNT", "order amount below minimum")
	// ErrDiscountScopeMismatch is returned when the discount scope doesn't match.
	ErrDiscountScopeMismatch = errors.BadRequest("DISCOUNT_SCOPE_MISMATCH", "discount scope does not match")
)
