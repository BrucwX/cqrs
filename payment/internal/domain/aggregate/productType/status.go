package productType

// Status 商品分类状态。
type Status int

const (
	StatusEnabled  Status = iota + 1 // 启用：可以挂新商品
	StatusDisabled                   // 停用：不再接新商品，存量商品不受影响
)
