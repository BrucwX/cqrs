package courseSlotChange

// ChangeType 变更类型
type ChangeType int

const (
	TypeReschedule ChangeType = iota + 1 // 改期调课（时间变动）
	TypeSubstitute                       // 临时代课（讲师变动）
	TypeRelocate                         // 临时换教室（地点变动）
	TypeComposite                        // 复合变动（多项同时调整）
)
