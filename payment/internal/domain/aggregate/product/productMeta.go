package product

import (
	"strings"

	"payment/pkg/uuid"
)

// --- 公共元数据（所有商品类型共享的部分） ---

// ProductMeta 是所有商品都有的那部分：标识、名字、描述、标价、生命周期状态。
//
// 它是**值**。具体商品类型按值内嵌它，于是天然获得 ID()/Name()/…、下面这套
// 「所有商品都一样」的行为（上下架、改名、改描述）以及密封凭据 isProductKind()
// —— 不用每种类型各写一遍；也可以用它的方法做实现的一部分。
//
// price 的**单位含义由具体类型决定**：按单位卖的是「每单位多少分」，一份一份卖的是
// 「一份的价」。金额一律「分」的 int64。
type ProductMeta struct {
	id          int64
	name        string
	description string
	price       int64 // 标价（分）：含义随具体类型而定
	status      Status
}

// validateCommon 建商品时共同的校验。
func validateCommon(name string, price int64) error {
	if strings.TrimSpace(name) == "" {
		return ErrProductNameRequired
	}
	if price <= 0 {
		return ErrInvalidPrice
	}
	return nil
}

// newMeta 建商品时用：默认「未上架」，要 PutOnSale 之后才能被下单。
func newMeta(name string, description string, price int64) ProductMeta {
	return ProductMeta{
		id:          uuid.NewInt64(),
		name:        name,
		description: description,
		price:       price,
		status:      StatusDraft,
	}
}

// ReconstituteMeta 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteMeta(
	id int64,
	name string,
	description string,
	price int64,
	status Status,
) ProductMeta {
	return ProductMeta{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		status:      status,
	}
}

// --- 所有商品都一样的行为 ---

// PutOnSale 上架：只有上架的商品能被下单。
func (m *ProductMeta) PutOnSale() error {
	if m.status == StatusOnSale {
		return ErrProductAlreadyOnSale
	}
	m.status = StatusOnSale
	return nil
}

// TakeOffShelf 下架：之后不再接新单，已经卖出去的不受影响。
func (m *ProductMeta) TakeOffShelf() error {
	if m.status == StatusOffShelf {
		return ErrProductAlreadyOffShelf
	}
	m.status = StatusOffShelf
	return nil
}

// ChangePrice 调价（分）：只影响之后下的单，已下的单用快照价。
//
// ⚠️ price 的单位含义看具体类型：按单位卖的是「每单位多少分」，一份一份卖的是「一份多少分」。
func (m *ProductMeta) ChangePrice(price int64) error {
	if price <= 0 {
		return ErrInvalidPrice
	}
	m.price = price
	return nil
}

// Rename 改名：只影响之后的展示，订单行项持有的是下单时的名字快照。
func (m *ProductMeta) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrProductNameRequired
	}
	m.name = name
	return nil
}

// UpdateDescription 改描述（可以为空，所以不做校验）。
func (m *ProductMeta) UpdateDescription(description string) {
	m.description = description
}

// --- 属性只读访问器 (Getters) ---

func (m ProductMeta) ID() int64           { return m.id }
func (m ProductMeta) Name() string        { return m.name }
func (m ProductMeta) Description() string { return m.description }
func (m ProductMeta) Price() int64        { return m.price }
func (m ProductMeta) Status() Status      { return m.status }
func (m ProductMeta) Meta() ProductMeta   { return m }

// IsOnSale 是否在售（可以被下单）。
func (m ProductMeta) IsOnSale() bool { return m.status == StatusOnSale }

// 私有密封标记：直接挂在 ProductMeta 上，只要内嵌了它就天然获得了密封凭据
func (*ProductMeta) isProductKind() {}
