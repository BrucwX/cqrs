package refund

import "time"

// --- 基础元数据（所有状态共享的部分） ---

// RefundMeta 是退款单的公共元数据：不论待审批还是已打款，这几项都一样。
//
// 它是**值**（不是指针）。各状态结构体直接内嵌它，于是天然获得 ID()/OrderID()/CustomerID()/Meta()
// 以及密封凭据 isRefundState() —— 不用每个状态各写一遍。
//
// 金额一律用「分」的 int64。
type RefundMeta struct {
	id         int64
	orderID    int64  // 退的是哪张订单
	customerID int64  // 申请人（必须与订单客户一致，跨聚合判定交给上层）
	amount     int64  // 申请退款金额（分）
	reason     string // 申请原因
	appliedAt  time.Time
}

// ReconstituteMeta 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteMeta(
	id int64,
	orderID int64,
	customerID int64,
	amount int64,
	reason string,
	appliedAt time.Time,
) RefundMeta {
	return RefundMeta{
		id:         id,
		orderID:    orderID,
		customerID: customerID,
		amount:     amount,
		reason:     reason,
		appliedAt:  appliedAt,
	}
}

// 通过 Getter 暴露公共元数据
func (m RefundMeta) ID() int64            { return m.id }
func (m RefundMeta) OrderID() int64       { return m.orderID }
func (m RefundMeta) CustomerID() int64    { return m.customerID }
func (m RefundMeta) Amount() int64        { return m.amount }
func (m RefundMeta) Reason() string       { return m.reason }
func (m RefundMeta) AppliedAt() time.Time { return m.appliedAt }
func (m RefundMeta) Meta() RefundMeta     { return m }

// 私有密封标记：直接挂在 RefundMeta 上，只要内嵌了 RefundMeta 就天然获得了密封凭据
func (RefundMeta) isRefundState() {}
