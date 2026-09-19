package product

// --- 聚合根接口 ---

// Product 是商品聚合根的接口：公共元数据 + 密封。
//
// 「这件商品怎么卖」由**具体类型**表达，而不是一个 struct 加一堆开关：
//
//	ServiceProduct    服务类：**没有库存这个概念**，只有「一份的价」（会员、开通服务）
//	LimitedProduct    有限类：有库存（**整数**个计价单位），卖完为止
//	UnlimitedProduct  无限类：没有库存上限，但仍按**计价单位**卖（流量按 GB、课按份）
//
// 后两类都有一个**计价单位** unit（件 / 千克 / GB / 次…）—— price 是「每个单位多少分」；
// 服务类没有它，因为它就是「一份」多少钱。
//
// 于是「服务类商品哪来的库存」这类问题在**编译期**就不存在：remaining / sold 只是
// LimitedProduct 自己的字段，另两类连字段都没有，更不会有出库 / 退单 / 进货这些动作。
//
// 没有 Stock 值对象 —— 有库存的那一类自己拿自己的两个数，规则写在方法里。
//
// ⚠️ 这个包里的两处「类型」不是一回事，别看混了：
//   - `Type`（值对象）= 商品**类别 / 分类**（course / digital；券的适用范围按它的 code 引用）
//   - 本接口的具体类型 = 商品的**种类**（有没有库存、有没有计价单位）
//
// ⚠️ 生命周期状态（Draft / OnSale / OffShelf）**不是**区分维度：它在 ProductMeta 里，
// 三种商品都能上架下架 —— 客户明确说过「不是按照商品 status 来区分」。
//
// 卖出去之后的数量与单价由 orderItem 存快照，商品改名 / 调价 / 下架都不影响已下的单。
type Product interface {
	ID() int64         // 标识符放在第一位
	Name() string      // 展示名
	Price() int64      // 标价（分）；单位含义随种类而定，见各自类型
	ProductType() Type // 商品类别（值对象，不是种类）
	Status() Status
	IsOnSale() bool
	// IsStockLimited 是否有**库存上限**（只有有限类商品有）。
	// 「限量」标签、下单前预检看它；服务类 / 无限类返回 false。
	IsStockLimited() bool
	// RemainingStock 剩余可售数量；没有库存上限（或压根没有库存概念）的商品返回 0，
	// 所以**先问 IsStockLimited 再读它**，别拿 0 当「卖光了」。
	//
	// 读可以统一问，**扣减不能**：ReserveStock / ReleaseStock / RestockStock 只长在
	// 有限类商品上，无限量的东西没有库存可扣。
	RemainingStock() int64
	// HasStockFor 库存够不够卖 need 个 —— 「够不够」是商品自己的规矩，
	// 所以判定收在这里，调用方（下单预检、客服下单）不用自己写比大小的式子。
	// 没有库存上限的商品永远够。
	HasStockFor(need int64) bool
	Meta() ProductMeta // 提取公共元数据快照
	isProductKind()    // 密封约束，防止包外随意实现
}

// 编译期断言：每种商品类型都必须满足 Product —— 漏了内嵌 ProductMeta 在这里就会编译失败。
var (
	_ Product = (*ServiceProduct)(nil)
	_ Product = (*LimitedProduct)(nil)
	_ Product = (*UnlimitedProduct)(nil)
)
