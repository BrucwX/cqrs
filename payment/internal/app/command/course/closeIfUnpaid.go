package course

import (
	"context"
	"time"

	orderagg "payment/internal/domain/aggregate/order"
)

// CloseCourseIfUnpaid 超时关单：延时队列唤醒后逐条处理。
//
// 订单还在支付窗口内、或已经不是待支付了（付了 / 取消 / 已关单），都直接返回 nil。
func (h *Handler) CloseCourseIfUnpaid(ctx context.Context, orderID int64) error {
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
	items, err := h.orderItemCommand.ListByOrder(ctx, orderID)
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
	// （"offlineCourse" = 课程商品的 Type()，跟下单那边入队用的是同一个字面量。）
	courses := products.GroupByType()["offlineCourse"]
	if len(courses) == 0 {
		return nil
	}

	return h.courseService.UnenrollCourse(ctx, receivers.AccountIDs(), courses.TargetIds())
}
