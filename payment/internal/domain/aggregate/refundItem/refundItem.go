package refundItem

import (
	"payment/pkg/uuid"
)

// --- 聚合根 (Aggregate Root) ---

// RefundItem 是「退款单里的一行」的聚合根：针对某个订单行项退几件、退多少钱。
//
// 它只记自己这一行的数量与金额。两条跨聚合的判定不在这里做：
//   - 「不能超过已购 / 已退数量」→ 由 OrderItem.RefundQuantity 兜底；
//   - 「整单退款总额不能超过实付」→ 由 Order 侧汇总判断。
//
// 金额一律用「分」的 int64。
type RefundItem struct {
	id          int64
	refundID    int64
	orderItemID int64
	quantity    int32 // 退几件
	amount      int64 // 退多少钱（分）
}

// NewRefundItem 发起退款时建行项。
func NewRefundItem(refundID int64, orderItemID int64, quantity int32, amount int64) (*RefundItem, error) {
	if refundID <= 0 {
		return nil, ErrRefundRequired
	}
	if orderItemID <= 0 {
		return nil, ErrOrderItemRequired
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	return &RefundItem{
		id:          uuid.NewInt64(),
		refundID:    refundID,
		orderItemID: orderItemID,
		quantity:    quantity,
		amount:      amount,
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(
	id int64,
	refundID int64,
	orderItemID int64,
	quantity int32,
	amount int64,
) *RefundItem {
	return &RefundItem{
		id:          id,
		refundID:    refundID,
		orderItemID: orderItemID,
		quantity:    quantity,
		amount:      amount,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// Adjust 审批时调整这一行的退款（部分同意 / 只退件数的一部分）。
//
// 数量与金额一起调 —— 件数变了金额不变是算错账，所以不允许单改一个。
func (r *RefundItem) Adjust(quantity int32, amount int64) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	r.quantity = quantity
	r.amount = amount
	return nil
}

// BelongsTo 这一行是否属于某张退款单（调用方拼装/校验用）。
func (r *RefundItem) BelongsTo(refundID int64) bool { return r.refundID == refundID }

// --- 属性只读访问器 (Getters) ---

func (r *RefundItem) ID() int64          { return r.id }
func (r *RefundItem) RefundID() int64    { return r.refundID }
func (r *RefundItem) OrderItemID() int64 { return r.orderItemID }
func (r *RefundItem) Quantity() int32    { return r.quantity }
func (r *RefundItem) Amount() int64      { return r.amount }
