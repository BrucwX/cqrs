package couponType

// Status 券种类的上架状态。
type Status int

const (
	StatusDraft    Status = iota + 1 // 未上架：刚建好，不允许发券
	StatusActive                     // 已上架：可以发券
	StatusInactive                   // 已下架：不再发新券，已发的券不受影响
)
