package repo

import (
	"context"

	"payment/internal/domain/aggregate/order"
)

// PayCommand 支付端口：把一张订单付掉（让订单从待支付变成已支付）。
//
// ⚠️ 它不是某个聚合的存储接口：实现要调支付渠道，再把结果落库。
type PayCommand interface {
	Pay(ctx context.Context, customerID int64, orderID int64, paychannel order.PayChannel) error
}
