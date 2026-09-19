package order

import (
	"context"
	"time"

	orderagg "payment/internal/domain/aggregate/order"
)

// defaultPaymentWindow 没传支付时限时的默认值：30 分钟不付就会被关单。
const defaultPaymentWindow = 30 * time.Minute

// CreateOrderIn 下单入参。
type CreateOrderIn struct {
	CustomerID    int64         // 下单的客户（不是行上的受赠人；原价按行项合计）
	OrderItemIDs  []int64       // 已经拼好的草稿行（购物车式的行项）
	PaymentWindow time.Duration // 支付时限，<= 0 用默认 30 分钟
}

// CreateOrder 下单：把一批**草稿行**变成一张待支付订单。
//
// 顺序：
//
//  1. 预检（domain/service/orderItemsCheck）：都还是草稿、库存够、受赠人能收
//  2. 建单：原价 = 这批行项的小计合计，状态是「待支付」
//  3. 行项归单：草稿行 Confirm 到这张订单上（值语义，各自返回新的）
//  4. 落库：行项 + 订单
//
// ⚠️ 这里**不扣库存**：下单只是意向，等订单**支付成功**之后再扣（那时由支付用例
// 去 ReserveStock）。所以上面那个库存预检只是「提前说卖完了」，不构成预留。
//
// 全程一个事务：任何一步出错都整体回滚 —— 「行项归了单但订单没建出来」这种半成品不能留。
func (h *Handler) CreateOrder(ctx context.Context, in CreateOrderIn) (err error) {
	if len(in.OrderItemIDs) == 0 {
		return ErrOrderItemsRequired
	}

	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer h.tx.Rollback(ctx) // 事务没提交就回滚

	items, err := h.orderItemCommand.MustGetOrderItems(ctx, in.OrderItemIDs)
	if err != nil {
		return err
	}

	// 预检：规则都在 domain/service 里，这里只问一句「行不行」。
	ok, err := h.itemsCheck.CanOrderCreate(ctx, in.OrderItemIDs)
	if err != nil {
		return err
	}
	if !ok {
		return ErrOrderNotCreatable
	}

	now := time.Now()
	window := in.PaymentWindow
	if window <= 0 {
		window = defaultPaymentWindow
	}

	// 建单：原价按行项合计，优惠留到之后（券的抵扣由 coupon 聚合先算好再 ApplyDiscount）。
	pending, err := orderagg.NewOrder(in.CustomerID, items.Total(), now, window)
	if err != nil {
		return err
	}

	// 先归单（领域：草稿行 → 正式行），再把这批行项存下去（存储只管写）。
	confirmed, err := items.ConfirmAll(pending.ID())
	if err != nil {
		return err
	}
	if err := h.orderItemCommand.Save(ctx, confirmed); err != nil {
		return err
	}

	if err := h.orderCommand.Save(ctx, pending); err != nil {
		return err
	}

	return h.tx.Commit(ctx)
}
