package coupon

// Status 优惠券状态。
type Status int

const (
	StatusUnused  Status = iota + 1 // 未使用
	StatusUsed                      // 已核销
	StatusExpired                   // 已过期
	StatusVoided                    // 已作废（人工）
)
