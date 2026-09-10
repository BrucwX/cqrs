package absence

// ReviewStatus 请假审批状态
type ReviewStatus int

const (
	StatusPending  ReviewStatus = iota + 1 // 待审核
	StatusApproved                         // 已批准
	StatusRejected                         // 已驳回（转为旷课）
)

func (s ReviewStatus) String() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusApproved:
		return "APPROVED"
	case StatusRejected:
		return "REJECTED"
	default:
		return "UNKNOWN"
	}
}
