package orderItem

import (
	"maps"
	"slices"
)

// --- 订单行集合 ---

// OrderItems 是一组订单行：把 []OrderItem 起个名字，好在上面挂几个「一次看一批行」的
// **只读**批量视图（筛选、汇总）。
//
// ⚠️ 它不是聚合根：每一行自己就是聚合根，要改行项还是得改那一行（值语义，换一份新的）。
// 所以这里只有读，没有任何改数据的口子 —— 免得出现「改了一半的集合」这种半成品状态。
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

// CustomerIDs 这批行涉及到的**受赠人**（去重 + 升序）：同一批东西送给同一个人时只查一次客户。
func (items OrderItems) CustomerIDs() []int64 {
	seen := make(map[int64]struct{}, len(items))
	for _, item := range items {
		seen[item.CustomerID()] = struct{}{}
	}
	return slices.Sorted(maps.Keys(seen))
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
