package product

import (
	"strings"

	"payment/pkg/uuid"
)

// --- 公共元数据（四种售卖方式共享的部分） ---

// ProductMeta 是所有商品都有的那部分：标识、名字、描述、类别、标价、生命周期状态。
//
// 它是**值**。各售卖方式按值内嵌它，于是天然获得 ID()/Name()/…、下面这套
// 「所有商品都一样」的行为（上下架、改名、改描述、换类别）以及密封凭据 ——
// 不用每种类型各写一遍。
//
// price 的**单位含义由商品种类决定**：服务类 = 一份的价、有限类 / 无限类 = **每个计价单位**
// （件 / 千克 / GB…，单位存在各自的 unit 字段里）的价。金额一律「分」的 int64。
type ProductMeta struct {
	id          int64
	name        string
	description string
	productType Type  // 商品类别（值对象）
	price       int64 // 标价（分）：含义随售卖方式而定
	status      Status
}

// validateCommon 四种售卖方式建商品时共同的校验。
func validateCommon(name string, price int64, productType Type) error {
	if strings.TrimSpace(name) == "" {
		return ErrProductNameRequired
	}
	if price <= 0 {
		return ErrInvalidPrice
	}
	if !productType.valid() {
		return ErrInvalidProductType
	}
	return nil
}

// validateUnit 有限类 / 无限类商品的「计价单位」校验：必须填（件 / 千克 / GB / 次…）。
//
// 没单位的话，price 是「每个什么」多少分就说不清了 —— 这是这两类商品跟服务类的分别。
func validateUnit(unit string) error {
	if strings.TrimSpace(unit) == "" {
		return ErrProductUnitRequired
	}
	return nil
}

// newMeta 建商品时用：默认「未上架」，要 PutOnSale 之后才能被下单。
func newMeta(name string, description string, price int64, productType Type) ProductMeta {
	return ProductMeta{
		id:          uuid.NewInt64(),
		name:        name,
		description: description,
		productType: productType,
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
	productType Type,
	status Status,
) ProductMeta {
	return ProductMeta{
		id:          id,
		name:        name,
		description: description,
		productType: productType,
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
// price 的单位含义随售卖方式而定（一份 / 每件 / 每千克 / 每单位），调价时别搞错单位。
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

// ChangeType 换商品类别（值对象 Type，不是商品种类 —— 种类是它自己的类型，换不了）。
func (m *ProductMeta) ChangeType(productType Type) error {
	if !productType.valid() {
		return ErrInvalidProductType
	}
	m.productType = productType
	return nil
}

// --- 属性只读访问器 (Getters) ---

func (m ProductMeta) ID() int64           { return m.id }
func (m ProductMeta) Name() string        { return m.name }
func (m ProductMeta) Description() string { return m.description }
func (m ProductMeta) Price() int64        { return m.price }
func (m ProductMeta) ProductType() Type   { return m.productType }
func (m ProductMeta) Status() Status      { return m.status }
func (m ProductMeta) Meta() ProductMeta   { return m }

// IsOnSale 是否在售（可以被下单）。
func (m ProductMeta) IsOnSale() bool { return m.status == StatusOnSale }

// IsStockLimited 默认**没有**库存上限：只有 LimitedProduct 覆写成 true。
func (ProductMeta) IsStockLimited() bool { return false }

// RemainingStock 默认没有剩余可售数量（0）：只有 LimitedProduct 覆写成真值。
func (ProductMeta) RemainingStock() int64 { return 0 }

// HasStockFor 默认**永远够**（没有库存上限就没有卖不出去的一说）：只有 LimitedProduct 覆写。
func (ProductMeta) HasStockFor(int64) bool { return true }

// 私有密封标记：直接挂在 ProductMeta 上，只要内嵌了 ProductMeta 就天然获得了密封凭据
func (*ProductMeta) isProductKind() {}
