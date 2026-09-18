package order

// Status 订单状态。
type Status int

const (
	StatusPending  Status = iota + 1 // 待支付
	StatusPaid                       // 已支付
	StatusCanceled                   // 已取消（主动）
	StatusClosed                     // 已关闭（超时未支付）
	StatusRefunded                   // 已退款
)
