package autoClosePay

import (
	"context"
	"time"

	orderagg "payment/internal/domain/aggregate/order"
	"payment/internal/domain/aggregate/product"
)

// CloseIfUnpaid 把一张过了支付窗口还没付的订单关掉：本地关单 → 渠道侧关单 → 课程侧退课。
//
// 幂等：已经不是待支付（付了 / 取消 / 已关单）或还没到点，直接返回 nil ——
// 队列重复投递、定时器提前触发，都不会办第二遍事。
//
// ⚠️ 顺序是「先本地、后远程」：本地一进终态，晚到的支付回调就找不到「待支付」可改，
// 由 Confirm 那边按资金差异处理；反过来（先关渠道）会出现「渠道关了、本地还挂着待支付」。
func (h *closePayScheduler) CloseIfUnpaid(ctx context.Context, orderID int64) error {
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

	// 退课要的数据在行项上，先取出来（取不到行项就只是关单，见最后）。
	items, err := h.itemCommand.ListByOrder(ctx, orderID)
	if err != nil {
		return err
	}

	closed, err := pending.CloseIfExpired(now)
	if err != nil {
		return err
	}

	// 本地关单落库：订单一进终态，晚到的支付回调就改不动它了。
	txCtx, err := h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	if err := h.orderCommand.Save(txCtx, closed); err != nil {
		_ = h.tx.Rollback(txCtx)
		return err
	}
	if err := h.tx.Commit(txCtx); err != nil {
		return err
	}

	// 渠道侧也关一下（微信 closeorder；支付宝自动关 / 无副作用）。
	if err := h.pay.Close(ctx, orderID); err != nil {
		return err
	}

	if len(items) == 0 {
		return nil // 没有行项 = 没有报名要退（老数据 / 只关单）
	}

	// 行项上的受赠人 → 账号 id；行项上的商品 → 课程 id。
	receivers, err := h.customerCommand.MustGetCustomers(ctx, items.CustomerIDs())
	if err != nil {
		return err
	}
	products, err := h.productCommand.GetProsByItems(ctx, items)
	if err != nil {
		return err
	}

	// 只退课程类：这批行项里混了别的类型就不硬塞给课程侧。
	courses := products.GroupByType()[product.ProductType]
	if len(courses) == 0 {
		return nil
	}

	return h.courseService.UnenrollCourse(ctx, receivers.AccountIDs(), courses.TargetIds())
}
