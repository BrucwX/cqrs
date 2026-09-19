package order

import "context"

type PayCommand interface {
	Pay(ctx context.Context, customerID int64, orderID int64, paychannel PayChannel) error
}
