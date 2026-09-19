package product

// --- 服务类商品（没有库存这个概念） ---

// ServiceProduct 服务类商品：会员、开通某项服务这种**卖多少份都不受影响**的东西。
//
// 它没有 remaining / sold，也没有出库 / 退单 / 进货 —— 「卖光了」在服务类商品上压根不存在，
// 所以不该有地方能写出「先看看还剩多少」这种代码。它也没有计价单位：就是「一份多少钱」。
//
// 要按单位计价（每件 / 每千克 / 每 GB）就用 LimitedProduct 或 UnlimitedProduct。
type ServiceProduct struct {
	ProductMeta // 自动获得 ID()、Name()、Price()、ProductType()、Meta() 与 isProductKind()
}

// NewServiceProduct 建服务类商品：price 是**一份**的价（分）。
func NewServiceProduct(
	name string,
	description string,
	price int64,
	productType Type,
) (*ServiceProduct, error) {
	if err := validateCommon(name, price, productType); err != nil {
		return nil, err
	}
	return &ServiceProduct{
		ProductMeta: newMeta(name, description, price, productType),
	}, nil
}

// ReconstituteServiceProduct 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteServiceProduct(meta ProductMeta) *ServiceProduct {
	return &ServiceProduct{ProductMeta: meta}
}
