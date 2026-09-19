package order

import (
	"time"

	"payment/pkg/uuid"
)

// --- 待支付 ---

// PendingOrder 待支付状态：下单之后、付款之前。
type PendingOrder struct {
	OrderMeta           // 自动获得 ID()、CustomerID()、Meta() 与 isOrderState()
	expiresAt time.Time // 状态专属字段：支付截止时间
}

// NewOrder 下单：原价由调用方按商品明细汇总后传入，ttl 内不支付就会被关单。
func NewOrder(
	customerID int64,
	originalAmount int64,
	placedAt time.Time,
	ttl time.Duration,
) (PendingOrder, error) {
	if customerID <= 0 {
		return PendingOrder{}, ErrOrderCustomerRequired
	}
	if originalAmount <= 0 {
		return PendingOrder{}, ErrInvalidAmount
	}
	if ttl <= 0 {
		return PendingOrder{}, ErrInvalidPaymentWindow
	}

	return PendingOrder{
		OrderMeta: OrderMeta{
			id:             uuid.NewInt64(),
			customerID:     customerID,
			originalAmount: originalAmount,
			payableAmount:  originalAmount,
			placedAt:       placedAt,
		},
		expiresAt: placedAt.Add(ttl),
	}, nil
}

// ReconstitutePending 仓储恢复：字段全部由外部给定，不做校验。
func ReconstitutePending(meta OrderMeta, expiresAt time.Time) PendingOrder {
	return PendingOrder{OrderMeta: meta, expiresAt: expiresAt}
}

func (o PendingOrder) ExpiresAt() time.Time { return o.expiresAt }

// IsPending 覆写元数据里的默认值（false）：待支付就是本状态。
func (o PendingOrder) IsPending() bool { return true }

// ApplyDiscount 记一笔优惠（券 / 活动算好之后传进来），返回更新后的订单。
//
// 优惠额对不对（券可用、门槛、有效期）由 Coupon 聚合判定，订单只认「钱数」：
// 单笔为正、累计不超过原价。**一张订单可以用多张券** —— 逐次调用本方法累加。
func (o PendingOrder) ApplyDiscount(discountAmount int64) (PendingOrder, error) {
	if discountAmount <= 0 || discountAmount > o.originalAmount-o.discountAmount {
		return o, ErrInvalidDiscount
	}

	o.discountAmount += discountAmount
	o.payableAmount = o.originalAmount - o.discountAmount
	return o, nil
}

// ClearDiscount 清空优惠（换券、取消用券时调用）。
//
// 券本身要回到可用状态由调用方负责（Coupon.Release）—— 订单不认识券。
func (o PendingOrder) ClearDiscount() PendingOrder {
	o.discountAmount = 0
	o.payableAmount = o.originalAmount
	return o
}

// HasDiscount 是否已经有优惠（用了券或参加了活动）。
func (o PendingOrder) HasDiscount() bool { return o.discountAmount > 0 }

// MarkPaid 支付成功 → PaidOrder。
//
// 三个硬性约束：超时的订单不能再支付（要先 CloseIfExpired 关掉）；
// 实付金额必须与应付一致（部分支付不属于本聚合的能力）；
// 必须说清**走的是哪条支付途径**（对账、退款原路返回都要靠它）。
func (o PendingOrder) MarkPaid(
	paidAmount int64,
	payChannel PayChannel,
	tradeNo string,
	now time.Time,
) (PaidOrder, error) {
	if !now.Before(o.expiresAt) {
		return PaidOrder{}, ErrOrderExpired
	}
	if paidAmount != o.payableAmount {
		return PaidOrder{}, ErrPaymentAmountMismatch
	}
	if !payChannel.valid() {
		return PaidOrder{}, ErrInvalidPayChannel
	}

	return PaidOrder{
		OrderMeta:  o.OrderMeta,
		paidAmount: paidAmount,
		paidAt:     now,
		tradeNo:    tradeNo,
		payChannel: payChannel,
	}, nil
}

// Cancel 主动取消 → CanceledOrder（reason 可以留空）。
func (o PendingOrder) Cancel(reason string, now time.Time) CanceledOrder {
	return CanceledOrder{
		OrderMeta:  o.OrderMeta,
		canceledAt: now,
		reason:     reason,
	}
}

// CloseIfExpired 超时关单 → ClosedOrder；没到期则报错，避免被误关。
func (o PendingOrder) CloseIfExpired(now time.Time) (ClosedOrder, error) {
	if now.Before(o.expiresAt) {
		return ClosedOrder{}, ErrOrderNotExpired
	}
	return ClosedOrder{
		OrderMeta: o.OrderMeta,
		closedAt:  now,
	}, nil
}

// IsExpired 是否已过支付时限（只读，不改状态）。
func (o PendingOrder) IsExpired(now time.Time) bool { return !now.Before(o.expiresAt) }
