package buyer

import "errors"

var (
	// --- 联系方式 (Contact) 校验 ---

	// ErrBuyerNameRequired 买家姓名不能为空。
	ErrBuyerNameRequired = errors.New("buyer name is required")
	// ErrBuyerPhoneRequired 买家手机号不能为空。
	ErrBuyerPhoneRequired = errors.New("buyer phone is required")

	// --- 买家 (Buyer) 校验 ---

	// ErrBuyerContactRequired 买家必须有合法的联系方式。
	ErrBuyerContactRequired = errors.New("buyer contact is required")
	// ErrBuyerAlreadyFrozen 买家已经冻结了。
	ErrBuyerAlreadyFrozen = errors.New("buyer is already frozen")
	// ErrBuyerNotFrozen 买家没被冻结，不能解冻。
	ErrBuyerNotFrozen = errors.New("buyer is not frozen")

	// --- 积分 ---

	// ErrInvalidPoints 积分数量必须为正（AdjustPoints 只要求非 0）。
	ErrInvalidPoints = errors.New("points must be positive")
	// ErrInsufficientPoints 积分余额不足。
	ErrInsufficientPoints = errors.New("insufficient points")

	// --- 仓库层用 ---

	// ErrBuyerRequired 传入的买家为空。
	ErrBuyerRequired = errors.New("buyer is required")
	// ErrBuyerNotFound 指定的买家不存在。
	ErrBuyerNotFound = errors.New("buyer not found")
)
