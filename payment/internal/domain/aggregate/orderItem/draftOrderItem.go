package orderItem

// --- 草稿行（还没归到订单上） ---

// DraftOrderItem 草稿行：正在拼装的一行，**不属于任何订单**（下单前的购物车式行项）。
//
// 它没有 orderID 字段，也没有 OrderID() 方法 —— 「这行是哪张订单的」这种问题在草稿行上
// 压根问不出来。要挂到订单上，只能 Confirm 换成 ConfirmedOrderItem。
type DraftOrderItem struct {
	OrderItemMeta
}

// NewDraftOrderItem 拼一行草稿：customerID 是这一行**给谁**（自己买给自己就传自己的客户 ID）；
// 商品名与单价由调用方从商品上取好传进来（快照）。
func NewDraftOrderItem(
	customerID int64,
	productID int64,
	productName string,
	unitPrice int64,
	quantity int32,
) (DraftOrderItem, error) {
	if err := validateCommon(customerID, productID, productName, unitPrice, quantity); err != nil {
		return DraftOrderItem{}, err
	}
	return DraftOrderItem{
		OrderItemMeta: newMeta(customerID, productID, productName, unitPrice, quantity),
	}, nil
}

// ReconstituteDraftOrderItem 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteDraftOrderItem(meta OrderItemMeta) DraftOrderItem {
	return DraftOrderItem{OrderItemMeta: meta}
}

// IsDraft 覆写元数据里的默认值（false）：还没归到订单上就是草稿。
func (i DraftOrderItem) IsDraft() bool { return true }

// ChangeQuantity 改数量 → 返回更新后的草稿行。
func (i DraftOrderItem) ChangeQuantity(quantity int32) (DraftOrderItem, error) {
	meta, err := i.OrderItemMeta.ChangeQuantity(quantity)
	if err != nil {
		return i, err
	}
	i.OrderItemMeta = meta
	return i, nil
}

// Confirm 把草稿行归到某张订单上 → ConfirmedOrderItem（终态，之后不再变形）。
func (i DraftOrderItem) Confirm(orderID int64) (ConfirmedOrderItem, error) {
	if orderID <= 0 {
		return ConfirmedOrderItem{}, ErrOrderRequired
	}
	return ConfirmedOrderItem{
		OrderItemMeta: i.OrderItemMeta,
		orderID:       orderID,
	}, nil
}
