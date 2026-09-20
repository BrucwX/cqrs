package course

import (
	"context"
	"time"

	orderagg "payment/internal/domain/aggregate/order"
)

// defaultPaymentWindow 没传支付时限时的默认值：30 分钟不付就会被关单。
const defaultPaymentWindow = 30 * time.Minute

// CreateOrderIn 下单入参。
type CreateOrderIn struct {
	CustomerID   int64   // 下单的客户（不是行上的受赠人；原价按行项合计）
	OrderItemIDs []int64 // 已经拼好的草稿行（购物车式的行项）
}

// CreateOrder 下单：把一批**草稿行**变成一张待支付订单。

func (h *Handler) CreateCourseOrder(ctx context.Context, in CreateOrderIn) (err error) {
	if len(in.OrderItemIDs) == 0 {
		return orderagg.ErrOrderItemsRequired
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
		return orderagg.ErrOrderNotCreatable
	}

	now := time.Now()
	window := defaultPaymentWindow

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
