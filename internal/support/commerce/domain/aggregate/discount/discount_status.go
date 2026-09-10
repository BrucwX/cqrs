package discount

// DiscountStatus is the lifecycle state of a discount.
type DiscountStatus int32

const (
	DiscountStatusUnspecified DiscountStatus = 0
	DiscountStatusActive      DiscountStatus = 1 // 生效中
	DiscountStatusExpired     DiscountStatus = 2 // 已过期
	DiscountStatusDisabled    DiscountStatus = 3 // 已禁用
	DiscountStatusDeleted     DiscountStatus = 4
)
