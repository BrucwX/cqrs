package product

// --- 无限类商品（没有库存上限、按计价单位卖） ---

// UnlimitedProduct 无限类商品：想卖多少卖多少，但仍然是**按计价单位**卖的东西
// （流量按 GB、线上课程按份、包月服务按次）。
//
// 它与服务类商品的差别只有一处，但这一处是要看数据才分得清的：
// 无限类有**计价单位** unit（price 是「每个单位多少分」），服务类只有「一份的价」。
// 它与有限类商品的差别是：**没有 remaining / sold**，所以「卖光了」不存在，
// 也不会有人能写出「先看看还剩多少」这种代码。
type UnlimitedProduct struct {
	ProductMeta
	unit string // 计价单位（展示 + 计价用）：GB / 份 / 次
}

// NewUnlimitedProduct 建无限类商品：price 是**每个计价单位**的价（分）；unit 必填。
func NewUnlimitedProduct(
	name string,
	description string,
	price int64,
	productType Type,
	unit string,
) (*UnlimitedProduct, error) {
	if err := validateCommon(name, price, productType); err != nil {
		return nil, err
	}
	if err := validateUnit(unit); err != nil {
		return nil, err
	}
	return &UnlimitedProduct{
		ProductMeta: newMeta(name, description, price, productType),
		unit:        unit,
	}, nil
}

// ReconstituteUnlimitedProduct 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteUnlimitedProduct(meta ProductMeta, unit string) *UnlimitedProduct {
	return &UnlimitedProduct{ProductMeta: meta, unit: unit}
}

// Unit 计价单位（GB / 份 / 次…）。
func (p *UnlimitedProduct) Unit() string { return p.unit }
