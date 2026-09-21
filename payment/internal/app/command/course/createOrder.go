package course

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
const defaultPaymentWindow = 3 * time.Minute

// CreateCOrderIn 下单入参。
type CreateCOrderIn struct {
	CustomerID     int64               // 下单的客户（不是行上的受赠人；原价按行项合计）
	OrderItemIDs   []int64             // 已经拼好的草稿行（购物车式的行项）
	Channel        orderagg.PayChannel // 想走的支付途径（微信 / 支付宝 / 银行卡…）
	Scene          string              // 渠道内场景：app / h5 / native / jsapi（可选，渠道自己解读）
	IdempotencyKey string              // 等幂key
}

// CreateCOrderOut 下单结果：订单号 + 「怎么付」。
type CreateCOrderOut struct {
	OrderID int64             // 订单号：前端轮询支付状态、重新发起支付都靠它
	Pay     appRepo.PayResult // 二维码 / 跳转链接 / 预支付参数（交给前端）
}

// CreateCOrder 下单 + 发起支付：把草稿行变成一张待支付订单，再调渠道侧 Prepare，
// 把订单号与「怎么付」（二维码 / URL / 预支付参数）一起交给前端。
//
// 一步到位是刻意的：前端流程就是「选好商品 + 选好支付方式 → 一次请求」，
// 不该为了拿二维码再调一次；而且订单**先落库再发起支付**，渠道回调回来才找得到单。
func (h *Handler) CreateCOrder(ctx context.Context, in CreateCOrderIn) (out CreateCOrderOut, err error) {
	if err = h.idservice.CheckIdempotency(in.IdempotencyKey); err != nil {
		return
	}

	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return CreateCOrderOut{}, err
	}
	defer func() {
		if err != nil {
			_ = h.tx.Rollback(ctx)
		}
	}()

	// 取数：要买的行项
	items, err := h.orderItemCommand.MustGetOrderItems(ctx, in.OrderItemIDs)
	if err != nil {
		return CreateCOrderOut{}, err
	}

	// 顺序要紧：先挑出「这批行项还能不能用」（全草稿），再动课程那边的副作用（报名 / 减库存），
	// 免得为一批用不了的行项白报名。
	if !items.AllDraft() {
		return CreateCOrderOut{}, orderitem.ErrOrderItemAlreadyConfirmed
	}

	// 能不能买（只买课：课程侧校验并报名）
	if err = h.CheckCanBuy(ctx, items); err != nil {
		return CreateCOrderOut{}, err
	}

	// 建单：原价按行项合计（先不考虑优惠）。
	now := time.Now()
	pending, err := orderagg.NewOrder(in.CustomerID, items.Total(), now, defaultPaymentWindow)
	if err != nil {
		return CreateCOrderOut{}, err
	}

	// 先归单（领域：草稿行 → 正式行），再把这批行项存下去（存储只管写）。
	confirmed, err := items.ConfirmAll(pending.ID())
	if err != nil {
		return CreateCOrderOut{}, err
	}
	if err := h.orderItemCommand.Save(ctx, confirmed); err != nil {
		return CreateCOrderOut{}, err
	}
	if err := h.orderCommand.Save(ctx, pending); err != nil {
		return CreateCOrderOut{}, err
	}

	// 开启超时关单：支付窗口过了还没付，就把这张订单关掉（异步，不在本事务里）。
	// id 带订单号：之后可以按 id 精确取消（比如用户提前付了、或主动取消）。
	if _, err = h.scheduler.AfterFunc(
		fmt.Sprintf("close-unpaid-order-%d", pending.ID()),
		defaultPaymentWindow,
		func(ctx context.Context) error {
			return h.closeIfUnpaid(ctx, pending.ID())
		},
	); err != nil {
		return CreateCOrderOut{}, err
	}

	// 发起支付，把「怎么付」拿回来交给前端。发起失败不影响刚建好的订单：
	// 它还是待支付，上面那个任务到点会关掉它（重新发起支付以后需要的话另开用例）。
	payResult, err := h.pay.Prepare(ctx, appRepo.PayIntent{
		OrderID:    pending.ID(),
		CustomerID: pending.CustomerID(),
		Amount:     pending.PayableAmount(),
		Channel:    in.Channel,
		Scene:      in.Scene,
	})
	if err != nil {
		return CreateCOrderOut{}, err
	}

	if err = h.tx.Commit(ctx); err != nil {
		return CreateCOrderOut{}, err
	}

	return CreateCOrderOut{OrderID: pending.ID(), Pay: payResult}, nil
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
