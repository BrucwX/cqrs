package buy

import (
	"context"
	"time"

	orderagg "payment/internal/domain/aggregate/order"
	orderitem "payment/internal/domain/aggregate/orderItem"
	"payment/internal/domain/aggregate/product"
)

// defaultPaymentWindow 没传支付时限时的默认值：30 分钟不付就会被关单。
const defaultPaymentWindow = 30 * time.Minute

// CreateOrderIn 下单入参。
type BuyCourseIn struct {
	CustomerID   int64   // 下单的客户（不是行上的受赠人；原价按行项合计）
	OrderItemIDs []int64 // 已经拼好的草稿行（购物车式的行项）
	pc           orderagg.PayChannel
}

func (h *Handler) BuyCourse(ctx context.Context, in BuyCourseIn) (err error) {
	h.tx.Begin(ctx)

	defer h.tx.Rollback(ctx)

	// 取数
	items, err := h.orderItemCommand.MustGetOrderItems(ctx, in.OrderItemIDs)
	if err != nil {
		return err
	}
	c, err2 := h.customerCommand.MustGetCustomer(ctx, in.CustomerID)
	if err2 != nil {
		return err2
	}

	// 检测能不能买，以及减库存
	err = h.CheckCanBuy(ctx, items)
	if err != nil {
		return err
	}

	// 创建订单
	po, err := h.CreateOrder(ctx, in.CustomerID, items)

	// 支付
	h.pay.Pay(ctx, c.ID(), po.ID(), in.pc)

	return h.tx.Commit(ctx)
}

// CheckCanBuy 买之前的校验：把 items 涉及到的商品和客户都取出来，交给课程那边校验报名。
//
// ⚠️ 这个用例**只买课**：混进别的商品类型就直接报错，不做「按类型分派」那一套。
// 哪天真要多卖几类，再来考虑每种类型一个下游能力的写法。
//
// 取数一次取一批（别在循环里一个个查）：商品按行项取（GetProsByItems），
// 客户按行项上的受赠人 id 取（MustGetCustomers）。
func (h *Handler) CheckCanBuy(ctx context.Context, items orderitem.OrderItems) error {
	// items 里涉及到的所有商品。
	products, err := h.productCommand.GetProsByItems(ctx, items)
	if err != nil {
		return err
	}

	// 只买课：这批商品必须都是课程类。
	if !products.AllOfType(product.ProductType) {
		return product.ErrUnsupportedProductType
	}

	// items 里涉及到的所有客户。
	receivers, err := h.customerCommand.MustGetCustomers(ctx, items.CustomerIDs())
	if err != nil {
		return err
	}

	// 课程那边：校验并报名（账号 id 与课程 id 都是 string）。
	if err := h.courseService.CheckAndEnrollCourse(ctx, receivers.AccountIDs(), products.TargetIds()); err != nil {
		return err
	}

	return nil
}

// CreateOrder 下单：把一批**草稿行**变成一张待支付订单。
func (h *Handler) CreateOrder(ctx context.Context, cID int64, items orderitem.OrderItems) (po *orderagg.PendingOrder, err error) {

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
