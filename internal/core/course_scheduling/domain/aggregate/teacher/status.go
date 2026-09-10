package teacher

// Status 讲师任职状态枚举
type Status int

const (
	StatusActive     Status = iota + 1 // 在职（正常可排课）
	StatusOnLeave                      // 请假 / 休假中
	StatusTerminated                   // 已离职 / 解约
)

func (s Status) String() string {
	switch s {
	case StatusActive:
		return "ACTIVE"
	case StatusOnLeave:
		return "ON_LEAVE"
	case StatusTerminated:
		return "TERMINATED"
	default:
		return "UNKNOWN"
	}
}
