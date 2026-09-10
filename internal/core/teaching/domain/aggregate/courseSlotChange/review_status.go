package courseSlotChange

// ReviewStatus 审批状态
type ReviewStatus int

const (
	StatusPending   ReviewStatus = iota + 1 // 待审核
	StatusApproved                          // 已生效
	StatusRejected                          // 已驳回
	StatusWithdrawn                         // 已撤销/撤回
)

func (s ReviewStatus) String() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusApproved:
		return "APPROVED"
	case StatusRejected:
		return "REJECTED"
	case StatusWithdrawn:
		return "WITHDRAWN"
	default:
		return "UNKNOWN"
	}
}
