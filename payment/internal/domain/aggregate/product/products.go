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
