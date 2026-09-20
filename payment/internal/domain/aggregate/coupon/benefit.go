package coupon

import "strings"

// Kind 券类型，决定这份权益怎么算。
type Kind int

const (
	KindCash      Kind = iota + 1 // 无门槛立减：直接减 amount
	KindThreshold                 // 满减：订单金额达到 threshold 后减 amount
	KindDiscount                  // 折扣：按 percent 打折，最多减 maxDiscount
)

// Valid 是否是已知的券类型（0 值与非枚举值都算非法）。
func (k Kind) Valid() bool {
	return k >= KindCash && k <= KindDiscount
}

// --- 值对象 (Value Object) ---

// Benefit 是券的权益规格：这张券长什么样、能减多少。
//
// 它是**发券时快照进券里**的规格 —— 之后改规格、下架之类的动作都不影响已经发出去的券
// （券自己拿着这一份，要是想改这张券的面额，就改这张券）。
//
// 金额一律用「分」的 int64 表示。
type Benefit struct {
	title       string // 券名 / 发放来源说明
	kind        Kind   // 券类型：无门槛立减 / 满减 / 折扣
	amount      int64  // 立减、满减金额（分）
	threshold   int64  // 使用门槛（分）：0 表示无门槛
	percent     int32  // 折扣百分比 1..99，仅折扣券使用
	maxDiscount int64  // 折扣券封顶金额（分）：0 表示不封顶
}

// NewBenefit 构造权益规格，按类型校验字段：
//   - KindCash：amount > 0，threshold 必须为 0；
//   - KindThreshold：amount > 0，threshold > 0 且不低于 amount；
//   - KindDiscount：percent 在 1..99，amount/threshold 必须是 0，maxDiscount >= 0。
func NewBenefit(
	title string,
	kind Kind,
	amount int64,
	threshold int64,
	percent int32,
	maxDiscount int64,
) (Benefit, error) {
	if strings.TrimSpace(title) == "" {
		return Benefit{}, ErrBenefitTitleRequired
	}

	switch kind {
	case KindCash:
		if amount <= 0 {
			return Benefit{}, ErrInvalidAmount
		}
		if threshold != 0 {
			return Benefit{}, ErrInvalidThreshold
		}
	case KindThreshold:
		if threshold <= 0 || amount > threshold {
			return Benefit{}, ErrInvalidThreshold
		}
		if amount <= 0 {
			return Benefit{}, ErrInvalidAmount
		}
	case KindDiscount:
		if percent <= 0 || percent >= 100 {
			return Benefit{}, ErrInvalidPercent
		}
		if amount != 0 || threshold != 0 {
			return Benefit{}, ErrInvalidAmount
		}
		if maxDiscount < 0 {
			return Benefit{}, ErrInvalidAmount
		}
	default:
		return Benefit{}, ErrInvalidKind
	}

	return Benefit{
		title:       title,
		kind:        kind,
		amount:      amount,
		threshold:   threshold,
		percent:     percent,
		maxDiscount: maxDiscount,
	}, nil
}

// ReconstituteBenefit 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteBenefit(
	title string,
	kind Kind,
	amount int64,
	threshold int64,
	percent int32,
	maxDiscount int64,
) Benefit {
	return Benefit{
		title:       title,
		kind:        kind,
		amount:      amount,
		threshold:   threshold,
		percent:     percent,
		maxDiscount: maxDiscount,
	}
}

// DiscountFor 按规格算出优惠额。
//
// 只管金额算得对不对（门槛归 MeetsThreshold），不看券的状态与有效期；
// 折扣额取 min(按比例, maxDiscount)，且不超过订单金额本身。
func (b Benefit) DiscountFor(orderAmount int64) int64 {
	var discount int64
	switch b.kind {
	case KindCash, KindThreshold:
		discount = b.amount
	case KindDiscount:
		discount = orderAmount * int64(b.percent) / 100
		if b.maxDiscount > 0 && discount > b.maxDiscount {
			discount = b.maxDiscount
		}
	}
	if discount > orderAmount {
		discount = orderAmount
	}
	return discount
}

// MeetsThreshold 订单金额是否达到用券门槛。
func (b Benefit) MeetsThreshold(orderAmount int64) bool {
	return orderAmount >= b.threshold
}

// --- 属性只读访问器 (Getters) ---

func (b Benefit) Title() string      { return b.title }
func (b Benefit) Kind() Kind         { return b.kind }
func (b Benefit) Amount() int64      { return b.amount }
func (b Benefit) Threshold() int64   { return b.threshold }
func (b Benefit) Percent() int32     { return b.percent }
func (b Benefit) MaxDiscount() int64 { return b.maxDiscount }
