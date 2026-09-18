package couponType

import (
	"strings"
	"time"

	"payment/pkg/uuid"
)

// --- 聚合根 (Aggregate Root) ---

// CouponType 是「一种优惠券」的聚合根：券的权益规格、适用商品类别、有效期规则与发放配额。
//
// 它管「能不能发、发出去长什么样」；一张券发出去之后的生命周期（核销、退回、
// 过期）归 coupon.Coupon 管 —— 两个聚合靠 ID 单向引用，方向是
// **coupon 包 import 本包**（券持有本包的定义），所以本包不认识
// coupon.Coupon，发券时只把有效期窗口交出去。
//
// 配额意味着这里会有并发抢占（和课程名额 Course.ReserveSeat 一个道理），
// 所以发券必须经过本聚合，不能由 app 层自己造券。
type CouponType struct {
	id       int64
	name     string
	benefit  Benefit  // 券的权益规格，发券时快照进每一张券
	validity Validity // 有效期规则：固定窗口 / 领取后 N 天
	scope    Scope    // 可用的商品类别：全部 / 限定某几个类别
	quota    Quota    // 发放配额：是否限量 / 发行总量 / 已发数量
	status   Status
}

// NewCouponType 建券种类：默认「未上架」，要 Activate 之后才能发券。
func NewCouponType(
	name string,
	benefit Benefit,
	validity Validity,
	scope Scope,
	quota Quota,
) (*CouponType, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrCouponTypeNameRequired
	}
	if !benefit.Kind().Valid() {
		return nil, ErrCouponTypeBenefitRequired
	}
	if !validity.valid() {
		return nil, ErrInvalidValidityRule
	}
	if !scope.valid() {
		return nil, ErrInvalidScope
	}
	if !quota.valid() {
		return nil, ErrInvalidQuota
	}

	return &CouponType{
		id:       uuid.NewInt64(),
		name:     name,
		benefit:  benefit,
		validity: validity,
		scope:    scope,
		quota:    quota,
		status:   StatusDraft,
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(
	id int64,
	name string,
	benefit Benefit,
	validity Validity,
	scope Scope,
	quota Quota,
	status Status,
) *CouponType {
	return &CouponType{
		id:       id,
		name:     name,
		benefit:  benefit,
		validity: validity,
		scope:    scope,
		quota:    quota,
		status:   status,
	}
}

// --- 上架 / 下架 ---

// Activate 上架：只有上架的券种类能发券。
func (t *CouponType) Activate() error {
	if t.status == StatusActive {
		return ErrCouponTypeAlreadyActive
	}
	t.status = StatusActive
	return nil
}

// Deactivate 下架：不再发新券，已经发出去的券不受影响。
func (t *CouponType) Deactivate() error {
	if t.status != StatusActive {
		return ErrCouponTypeNotActive
	}
	t.status = StatusInactive
	return nil
}

// --- 发券 ---

// Issue 发券：占一个名额，按有效期规则算出这张券的有效期窗口。
//
// 本包不认识具体券，所以只把窗口交出去，建券由调用方（app 层）完成：
//
//	from, to, err := ct.Issue(now)
//	c, err := coupon.NewCoupon(ownerID, ct.ID(), from, to)
//
// 名额（quota.issued）与建出来的券必须在同一个事务里落库，否则会出现「名额扣了、券没发出去」。
func (t *CouponType) Issue(now time.Time) (time.Time, time.Time, error) {
	if t.status != StatusActive {
		return time.Time{}, time.Time{}, ErrCouponTypeNotIssuable
	}

	from, to := t.validity.Resolve(now)
	if !now.Before(to) {
		// 固定窗口的活动已经结束；相对窗口不会命中（days > 0）。
		return time.Time{}, time.Time{}, ErrCouponTypeExpired
	}

	// 先判完再占名额：Reserve 返回新配额，失败不产生副作用。
	next, err := t.quota.Reserve(1)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	t.quota = next

	return from, to, nil
}

// ChangeScope 改适用商品类别（比如从「限某类」放开成通用）。
//
// ⚠️ 适用商品类别是**活规则**：券上只存 couponTypeID、不留副本，核销时按本聚合当前的
// 范围判，所以改范围会立刻影响到已经发出去的券。要「发出去就锁定」得在券上存快照。
func (t *CouponType) ChangeScope(scope Scope) error {
	if !scope.valid() {
		return ErrInvalidScope
	}
	t.scope = scope
	return nil
}

// CoversProductType 某类商品能不能用这种券（核销前的准入判定）。
func (t *CouponType) CoversProductType(productTypeID int64) bool {
	return t.scope.Covers(productTypeID)
}

// IsSoldOut 是否已发完（不限量永远不发完）。
//
// 留个便利方法，等价于 `Quota().IsExhausted()`。
func (t *CouponType) IsSoldOut() bool {
	return t.quota.IsExhausted()
}

// --- 属性只读访问器 (Getters) ---

func (t *CouponType) ID() int64          { return t.id }
func (t *CouponType) Name() string       { return t.name }
func (t *CouponType) Benefit() Benefit   { return t.benefit }
func (t *CouponType) Validity() Validity { return t.validity }
func (t *CouponType) Scope() Scope       { return t.scope }
func (t *CouponType) Quota() Quota       { return t.quota }
func (t *CouponType) Status() Status     { return t.status }
func (t *CouponType) IsActive() bool     { return t.status == StatusActive }
