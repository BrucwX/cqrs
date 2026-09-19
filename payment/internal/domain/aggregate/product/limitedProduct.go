package product

// --- 有限类商品（有库存上限、按计价单位卖） ---

// LimitedProduct 有限类商品：库存有限、卖完为止的东西（教材、实物、限量课时包）。
//
// 它比服务类商品多两样东西：
//   - **计价单位** unit：price 是「每个单位多少分」，单位写在这里（件 / 千克 / 米 / 次…），
//     不然「3 块钱」到底是 3 块一件还是 3 块一千克，只能靠注释猜
//   - **库存**：remaining（还剩多少个单位）+ sold（卖出多少个单位），都是**整数**
//
// 库存是整数，所以数量参数也是 `int64` —— 想买 1.5 个？方法签名根本收不进去，
// 不用在函数体里写「是不是整数」的运行时校验。
type LimitedProduct struct {
	ProductMeta
	unit      string // 计价单位（展示 + 计价用）：件 / 千克 / 米 / 次
	remaining int64  // 剩余可售数量（多少个单位）
	sold      int64  // 累计卖出数量（退单会减回去）
}

// NewLimitedProduct 建有限类商品：price 是**每个计价单位**的价（分）；
// unit 必填；remaining 是最初可售数量（>= 0）。
func NewLimitedProduct(
	name string,
	description string,
	price int64,
	productType Type,
	unit string,
	remaining int64,
) (*LimitedProduct, error) {
	if err := validateCommon(name, price, productType); err != nil {
		return nil, err
	}
	if err := validateUnit(unit); err != nil {
		return nil, err
	}
	if remaining < 0 {
		return nil, ErrInvalidStock
	}
	return &LimitedProduct{
		ProductMeta: newMeta(name, description, price, productType),
		unit:        unit,
		remaining:   remaining,
	}, nil
}

// ReconstituteLimitedProduct 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteLimitedProduct(
	meta ProductMeta,
	unit string,
	remaining int64,
	sold int64,
) *LimitedProduct {
	return &LimitedProduct{ProductMeta: meta, unit: unit, remaining: remaining, sold: sold}
}

// Unit 计价单位（件 / 千克 / 米 / 次…）。
func (p *LimitedProduct) Unit() string { return p.unit }

// RemainingStock 剩余可售数量（多少个单位）—— 覆写元数据里的默认值 0。
func (p *LimitedProduct) RemainingStock() int64 { return p.remaining }

// IsStockLimited 覆写元数据里的默认值（false）：有限类商品就是有库存上限的那种。
func (p *LimitedProduct) IsStockLimited() bool { return true }

// HasStockFor 覆写元数据里的默认值（true）：剩余数量够得着才卖。
//
// ⚠️ 这是**预检**（读一眼就回答）：真要卖就得 ReserveStock，它才挡得住并发 ——
// 两者之间有任何时间差，都可能有别人把最后一件买走。
func (p *LimitedProduct) HasStockFor(need int64) bool { return need <= p.remaining }

// Sold 累计卖出数量（多少个单位）。
func (p *LimitedProduct) Sold() int64 { return p.sold }

// IsExhausted 是否已售罄。
func (p *LimitedProduct) IsExhausted() bool { return p.remaining <= 0 }

// ReserveStock 出库 quantity 个单位（下单时调用）：剩余减、销售量加。
//
// 剩余不够直接拒绝 —— 有限类商品就得有卖光的时候。
func (p *LimitedProduct) ReserveStock(quantity int64) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if quantity > p.remaining {
		return ErrInsufficientStock
	}
	p.remaining -= quantity
	p.sold += quantity
	return nil
}

// ReleaseStock 退单 quantity 个单位（订单取消 / 退款时调用）：剩余加回来、销售量减回去。
//
// 退的数量不能超过累计销售量（不然销量就成负数了）。
func (p *LimitedProduct) ReleaseStock(quantity int64) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if quantity > p.sold {
		return ErrInvalidStock
	}
	p.remaining += quantity
	p.sold -= quantity
	return nil
}

// RestockStock 进货 quantity 个单位（补货）：只加剩余，不动销售量。
func (p *LimitedProduct) RestockStock(quantity int64) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	p.remaining += quantity
	return nil
}
