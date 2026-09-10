package student

// Status 学员账号状态
type Status int

const (
	StatusActive    Status = iota + 1 // 正常可参训
	StatusBlocked                     // 封禁 / 限制参训
	StatusCancelled                   // 已注销 / 离职归档
)
