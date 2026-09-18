package couponType

import "time"

// --- 值对象 (Value Object) ---

// ValidityKind 有效期规则的种类。
type ValidityKind int

const (
	ValidityFixed    ValidityKind = iota + 1 // 固定窗口：全场统一在 from..to 之间有效
	ValidityRelative                         // 相对窗口：从领取那一刻起算 N 天
)

// Validity 是「券的有效期规则」值对象。
//
// 两种口径二选一：
//   - 固定窗口：活动券，统一在某个区间内有效；
//   - 相对窗口：领券时才算得出来，就是「领取后 N 天内有效」（按自然日，含领取当天）。
type Validity struct {
	kind ValidityKind
	from time.Time
	to   time.Time
	days int
}

// NewFixedValidity 固定窗口。
func NewFixedValidity(from, to time.Time) (Validity, error) {
	if !from.Before(to) {
		return Validity{}, ErrInvalidValidityRule
	}
	return Validity{kind: ValidityFixed, from: from, to: to}, nil
}

// NewRelativeValidity 领取后 days 天内有效。
func NewRelativeValidity(days int) (Validity, error) {
	if days <= 0 {
		return Validity{}, ErrInvalidValidityRule
	}
	return Validity{kind: ValidityRelative, days: days}, nil
}

// ReconstituteValidity 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteValidity(kind ValidityKind, from time.Time, to time.Time, days int) Validity {
	return Validity{kind: kind, from: from, to: to, days: days}
}

// Resolve 以 now（发券时刻）为基准算出券的有效期区间。
func (v Validity) Resolve(now time.Time) (time.Time, time.Time) {
	if v.kind == ValidityFixed {
		return v.from, v.to
	}
	return now, now.AddDate(0, 0, v.days)
}

// valid 是否是已知且自洽的有效期规则。
func (v Validity) valid() bool {
	switch v.kind {
	case ValidityFixed:
		return v.from.Before(v.to)
	case ValidityRelative:
		return v.days > 0
	default:
		return false
	}
}

// --- 属性只读访问器 (Getters) ---

func (v Validity) Kind() ValidityKind { return v.kind }
func (v Validity) From() time.Time    { return v.from }
func (v Validity) To() time.Time      { return v.to }
func (v Validity) Days() int          { return v.days }
