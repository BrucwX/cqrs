package order

import (
	"github.com/go-kratos/kratos/v3/errors"
)

var (
	ErrOrderInvalidArgument = errors.BadRequest("ORDER_INVALID_ARGUMENT", "invalid order argument")
	ErrOrderNotFound        = errors.NotFound("ORDER_NOT_FOUND", "order not found")
	ErrOrderAlreadyPaid     = errors.BadRequest("ORDER_ALREADY_PAID", "order already paid")
	ErrOrderCanceled        = errors.BadRequest("ORDER_CANCELED", "order is canceled")
)
