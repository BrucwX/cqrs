package buy

import (
	"context"
	"fmt"
	"time"

	appRepo "payment/internal/app/command/appRepo"
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

func (h *Handler) BuyCourse(ctx context.Context, in BuyCourseIn) (result appRepo.PayResult, err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return appRepo.PayResult{}, err
	}
	defer func() {
		if err != nil {
			_ = h.tx.Rollback(ctx)
		}
	}()

	// 取数：要买的行项 + 下单的人
	items, err := h.orderItemCommand.MustGetOrderItems(ctx, in.OrderItemIDs)
	if err != nil {
		return appRepo.PayResult{}, err
	}
	c, err := h.customerCommand.MustGetCustomer(ctx, in.CustomerID)
	if err != nil {
		return appRepo.PayResult{}, err
	}

	// 能不能买（只买课：课程侧校验并报名）
	if err = h.CheckCanBuy(ctx, items); err != nil {
		return appRepo.PayResult{}, err
	}

	// 建单（草稿行 → 待支付订单）
	po, err := h.CreateOrder(ctx, in.CustomerID, items)
	if err != nil {
		return appRepo.PayResult{}, err
	}

	// 发起支付：把「怎么付」（二维码 / 跳转链接 / 预支付参数）交回给调用方。
	result, err = h.pay.Prepare(ctx, appRepo.PayIntent{
		OrderID:    po.ID(),
		CustomerID: c.ID(),
		Amount:     po.PayableAmount(),
		Channel:    in.pc,
	})
	if err != nil {
		return appRepo.PayResult{}, err
	}

	if err = h.tx.Commit(ctx); err != nil {
		return appRepo.PayResult{}, err
	}

	// 开启超时关单：支付窗口过了还没付，就把这张订单关掉（异步，不在本事务里）。
	// id 带订单号：之后可以按 id 精确取消（比如用户在窗口内提前付了、或主动取消）。
	if _, err := h.scheduler.AfterFunc(
		fmt.Sprintf("close-unpaid-order-%d", po.ID()),
		defaultPaymentWindow,
		func(ctx context.Context) error {
			return h.closeIfUnpaid(ctx, po.ID())
		},
	); err != nil {
		return appRepo.PayResult{}, err
	}

	return result, nil
}

// closeIfUnpaid 支付窗口到了还没付的订单，关掉它：
// 领域关单（CloseIfExpired）→ 落库 → 渠道侧关单（pay.Close）。
// 已经不是待支付（付了 / 取消 / 已关单）或还没到点，就什么都不做。
func (h *Handler) closeIfUnpaid(ctx context.Context, orderID int64) error {
	o, err := h.orderCommand.MustGetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	pending, ok := o.(orderagg.PendingOrder)
	if !ok {
		return nil // 付了 / 取消 / 已关单，没得可关
	}

	now := time.Now()
	if !pending.IsExpired(now) {
		return nil // 定时器提前触发，还没到点
	}

	closed, err := pending.CloseIfExpired(now)
	if err != nil {
		return err
	}
	if err := h.orderCommand.Save(ctx, closed); err != nil {
		return err
	}

	// 渠道侧也关一下（微信 closeorder；支付宝自动关 / 无副作用）。
	return h.pay.Close(ctx, orderID)
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

	// 建单：原价按行项合计。
	// 先不考虑优惠
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
