package refund

// Status 退款单状态。
type Status int

const (
	StatusApplied  Status = iota + 1 // 已申请，待审批
	StatusApproved                   // 审批通过，待打款
	StatusRejected                   // 审批拒绝
	StatusRefunded                   // 已退款（终态）
	StatusCanceled                   // 买家已撤销
)
