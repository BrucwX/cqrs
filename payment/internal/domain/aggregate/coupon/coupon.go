package coupon

import (
	"strings"
	"time"

	"payment/pkg/uuid"

	"payment/internal/domain/aggregate/couponType"
)

// --- 聚合根 (Aggregate Root) ---

// Coupon 是「一张已经发到用户手上的优惠券」的聚合根。
//
// 券种类（权益规格、适用商品类别、有效期规则、发放配额）活在另一个聚合 couponType 里，
// 本聚合只存一个 couponTypeID —— **不复制权益与适用范围的副本**，要用的时候由调用方
// 按这个 ID 把种类查出来，一起传进判定方法。
//
// 券自己持有的事实只有：券码、持有人、有效期（发券时按种类规则算好，之后固定）、
// 生命周期状态与核销记录。
//
// 金额一律用「分」的 int64 表示，不引入浮点。
type Coupon struct {
	id           int64
	code         string // 券码：对外展示、核销时按它找券
	ownerID      int64  // 持有人
	couponTypeID int64  // 来源券种类：权益 / 适用范围都按它现查

	validFrom time.Time // 发券时按种类的有效期规则算出来的窗口
	validTo   time.Time
	status    Status // 生命周期状态
	orderID   int64  // 核销到哪张订单
	usedAt    time.Time
}

// NewCoupon 发券：由调用方拿券种类 Issue 出来的有效期窗口调过来，券从「未使用」开始。
//
// 校验只剩「券自己这半边」：持有人、来源种类、有效期窗口；权益与适用范围不落库。
func NewCoupon(
	ownerID int64,
	couponTypeID int64,
	validFrom time.Time,
	validTo time.Time,
) (*Coupon, error) {
	if ownerID <= 0 {
		return nil, ErrCouponOwnerRequired
	}
	if couponTypeID <= 0 {
		return nil, ErrCouponTypeRequired
	}
	if !validFrom.Before(validTo) {
		return nil, ErrInvalidValidity
	}

	return &Coupon{
		id:           uuid.NewInt64(),
		code:         generateCode(),
		ownerID:      ownerID,
		couponTypeID: couponTypeID,
		validFrom:    validFrom,
		validTo:      validTo,
		status:       StatusUnused,
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(
	id int64,
	code string,
	ownerID int64,
	couponTypeID int64,
	validFrom time.Time,
	validTo time.Time,
	status Status,
	orderID int64,
	usedAt time.Time,
) *Coupon {
	return &Coupon{
		id:           id,
		code:         code,
		ownerID:      ownerID,
		couponTypeID: couponTypeID,
		validFrom:    validFrom,
		validTo:      validTo,
		status:       status,
		orderID:      orderID,
		usedAt:       usedAt,
	}
}

// generateCode 生成券码（16 位大写十六进制）。
func generateCode() string {
	raw := strings.ReplaceAll(uuid.NewString(), "-", "")
	return strings.ToUpper(raw[:16])
}

// --- 核心领域行为 (Domain Behaviors) ---

// ensureRedeemable 判定券自身此刻是否处于可核销状态（不看金额、不看商品）。
func (c *Coupon) ensureRedeemable(now time.Time) error {
	switch c.status {
	case StatusUsed:
		return ErrCouponAlreadyUsed
	case StatusVoided:
		return ErrCouponVoided
	case StatusExpired:
		return ErrCouponExpired
	}
	if now.Before(c.validFrom) {
		return ErrCouponNotStarted
	}
	if !now.Before(c.validTo) {
		return ErrCouponExpired
	}
	return nil
}

// ensureUsable 判定券此刻能否用在「某类商品 + orderAmount」上（不改状态）。
//
// 券这半边看状态与有效期；种类那半边看适用商品类别与门槛 —— 所以种类要传进来，
// 不能靠券上存副本（券只留 couponTypeID）。
func (c *Coupon) ensureUsable(
	ct *couponType.CouponType,
	orderAmount int64,
	productTypeID int64,
	now time.Time,
) error {
	if err := c.ensureRedeemable(now); err != nil {
		return err
	}
	if orderAmount <= 0 {
		return ErrInvalidOrderAmount
	}
	if !ct.Scope().Covers(productTypeID) {
		return ErrCouponNotApplicable
	}
	if !ct.Benefit().MeetsThreshold(orderAmount) {
		return ErrThresholdNotMet
	}
	return nil
}

// DiscountFor 预演这张券在某类商品、某个订单金额上能减多少（只读，不改状态）。
func (c *Coupon) DiscountFor(
	ct *couponType.CouponType,
	orderAmount int64,
	productTypeID int64,
	now time.Time,
) (int64, error) {
	if err := c.ensureUsable(ct, orderAmount, productTypeID, now); err != nil {
		return 0, err
	}
	return ct.Benefit().DiscountFor(orderAmount), nil
}

// IsUsableAt 券此刻是否可用于「某类商品 + orderAmount」（只判不读优惠额，给「展示可用券列表」用）。
func (c *Coupon) IsUsableAt(
	ct *couponType.CouponType,
	orderAmount int64,
	productTypeID int64,
	now time.Time,
) bool {
	return c.ensureUsable(ct, orderAmount, productTypeID, now) == nil
}

// Redeem 核销：把券打到指定订单上，只做状态流转与登记 —— 优惠额由 DiscountFor 先算好
// （先判后写），这里不再算钱。
func (c *Coupon) Redeem(orderID int64, now time.Time) error {
	if orderID <= 0 {
		return ErrCouponOrderRequired
	}
	if err := c.ensureRedeemable(now); err != nil {
		return err
	}

	c.status = StatusUsed
	c.orderID = orderID
	c.usedAt = now
	return nil
}

// Release 退回已核销的券（订单取消/退款时调用）：没过期就回到「未使用」，过期就落到「已过期」。
func (c *Coupon) Release(now time.Time) error {
	if c.status != StatusUsed {
		return ErrCouponNotUsed
	}
	c.orderID = 0
	c.usedAt = time.Time{}
	if !now.Before(c.validTo) {
		c.status = StatusExpired
		return nil
	}
	c.status = StatusUnused
	return nil
}

// Expire 到期作废（由定时任务扫描后调用）。
func (c *Coupon) Expire(now time.Time) error {
	if c.status == StatusUsed || c.status == StatusVoided {
		return ErrCouponNotExpirable
	}
	if now.Before(c.validTo) {
		return ErrCouponNotExpired
	}
	c.status = StatusExpired
	return nil
}

// Void 人工作废（风控/发错券）。
func (c *Coupon) Void() error {
	if c.status == StatusVoided {
		return ErrCouponVoided
	}
	c.status = StatusVoided
	return nil
}

// EffectiveStatus 按当前时间看券的真实状态：过期的未使用券不必等定时任务也应是「已过期」。
func (c *Coupon) EffectiveStatus(now time.Time) Status {
	if c.status == StatusUnused && !now.Before(c.validTo) {
		return StatusExpired
	}
	return c.status
}

// --- 属性只读访问器 (Getters) ---

func (c *Coupon) ID() int64            { return c.id }
func (c *Coupon) Code() string         { return c.code }
func (c *Coupon) OwnerID() int64       { return c.ownerID }
func (c *Coupon) CouponTypeID() int64  { return c.couponTypeID }
func (c *Coupon) ValidFrom() time.Time { return c.validFrom }
func (c *Coupon) ValidTo() time.Time   { return c.validTo }
func (c *Coupon) Status() Status       { return c.status }
func (c *Coupon) OrderID() int64       { return c.orderID }
func (c *Coupon) UsedAt() time.Time    { return c.usedAt }
