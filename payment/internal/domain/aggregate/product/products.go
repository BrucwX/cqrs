package product

// --- 商品集合（一批商品） ---

// Products 是一批商品。
//
// 跟 orderItem.OrderItems 一个路子：一批东西凑在一起才成立的判定／统计，规则写在集合上，
// 调用方（领域服务）只管取数，不用自己写循环。
type Products []Product

// AllOfType 这一批商品是不是都属于 t 这一类。
//
// 空集合算「是」—— 没有东西违反。
func (ps Products) AllOfType(t Type) bool {
	for _, p := range ps {
		if p.Type() != t {
			return false
		}
	}
	return true
}

// TargetIds 这一批商品**各自**指向的目标 id（有几个商品就有几个，顺序照原样）。
//
// 拿去到别的上下文查东西时用，比如按这批课程 id 去查课程侧的信息。
func (ps Products) TargetIds() []string {
	ids := make([]string, 0, len(ps))
	for _, p := range ps {
		ids = append(ids, p.TargetId())
	}
	return ids
}

// GroupByType 按**商品类型**把这一批商品分开：key 是分类，value 是这一类的那几件。
//
// 留着的用处：哪天真要多卖几类（不同类型「怎么校验、买完要做什么」不一样），
// 调用方先分组、再一类一类派给对应能力。
func (ps Products) GroupByType() map[Type]Products {
	groups := make(map[Type]Products)
	for _, p := range ps {
		t := p.Type()
		groups[t] = append(groups[t], p)
	}
	return groups
}
