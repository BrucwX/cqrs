package coupon

import "errors"

var (
	// --- 券实例 (Coupon) 校验 ---

	// ErrCouponOwnerRequired 券必须有持有人。
	ErrCouponOwnerRequired = errors.New("coupon owner is required")
	// ErrCouponTypeRequired 券必须能追溯到一种券种类（couponType）。
	ErrCouponTypeRequired = errors.New("coupon source type is required")
	// ErrCouponOrderRequired 核销必须指明订单。
	ErrCouponOrderRequired = errors.New("coupon order is required")
	// ErrInvalidValidity 有效期不合法（开始必须早于结束）。
	ErrInvalidValidity = errors.New("coupon validity window is invalid")
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
