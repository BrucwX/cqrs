package customer

// --- 客户集合（一批客户） ---

// Customers 是一批客户。
//
// 跟 product.Products / orderItem.OrderItems 一个路子：一起取回来（比如一批行项涉及的
// 受赠人）之后，批量的判定／挑选规则写在集合上，调用方（领域服务）只管取数、不用自己写循环。
type Customers []Customer

// AccountIDs 这一批客户**各自**的登录账号 id（有几个客户就有几个，顺序照原样）。
//
// 拿到账号体系（user 那一侧）批量查东西时用：聚合这边的字段叫 AccountID，
// 就是「这个客户对应哪个登录账号」。
func (cs Customers) AccountIDs() []string {
	ids := make([]string, 0, len(cs))
	for _, c := range cs {
		ids = append(ids, c.AccountID())
	}
	return ids
}
