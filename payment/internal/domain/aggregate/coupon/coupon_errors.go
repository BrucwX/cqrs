package coupon

import "errors"

var (
	// --- 权益规格 (Benefit) 校验 ---

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
	// ErrBenefitRequired 券必须有合法的权益规格。
	ErrBenefitRequired = errors.New("coupon benefit is required")

	// --- 适用范围 (Scope) / 有效期规则 (Validity) 校验 ---

	// ErrInvalidScope 适用范围不合法：限定范围至少要指定一个商品分类。
	ErrInvalidScope = errors.New("coupon scope is invalid")
	// ErrInvalidValidityRule 有效期规则不合法（固定窗口要 from < to，相对窗口要 days > 0）。
	ErrInvalidValidityRule = errors.New("coupon validity rule is invalid")

	// --- 券实例 (Coupon) 校验 ---

	// ErrCouponOwnerRequired 券必须有持有人。
	ErrCouponOwnerRequired = errors.New("coupon owner is required")
	// ErrCouponOrderRequired 核销必须指明订单。
	ErrCouponOrderRequired = errors.New("coupon order is required")
	// ErrInvalidOrderAmount 用券的订单金额必须大于 0。
	ErrInvalidOrderAmount = errors.New("order amount must be greater than zero")
	// ErrThresholdNotMet 订单金额未达到用券门槛（门槛来自权益快照）。
	ErrThresholdNotMet = errors.New("order amount does not reach the coupon threshold")
	// ErrCouponNotApplicable 券不适用于这类商品（适用范围来自发券时的快照）。
	ErrCouponNotApplicable = errors.New("coupon is not applicable to this product type")

	// --- 券实例状态流转 ---

	// ErrCouponNotStarted 券还没到生效时间。
	ErrCouponNotStarted = errors.New("coupon is not effective yet")
	// ErrCouponExpired 券已过期。
	ErrCouponExpired = errors.New("coupon has expired")
	// ErrCouponAlreadyUsed 券已核销。
	ErrCouponAlreadyUsed = errors.New("coupon has already been used")
	// ErrCouponVoided 券已作废。
	ErrCouponVoided = errors.New("coupon has been voided")
	// ErrCouponNotUsed 只有已核销的券才能退回。
	ErrCouponNotUsed = errors.New("coupon has not been used")
	// ErrCouponNotExpired 还没到有效期终点，不能提前置为过期。
	ErrCouponNotExpired = errors.New("coupon has not expired yet")
	// ErrCouponNotExpirable 已核销 / 已作废的券不参与过期处理。
	ErrCouponNotExpirable = errors.New("coupon can not be expired in its current status")

	// --- 仓库层用 ---

	// ErrCouponRequired 传入的券为空。
	ErrCouponRequired = errors.New("coupon is required")
	// ErrCouponNotFound 指定的券不存在。
	ErrCouponNotFound = errors.New("coupon not found")
)
