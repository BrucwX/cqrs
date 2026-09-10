package absence

// AbsenceType 缺勤分类
type AbsenceType int

const (
	TypePersonalLeave AbsenceType = iota + 1 // 事假 / 病假（事前申请）
	TypeOfficialDuty                         // 公假 / 因公出差
	TypeUnexcused                            // 旷课 / 未出勤（事后考勤点名）
)

func (t AbsenceType) String() string {
	switch t {
	case TypePersonalLeave:
		return "PERSONAL_LEAVE"
	case TypeOfficialDuty:
		return "OFFICIAL_DUTY"
	case TypeUnexcused:
		return "UNEXCUSED"
	default:
		return "UNKNOWN"
	}
}
