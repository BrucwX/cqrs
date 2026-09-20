package repo

import (
	"context"

	"payment/internal/domain/aggregate/customer"
)

// CustomerCommand 客户的写侧仓储。
type CustomerCommand interface {
	MustGetCustomer(ctx context.Context, customerID int64) (customer.Customer, error)
	// MustGetCustomers 按一批 id 取**客户集合**（批量取）：比如「这批行项的受赠人是谁」，
	// 一次把人都取回来，别在循环里一个个查。
	//
	// 口径跟单个版本一致：有一个 id 查不到就报 ErrCustomerNotFound。
	MustGetCustomers(ctx context.Context, customerIDs []int64) (customer.Customers, error)
}
