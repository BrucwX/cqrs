package enrollment

// Status 课程注册履约状态枚举
type Status int

const (
	StatusEnrolled    Status = iota + 1 // 正常报班/在读
	StatusCompleted                     // 已结业/修完
	StatusDropped                       // 已退课/取消
	StatusNotSelected                   // 未选课（已进入名单/待选，尚未报班）
)

func (s Status) String() string {
	switch s {
	case StatusEnrolled:
		return "ENROLLED"
	case StatusCompleted:
		return "COMPLETED"
	case StatusDropped:
		return "DROPPED"
	case StatusNotSelected:
		return "NOT_SELECTED"
	default:
		return "UNKNOWN"
	}
}
