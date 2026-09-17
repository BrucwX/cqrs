package qualification

// Status 授课资质状态
type Status int

const (
	StatusActive  Status = iota + 1 // 正常有效
	StatusRevoked                   // 已吊销/被取消授课资格
	StatusExpired                   // 已过期
)

func (s Status) String() string {
	switch s {
	case StatusActive:
		return "ACTIVE"
	case StatusRevoked:
		return "REVOKED"
	case StatusExpired:
		return "EXPIRED"
	default:
		return "UNKNOWN"
	}
}
