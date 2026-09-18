package order

import (
	"time"

	"payment/pkg/uuid"
)

// --- 聚合根 (Aggregate Root) ---

// Order 是一次购买行为的订单聚合根。
//
// 金额一律用「分」的 int64 表示。买了哪些商品不在订单上，而在 orderItem
// （订单行项）；订单只记钱与状态。
//
// 优惠只记**累计金额**，不记用了哪几张券：券能不能用归 Coupon 判定，核销关系
// 记在券自己的 orderID 上（一张订单可以用多张券，逐次调 ApplyDiscount 累加）。
//
// 状态机（只允许下面这几条边）：
//
//	Pending ──MarkPaid──► Paid ──Refund──► Refunded
//	   ├──Cancel────────► Canceled
//	   └──CloseIfExpired─► Closed
type Order struct {
	id      int64
	buyerID int64

	originalAmount int64 // 原价（分）
	discountAmount int64 // 优惠金额（分）
	payableAmount  int64 // 应付 = 原价 - 优惠
	paidAmount     int64 // 实付（分）

	status Status

	placedAt   time.Time // 下单时间
	expiresAt  time.Time // 支付截止时间
	paidAt     time.Time
	canceledAt time.Time
	closedAt   time.Time
}

// NewOrder 下单：原价由调用方按商品明细汇总后传入，ttl 内不支付就会被关单。
func NewOrder(
	buyerID int64,
	originalAmount int64,
	placedAt time.Time,
	ttl time.Duration,
) (*Order, error) {
	if buyerID <= 0 {
		return nil, ErrOrderBuyerRequired
	}
	if originalAmount <= 0 {
		return nil, ErrInvalidAmount
	}
	if ttl <= 0 {
		return nil, ErrInvalidPaymentWindow
	}

	return &Order{
		id:             uuid.NewInt64(),
		buyerID:        buyerID,
		originalAmount: originalAmount,
		payableAmount:  originalAmount,
		status:         StatusPending,
		placedAt:       placedAt,
		expiresAt:      placedAt.Add(ttl),
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(
	id int64,
	buyerID int64,
	originalAmount int64,
	discountAmount int64,
	payableAmount int64,
	paidAmount int64,
	status Status,
	placedAt time.Time,
	expiresAt time.Time,
	paidAt time.Time,
	canceledAt time.Time,
	closedAt time.Time,
) *Order {
	return &Order{
		id:             id,
		buyerID:        buyerID,
		originalAmount: originalAmount,
		discountAmount: discountAmount,
		payableAmount:  payableAmount,
		paidAmount:     paidAmount,
		status:         status,
		placedAt:       placedAt,
		expiresAt:      expiresAt,
		paidAt:         paidAt,
		canceledAt:     canceledAt,
		closedAt:       closedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// ApplyDiscount 记一笔优惠（券 / 活动算好之后传进来），并重算应付金额。
//
// 优惠额对不对（券可用、门槛、有效期）由 Coupon 聚合判定，订单只认「钱数」：
// 单笔为正、累计不超过原价。**一张订单可以用多张券** —— 逐次调用本方法累加。
func (o *Order) ApplyDiscount(discountAmount int64) error {
	if o.status != StatusPending {
		return ErrOrderNotPending
	}
	if discountAmount <= 0 || discountAmount > o.originalAmount-o.discountAmount {
		return ErrInvalidDiscount
	}

	o.discountAmount += discountAmount
	o.payableAmount = o.originalAmount - o.discountAmount
	return nil
}

// ClearDiscount 清空优惠（换券、取消用券时调用）。
//
// 券本身要回到可用状态由调用方负责（Coupon.Release）—— 订单不认识券。
func (o *Order) ClearDiscount() error {
	if o.status != StatusPending {
		return ErrOrderNotPending
	}
	o.discountAmount = 0
	o.payableAmount = o.originalAmount
	return nil
}

// MarkPaid 支付成功。
//
// 两个硬性约束：超时的订单必须先 CloseIfExpired 关掉，不能直接支付；
// 实付金额必须与应付一致（部分支付不属于本聚合的能力）。
func (o *Order) MarkPaid(paidAmount int64, now time.Time) error {
	if o.status != StatusPending {
		return ErrOrderNotPending
	}
	if !now.Before(o.expiresAt) {
		return ErrOrderExpired
	}
	if paidAmount != o.payableAmount {
		return ErrPaymentAmountMismatch
	}

	o.paidAmount = paidAmount
	o.paidAt = now
	o.status = StatusPaid
	return nil
}

// Cancel 用户/后台主动取消，仅限未支付订单。
func (o *Order) Cancel(now time.Time) error {
	if o.status != StatusPending {
		return ErrOrderNotPending
	}
	o.status = StatusCanceled
	o.canceledAt = now
	return nil
}

// CloseIfExpired 超时关单（由定时任务调用）；没到期则报错，避免被误关。
func (o *Order) CloseIfExpired(now time.Time) error {
	if o.status != StatusPending {
		return ErrOrderNotPending
	}
	if now.Before(o.expiresAt) {
		return ErrOrderNotExpired
	}
	o.status = StatusClosed
	o.closedAt = now
	return nil
}

// Refund 把已支付订单推到「已退款」。
//
// 退款单本身的细节（金额、审批、打款时间）在 refund 聚合里，订单只记状态，
// 所以这里不需要时间参数。
func (o *Order) Refund() error {
	if o.status != StatusPaid {
		return ErrOrderNotPaid
	}
	o.status = StatusRefunded
	return nil
}

// IsExpired 订单是否已过支付时限（只读，不改状态）。
func (o *Order) IsExpired(now time.Time) bool {
	return o.status == StatusPending && !now.Before(o.expiresAt)
}

// HasDiscount 是否已经有优惠（用了券或参加了活动）。
func (o *Order) HasDiscount() bool { return o.discountAmount > 0 }

// --- 属性只读访问器 (Getters) ---

func (o *Order) ID() int64             { return o.id }
func (o *Order) BuyerID() int64        { return o.buyerID }
func (o *Order) OriginalAmount() int64 { return o.originalAmount }
func (o *Order) DiscountAmount() int64 { return o.discountAmount }
func (o *Order) PayableAmount() int64  { return o.payableAmount }
func (o *Order) PaidAmount() int64     { return o.paidAmount }
func (o *Order) Status() Status        { return o.status }
func (o *Order) PlacedAt() time.Time   { return o.placedAt }
func (o *Order) ExpiresAt() time.Time  { return o.expiresAt }
func (o *Order) PaidAt() time.Time     { return o.paidAt }
func (o *Order) CanceledAt() time.Time { return o.canceledAt }
func (o *Order) ClosedAt() time.Time   { return o.closedAt }
