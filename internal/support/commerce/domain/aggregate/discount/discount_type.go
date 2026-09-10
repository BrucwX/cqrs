package discount

// DiscountType is the type of discount.
type DiscountType int32

const (
	DiscountTypeUnspecified DiscountType = 0
	DiscountTypePercentage  DiscountType = 1 // 百分比折扣，如 8 折
	DiscountTypeFixed       DiscountType = 2 // 固定金额，如减 50 元
	DiscountTypeFreeTrial   DiscountType = 3 // 免费试用
)
