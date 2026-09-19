package orderItem

import (
	"strings"

	"payment/pkg/uuid"
)

// --- 基础元数据（草稿行与正式行共享的部分） ---

// OrderItemMeta 是订单行的公共元数据：不论草稿还是正式，这几项都一样。
//
// 它是**值**（不是指针）。两种状态直接内嵌它，于是天然获得 ID()/ProductID()/Quantity()/Meta()
// 以及密封凭据 isOrderItemState() —— 不用每种状态各写一遍。
//
// ⚠️ **这里没有 orderID** —— 草稿行不属于任何订单，orderID 只长在 ConfirmedOrderItem 上。
//
// 商品名与单价是**下单那一刻的快照**：商品事后改名、调价、下架，已经记下的行项都不受影响。
// 金额一律用「分」的 int64。
type OrderItemMeta struct {
	id          int64
	customerID  int64 // 这一行**给谁**（受赠人）：自己买给自己时＝下单的那位客户，送礼时不是
	productID   int64
	productName string // 下单时的商品名快照
	unitPrice   int64  // 下单时的单价快照（分）
	quantity    int32  // 购买数量
}

// validateCommon 两种状态建行项时共同的校验。
func validateCommon(customerID int64, productID int64, productName string, unitPrice int64, quantity int32) error {
	if customerID <= 0 {
		return ErrCustomerRequired
	}
	if productID <= 0 {
		return ErrProductRequired
	}
	if strings.TrimSpace(productName) == "" {
		return ErrProductNameRequired
	}
	if unitPrice <= 0 {
		return ErrInvalidUnitPrice
	}
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	return nil
}

// newMeta 草稿行建行项时用。
func newMeta(
	customerID int64,
	productID int64,
	productName string,
	unitPrice int64,
	quantity int32,
) OrderItemMeta {
	return OrderItemMeta{
		id:          uuid.NewInt64(),
		customerID:  customerID,
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}
}

// ReconstituteMeta 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteMeta(
	id int64,
	customerID int64,
	productID int64,
	productName string,
	unitPrice int64,
	quantity int32,
) OrderItemMeta {
	return OrderItemMeta{
		id:          id,
		customerID:  customerID,
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}
}

// --- 两种状态都一样的行为 ---

// Amount 小计 = 单价 × 数量（算出来的，不存字段，免得与数量漂移）。
func (m OrderItemMeta) Amount() int64 { return m.unitPrice * int64(m.quantity) }

// ChangeQuantity 改数量 → 返回更新后的元数据（值语义，不改旧值）。
//
// 两种状态都能改（草稿随便改；正式行改完要同步库存，那是上层的事）。
func (m OrderItemMeta) ChangeQuantity(quantity int32) (OrderItemMeta, error) {
	if quantity <= 0 {
		return m, ErrInvalidQuantity
	}
	m.quantity = quantity
	return m, nil
}

// --- 属性只读访问器 (Getters) ---

func (m OrderItemMeta) ID() int64           { return m.id }
func (m OrderItemMeta) CustomerID() int64   { return m.customerID }
func (m OrderItemMeta) ProductID() int64    { return m.productID }
func (m OrderItemMeta) ProductName() string { return m.productName }
func (m OrderItemMeta) UnitPrice() int64    { return m.unitPrice }
func (m OrderItemMeta) Quantity() int32     { return m.quantity }
func (m OrderItemMeta) Meta() OrderItemMeta { return m }

// IsDraft 默认**不是**草稿：只有 DraftOrderItem 覆写成 true。
func (OrderItemMeta) IsDraft() bool { return false }

// 私有密封标记：直接挂在 OrderItemMeta 上，只要内嵌了它就天然获得了密封凭据
func (OrderItemMeta) isOrderItemState() {}
