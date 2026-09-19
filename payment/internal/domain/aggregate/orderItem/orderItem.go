package orderItem

import (
	"strings"

	"payment/pkg/uuid"
)

// --- 聚合根 (Aggregate Root) ---

// OrderItem 是「订单里的一行」的聚合根：买了哪件商品、单价多少、几件。
//
// 商品名与单价是**下单那一刻的快照** —— 商品事后改名、调价、下架，
// 已经下过的单都不受影响。
//
// 它只记「买了什么、买了几件」；退了几件不在这里管 —— 退款数量由退款流程自己核
// （要防超额退款就拿本行的 quantity 去比对）。
//
// 金额一律用「分」的 int64。
type OrderItem struct {
	id          int64
	orderID     int64
	productID   int64
	productName string // 下单时的商品名快照
	unitPrice   int64  // 下单时的单价快照（分）
	quantity    int32  // 购买数量
}

// NewOrderItem 下单时建行项。
func NewOrderItem(
	orderID int64,
	productID int64,
	productName string,
	unitPrice int64,
	quantity int32,
) (*OrderItem, error) {
	if orderID <= 0 {
		return nil, ErrOrderRequired
	}
	if productID <= 0 {
		return nil, ErrProductRequired
	}
	if strings.TrimSpace(productName) == "" {
		return nil, ErrProductNameRequired
	}
	if unitPrice <= 0 {
		return nil, ErrInvalidUnitPrice
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	return &OrderItem{
		id:          uuid.NewInt64(),
		orderID:     orderID,
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(
	id int64,
	orderID int64,
	productID int64,
	productName string,
	unitPrice int64,
	quantity int32,
) *OrderItem {
	return &OrderItem{
		id:          id,
		orderID:     orderID,
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// Amount 小计 = 单价 × 数量（算出来的，不存字段，免得与数量漂移）。
func (i *OrderItem) Amount() int64 { return i.unitPrice * int64(i.quantity) }

// ChangeQuantity 改购买数量（比如下单后改件数）。
func (i *OrderItem) ChangeQuantity(quantity int32) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	i.quantity = quantity
	return nil
}

// BelongsTo 这一行是否属于某张订单（调用方拼装/校验用）。
func (i *OrderItem) BelongsTo(orderID int64) bool { return i.orderID == orderID }

// --- 属性只读访问器 (Getters) ---

func (i *OrderItem) ID() int64           { return i.id }
func (i *OrderItem) OrderID() int64      { return i.orderID }
func (i *OrderItem) ProductID() int64    { return i.productID }
func (i *OrderItem) ProductName() string { return i.productName }
func (i *OrderItem) UnitPrice() int64    { return i.unitPrice }
func (i *OrderItem) Quantity() int32     { return i.quantity }
