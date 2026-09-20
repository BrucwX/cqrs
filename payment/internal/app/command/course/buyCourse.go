package course

import (
	"context"
	"time"

	"payment/internal/domain/aggregate/customer"
	orderagg "payment/internal/domain/aggregate/order"
	orderitem "payment/internal/domain/aggregate/orderItem"
)

// defaultPaymentWindow 没传支付时限时的默认值：30 分钟不付就会被关单。
const defaultPaymentWindow = 30 * time.Minute

// CreateOrderIn 下单入参。
type BuyCourseIn struct {
	CustomerID   int64   // 下单的客户（不是行上的受赠人；原价按行项合计）
	OrderItemIDs []int64 // 已经拼好的草稿行（购物车式的行项）
}

func (h *Handler) Buy(ctx context.Context, in BuyCourseIn) (err error) {
	// 取数
	items, err := h.orderItemCommand.MustGetOrderItems(ctx, in.OrderItemIDs)
	if err != nil {
		return err
	}
	customer ,err = h.customerCommand.MustGetCustomer(ctx , in.CustomerID)
	if err != nil {
		return err
	}

	// 检测能不能买，以及减库存
	h.CheckCanBuy(ctx , items , customer)

	// 创建订单
	po, err := h.CreateCourseOrder(ctx, in.CustomerID, items)

	return nil
}

func (h *Handler) CheckCanBuy(ctx context.Context, items orderitem.OrderItems , customer customer.Customer)(bool,err){
	cids := items.CustomerIDs()

	err := h.courseService.CheckAndEnrollCourse(ctx ,,)
}

// CreateOrder 下单：把一批**草稿行**变成一张待支付订单。
func (h *Handler) CreateCourseOrder(ctx context.Context, cID int64, items orderitem.OrderItems) (po *orderagg.PendingOrder, err error) {

	// 订单不是全draft
	if !items.AllDraft() {
		return nil, orderitem.ErrOrderItemAlreadyConfirmed
	}

	// 建单：原价按行项合计，优惠留到之后（券的抵扣由 coupon 聚合先算好再 ApplyDiscount）。
	now := time.Now()
	window := defaultPaymentWindow
	pending, err := orderagg.NewOrder(cID, items.Total(), now, window)
	if err != nil {
		return nil, err
	}

	// 先归单（领域：草稿行 → 正式行），再把这批行项存下去（存储只管写）。
	confirmed, err := items.ConfirmAll(pending.ID())
	if err != nil {
		return nil, err
	}
	if err := h.orderItemCommand.Save(ctx, confirmed); err != nil {
		return nil, err
	}

	if err := h.orderCommand.Save(ctx, pending); err != nil {
		return nil, err
	}

	return &pending, nil
}
