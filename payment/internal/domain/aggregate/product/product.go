package product

// --- 聚合根接口 ---

// Product 是商品聚合根的接口：公共元数据 + 密封。
//
// 「这件商品怎么卖」由**具体类型**表达，而不是一个 struct 加一堆开关 —— 具体类型
// 由使用方自己定义：内嵌 ProductMeta 就自动拿到 ID()/Name()/… 与密封凭据，
// 库存 / 计价那些只有该类商品才有的东西自己加。
//
// 于是「服务类商品哪来的库存」这类问题在**编译期**就不存在：库存相关的方法只长在
// 真正有库存的那个类型上。
//
// ⚠️ 生命周期状态（Draft / OnSale / OffShelf）**不是**区分维度：它在 ProductMeta 里，
// 各种商品都能上架下架 —— 客户明确说过「不是按照商品 status 来区分」。
//
// 卖出去之后的数量与单价由 orderItem 存快照，商品改名 / 调价 / 下架都不影响已下的单。
type Product interface {
	ID() int64    // 标识符放在第一位
	Name() string // 展示名
	Price() int64 // 标价（分）；单位含义随具体类型而定
	Status() Status
	IsOnSale() bool
	Type() Type       // 分类值对象：课程 / 数码 / 服务包（券的适用范围按它匹配）
	TargetId() string // 具体类型的标识（如课程 / 服务包的 id），由具体类型自己定义
	isProductKind()   // 密封约束，防止包外随意实现（内嵌 ProductMeta 就自动满足）
}
