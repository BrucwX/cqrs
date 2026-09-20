package repo

import (
	"context"

	"payment/internal/domain/aggregate/customer"
)

// CustomerCommand 客户的写侧仓储。
type CustomerCommand interface {
	MustGetCustomer(ctx context.Context, customerID int64) (customer.Customer, error)
}
