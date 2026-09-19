package order

import "time"

// --- 基础元数据（所有状态共享的部分） ---

// OrderMeta 是订单的公共元数据：不论待支付还是已退款，这几项都一样。
//
// 它是**值**（不是指针）。各状态结构体直接内嵌它，于是天然获得 ID()/CustomerID()/Meta()
// 以及密封凭据 isOrderState() —— 不用每个状态各写一遍。
//
// 金额一律用「分」的 int64。
type OrderMeta struct {
	id             int64
	customerID     int64
	originalAmount int64 // 原价（分）
	discountAmount int64 // 累计优惠金额（分）
	payableAmount  int64 // 应付 = 原价 - 优惠
	placedAt       time.Time
}

// ReconstituteMeta 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteMeta(
	id int64,
	customerID int64,
	originalAmount int64,
	discountAmount int64,
	payableAmount int64,
	placedAt time.Time,
) OrderMeta {
	return OrderMeta{
		id:             id,
		customerID:     customerID,
		originalAmount: originalAmount,
		discountAmount: discountAmount,
		payableAmount:  payableAmount,
		placedAt:       placedAt,
	}
}

// 通过 Getter 暴露公共元数据
func (m OrderMeta) ID() int64             { return m.id }
func (m OrderMeta) CustomerID() int64     { return m.customerID }
func (m OrderMeta) OriginalAmount() int64 { return m.originalAmount }
func (m OrderMeta) DiscountAmount() int64 { return m.discountAmount }
func (m OrderMeta) PayableAmount() int64  { return m.payableAmount }
func (m OrderMeta) PlacedAt() time.Time   { return m.placedAt }

// IsPending 是否处于待支付状态：元数据自己不知道状态，所以默认 false，
// 由 PendingOrder 覆写成 true（外层方法遮蔽被提升的方法，跟状态类型的分工一致）。
func (OrderMeta) IsPending() bool { return false }

// 私有密封标记：直接挂在 OrderMeta 上，只要内嵌了 OrderMeta 就天然获得了密封凭据
func (OrderMeta) isOrderState() {}
