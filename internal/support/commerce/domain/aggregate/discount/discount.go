package discount

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// DiscountType is the type of discount.
type DiscountType int32

const (
	DiscountTypeUnspecified DiscountType = 0
	DiscountTypePercentage  DiscountType = 1 // 百分比折扣，如 8 折
	DiscountTypeFixed       DiscountType = 2 // 固定金额，如减 50 元
	DiscountTypeFreeTrial   DiscountType = 3 // 免费试用
)

// DiscountStatus is the lifecycle state of a discount.
type DiscountStatus int32

const (
	DiscountStatusUnspecified DiscountStatus = 0
	DiscountStatusActive      DiscountStatus = 1 // 生效中
	DiscountStatusExpired     DiscountStatus = 2 // 已过期
	DiscountStatusDisabled    DiscountStatus = 3 // 已禁用
	DiscountStatusDeleted     DiscountStatus = 4
)

// Discount applies to specific business contexts.
type DiscountScope int32

const (
	DiscountScopeUnspecified DiscountScope = 0
	DiscountScopeTeaching    DiscountScope = 1 // 教学课程
	DiscountScopeVenue       DiscountScope = 2 // 场地
	DiscountScopeProduct     DiscountScope = 3 // 商品
	DiscountScopeAll         DiscountScope = 4 // 全部
)

// Discount is a discount/promotion domain object.
type Discount struct {
	ID          string
	Name        string        // 促销名称，如 "新用户优惠"
	Code        string        // 优惠码，如 "WELCOME20"
	Type        DiscountType  // 折扣类型
	Value       int64         // 折扣值：百分比(80=8折) 或 固定金额(单位：分)
	MinAmount   int64         // 最低消费金额(单位：分)，0 表示无门槛
	MaxDiscount int64         // 最大优惠金额(单位：分)，0 表示无上限
	Scope       DiscountScope // 适用范围
	StartDate   time.Time     // 生效开始时间
	EndDate     time.Time     // 生效结束时间
	UsageLimit  int32         // 使用次数限制，0 表示不限
	UsedCount   int32         // 已使用次数
	Status      DiscountStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewDiscount creates a new discount.
func NewDiscount(name, code string, typ DiscountType, value int64, scope DiscountScope, startDate, endDate time.Time) *Discount {
	return &Discount{
		ID:        uuid.New().String(),
		Name:      name,
		Code:      code,
		Type:      typ,
		Value:     value,
		Scope:     scope,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    DiscountStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate checks if the discount domain object is valid.
func (d *Discount) Validate() error {
	if strings.TrimSpace(d.Name) == "" {
		return ErrDiscountInvalidArgument
	}
	if strings.TrimSpace(d.Code) == "" {
		return ErrDiscountInvalidArgument
	}
	if d.Type == DiscountTypeUnspecified {
		return ErrDiscountInvalidArgument
	}
	if d.Value <= 0 {
		return ErrDiscountInvalidArgument
	}
	if d.StartDate.After(d.EndDate) {
		return ErrDiscountInvalidArgument
	}
	return nil
}

// CalculateDiscount calculates the discount amount.
func CalculateDiscount(d *Discount, amount int64) int64 {
	var discountAmount int64
	switch d.Type {
	case DiscountTypePercentage:
		discountAmount = amount * (100 - d.Value) / 100
	case DiscountTypeFixed:
		discountAmount = d.Value
	default:
		return 0
	}
	if d.MaxDiscount > 0 && discountAmount > d.MaxDiscount {
		discountAmount = d.MaxDiscount
	}
	if discountAmount > amount {
		discountAmount = amount
	}
	return discountAmount
}
