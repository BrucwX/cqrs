package orderItem

// --- 正式行（属于某张订单） ---

// ConfirmedOrderItem 正式行：已经归到订单上的一行，从此有了 orderID。
//
// 它是行项的终态 —— 行项的形态只会在「拼装阶段」变（草稿 → 正式），
// 之后再没有别的状态：退了几件不在这里管（退款流程自己核），
// 所以谁都可以安全地拿 OrderID() 去组装订单。
type ConfirmedOrderItem struct {
	OrderItemMeta
	orderID int64 // 状态专属字段：挂在哪张订单上（草稿行没有这一栏）
}

// ReconstituteConfirmedOrderItem 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteConfirmedOrderItem(meta OrderItemMeta, orderID int64) ConfirmedOrderItem {
	return ConfirmedOrderItem{OrderItemMeta: meta, orderID: orderID}
}

// OrderID 属于哪张订单。
func (i ConfirmedOrderItem) OrderID() int64 { return i.orderID }

// BelongsTo 这一行是否属于某张订单（调用方拼装订单时校验用）。
func (i ConfirmedOrderItem) BelongsTo(orderID int64) bool { return i.orderID == orderID }

// ChangeQuantity 改数量 → 返回更新后的正式行。
//
// ⚠️ 正式行改数量意味着**库存也要跟着调**，那是上层（app / domain service）的事：
// 行项只认钱与数量，不认识库存。
func (i ConfirmedOrderItem) ChangeQuantity(quantity int32) (ConfirmedOrderItem, error) {
	meta, err := i.OrderItemMeta.ChangeQuantity(quantity)
	if err != nil {
		return i, err
	}
	i.OrderItemMeta = meta
	return i, nil
}
