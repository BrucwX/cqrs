package couponType

import "errors"

var (
	// --- 权益规格 (Benefit) 校验：建券种类时用到 ---

	// ErrInvalidKind 券类型不合法。
	ErrInvalidKind = errors.New("coupon kind is invalid")
	// ErrInvalidAmount 面额不合法：立减/满减必须大于 0，折扣券必须为 0。
	ErrInvalidAmount = errors.New("coupon amount is invalid for its kind")
	// ErrInvalidThreshold 门槛不合法：满减券门槛必须大于 0 且不低于面额，无门槛券必须为 0。
	ErrInvalidThreshold = errors.New("coupon threshold is invalid for its kind")
	// ErrInvalidPercent 折扣百分比不合法（应在 1..99）。
	ErrInvalidPercent = errors.New("coupon percent must be between 1 and 99")
	// ErrBenefitTitleRequired 券名不能为空。
	ErrBenefitTitleRequired = errors.New("coupon benefit title is required")

	// --- 券种类 (CouponType) 校验 ---

	// ErrCouponTypeNameRequired 券种类必须有名字。
	ErrCouponTypeNameRequired = errors.New("coupon type name is required")
	// ErrCouponTypeBenefitRequired 券种类必须有合法的权益规格。
	ErrCouponTypeBenefitRequired = errors.New("coupon type benefit is required")
	// ErrInvalidValidityRule 有效期规则不合法（固定窗口要 from < to，相对窗口要 days > 0）。
	ErrInvalidValidityRule = errors.New("coupon type validity rule is invalid")
	// ErrInvalidQuota 发放上限不能为负（0 表示不限量）。
	ErrInvalidQuota = errors.New("coupon type quota must not be negative")
	// ErrInvalidScope 适用商品类别不合法：限定范围至少要指定一个商品类别。
	ErrInvalidScope = errors.New("coupon type scope is invalid")

	// --- 上架 / 发券 ---

	// ErrCouponTypeAlreadyActive 券种类已经上架了。
	ErrCouponTypeAlreadyActive = errors.New("coupon type is already active")
	// ErrCouponTypeNotActive 券种类没上架，不能下架。
	ErrCouponTypeNotActive = errors.New("coupon type is not active")
	// ErrCouponTypeNotIssuable 未上架的券种类不能发券。
	ErrCouponTypeNotIssuable = errors.New("coupon type is not issuable in its current status")
	// ErrCouponTypeQuotaExhausted 券种类已发完。
	ErrCouponTypeQuotaExhausted = errors.New("coupon type quota is exhausted")
	// ErrCouponTypeExpired 券种类的有效期已经过去，不能再发。
	ErrCouponTypeExpired = errors.New("coupon type has expired")

	// --- 仓库层用 ---

	// ErrCouponTypeRequired 传入的券种类为空。
	ErrCouponTypeRequired = errors.New("coupon type is required")
	// ErrCouponTypeNotFound 指定的券种类不存在。
	ErrCouponTypeNotFound = errors.New("coupon type not found")
)
