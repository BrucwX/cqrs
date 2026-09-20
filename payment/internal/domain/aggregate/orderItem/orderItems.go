package orderItem

// --- 订单行集合 ---

// OrderItems 是一组订单行：把 []OrderItem 起个名字，好在上面挂几个「一次看一批行」的
// **只读**批量视图（筛选、汇总）。
//
// ⚠️ 它不是聚合根：每一行自己就是聚合根，要改行项还是得改那一行（值语义，换一份新的）。
// 所以这里只有读，以及**值语义的批量迁移**（`ConfirmAll` 返回一批新行项，旧的集合一个字节不动）——
// 没有就地改数据的口子，免得出「改了一半的集合」这种半成品状态。
type OrderItems []OrderItem

// Drafts 挑出还没归到订单上的草稿行（正在拼装的那些）。
func (items OrderItems) Drafts() OrderItems {
	drafts := make(OrderItems, 0, len(items))
	for _, item := range items {
		if item.IsDraft() {
			drafts = append(drafts, item)
		}
	}
	return drafts
}

// Confirmed 挑出已经归到订单上的正式行。
func (items OrderItems) Confirmed() OrderItems {
	confirmed := make(OrderItems, 0, len(items))
	for _, item := range items {
		if !item.IsDraft() {
			confirmed = append(confirmed, item)
		}
	}
	return confirmed
}

// AllDraft 每一行是不是都还是草稿（一条都没被归到订单上）—— 下单前「这些行还没被用掉」的检查。
//
// ⚠️ 空集合返回 true：「每一项都…」在空集上恒成立。所以「至少要有一行」这种前提得另外判
// （`len(items) == 0`），别拿它当「有东西可下单」。
func (items OrderItems) AllDraft() bool {
	for _, item := range items {
		if !item.IsDraft() {
			return false
		}
	}
	return true
}

// ConfirmAll 把这一批行项都归到 orderID 上 → 返回一批**正式行**（值语义，旧的集合不动）。
//
// 草稿行才认识怎么变成正式行，所以这里按具体类型取；**有一行不是草稿就整体报错**
// （同一批行项不能被两张订单共用）。全部成功才返回，不会交出「一半归了单」的集合。
func (items OrderItems) ConfirmAll(orderID int64) (OrderItems, error) {
	confirmed := make(OrderItems, 0, len(items))
	for _, item := range items {
		draft, ok := item.(DraftOrderItem)
		if !ok {
			return nil, ErrOrderItemAlreadyConfirmed
		}
		next, err := draft.Confirm(orderID)
		if err != nil {
			return nil, err
		}
		confirmed = append(confirmed, next)
	}
	return confirmed, nil
}

// Total 这一批行项的小计合计（各行单价×数量相加）—— 建订单时的原价就是它。
//
// 金额只长在公共元数据上，而接口刻意不暴露它；两种状态都内嵌了元数据，
// 所以这里按具体类型取（跟落库时取状态名一个路子）。
func (items OrderItems) Total() int64 {
	var total int64
	for _, item := range items {
		switch item := item.(type) {
		case DraftOrderItem:
			total += item.Amount()
		case ConfirmedOrderItem:
			total += item.Amount()
		}
	}
	return total
}

// CustomerIDs 这一批行**各自**给谁（有几个行项就有几个，顺序照原样）。
//
// 拿去批量查客户时用。⚠️ 不去重：同一个受赠人被多行指着就会出现多次，
// 要按人去重的话调用方自己收一下。
func (items OrderItems) CustomerIDs() []int64 {
	ids := make([]int64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.CustomerID())
	}
	return ids
}

// QuantityByProduct 按商品汇总数量：同一个商品出现在多行也会累加到一起，
// 下单前的库存预检（CanOrderCreate）就按这份汇总去跟剩余库存比。
//
// 草稿行也算进来 —— 预检就是在「订单还没建、行项还是草稿」的时候做的。
func (items OrderItems) QuantityByProduct() map[int64]int64 {
	need := make(map[int64]int64, len(items))
	for _, item := range items {
		need[item.ProductID()] += int64(item.Quantity())
	}
	return need
}
