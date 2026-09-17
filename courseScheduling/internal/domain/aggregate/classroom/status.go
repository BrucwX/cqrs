package classroom

// Status 教室可用状态值对象/枚举
type Status int

const (
	StatusAvailable        Status = iota + 1 // 正常可用
	StatusUnderMaintenance                   // 维修中 / 维护中
	StatusDecommissioned                     // 已弃用 / 报废
)
