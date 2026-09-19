package customer

// Status 客户状态。
type Status int

const (
	StatusActive Status = iota + 1 // 正常：可以下单、退款
	StatusFrozen                   // 已冻结（风控 / 违规）
)
