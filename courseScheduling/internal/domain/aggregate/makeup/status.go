package makeup

// Status 补课状态（预约即生效，无需审批）
type Status int

const (
	StatusBooked    Status = iota + 1 // 已预约（待补课）
	StatusCompleted                   // 已补课（现场核销）
	StatusCancelled                   // 已取消
)

func (s Status) String() string {
	switch s {
	case StatusBooked:
		return "BOOKED"
	case StatusCompleted:
		return "COMPLETED"
	case StatusCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}
