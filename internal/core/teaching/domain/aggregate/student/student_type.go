package student

// StudentType 学员来源类型
type StudentType int

const (
	TypeExternal StudentType = iota + 1 // 外部客户学员
	TypeInternal                        // 内部员工学员 (自身也可能是讲师)
)

func (t StudentType) String() string {
	switch t {
	case TypeExternal:
		return "EXTERNAL"
	case TypeInternal:
		return "INTERNAL"
	default:
		return "UNKNOWN"
	}
}
