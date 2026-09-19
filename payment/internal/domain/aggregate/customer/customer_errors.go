package customer

import "errors"

var (
	// --- 联系方式 (Contact) 校验 ---

	// ErrCustomerNameRequired 客户姓名不能为空。
	ErrCustomerNameRequired = errors.New("customer name is required")
	// ErrCustomerPhoneRequired 客户手机号不能为空。
	ErrCustomerPhoneRequired = errors.New("customer phone is required")

	// --- 客户 (Customer) 校验 ---

	// ErrCustomerContactRequired 客户必须有合法的联系方式。
	ErrCustomerContactRequired = errors.New("customer contact is required")
	// ErrCustomerAlreadyFrozen 客户已经冻结了。
	ErrCustomerAlreadyFrozen = errors.New("customer is already frozen")
	// ErrCustomerNotFrozen 客户没被冻结，不能解冻。
	ErrCustomerNotFrozen = errors.New("customer is not frozen")

	// --- 积分 ---

	// ErrInvalidPoints 积分数量必须为正（AdjustPoints 只要求非 0）。
	ErrInvalidPoints = errors.New("points must be positive")
	// ErrInsufficientPoints 积分余额不足。
	ErrInsufficientPoints = errors.New("insufficient points")

	// --- 仓库层用 ---

	// ErrCustomerRequired 传入的客户为空。
	ErrCustomerRequired = errors.New("customer is required")
	// ErrCustomerNotFound 指定的客户不存在。
	ErrCustomerNotFound = errors.New("customer not found")
)
