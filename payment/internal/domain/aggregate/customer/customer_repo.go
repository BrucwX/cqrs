package customer

import "context"

type CustomerCommand interface {
	MustGetCustomer(ctx context.Context, customerID int64) (Customer, error)
}
