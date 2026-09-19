package command

import (
	"context"
	"payment/internal/domain/aggregate/order"
)

type PayCommand interface {
	Pay(ctx context.Context, customerID int64, orderID int64, paychannel order.PayChannel) error
}
