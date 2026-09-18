package product

// Status 商品状态。
type Status int

const (
	StatusDraft    Status = iota + 1 // 未上架：刚建好，不能被下单
	StatusOnSale                     // 在售
	StatusOffShelf                   // 已下架：不再接新单
)
