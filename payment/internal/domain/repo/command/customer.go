package command

import (
	"context"
	"payment/internal/domain/aggregate/customer"
)

type CustomerCommand interface {
	MustGetCustomer(ctx context.Context, customerID int64) (customer.Customer, error)
}
