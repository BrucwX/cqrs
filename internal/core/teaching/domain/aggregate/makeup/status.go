package makeup

// Status 补课履约状态
type Status int

const (
	StatusPending   Status = iota + 1 // 待确认 (学员申请插班/补课)
	StatusApproved                    // 已准许 (准许进入目标场次补课，占目标课位)
	StatusCompleted                   // 已出勤核销 (补课现场核销成功，冲销缺勤)
	StatusRejected                    // 已驳回 (目标班级满员或时间冲突)
	StatusCancelled                   // 学员主动取消
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusApproved:
		return "APPROVED"
	case StatusCompleted:
		return "COMPLETED"
	case StatusRejected:
		return "REJECTED"
	case StatusCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}
