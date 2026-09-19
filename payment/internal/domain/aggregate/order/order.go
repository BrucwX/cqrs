package order

// --- 聚合根接口 ---

// Order 是订单聚合根的接口：公共元数据 + 密封。
//
// 「订单处于什么状态」由**具体类型**表达（PendingOrder / PaidOrder / …），不再有状态枚举。
// 于是「状态不对就不能做某事」这类规则在**编译期**就被拦住了：
// PaidOrder 没有 MarkPaid（不会再付一次）、CanceledOrder / ClosedOrder 什么动作都没有，
// 每种状态只暴露自己允许的动作。
//
// 状态迁移一律**返回新状态**（值语义，不改旧值）：
//
//	pending, _ := order.NewOrder(customerID, originalAmount, placedAt, ttl)
//	paid, err := pending.MarkPaid(paidAmount, wechatPay, tradeNo, now)
//	canceled := pending.Cancel("不买了", now)
//	closed, err := pending.CloseIfExpired(now)
//
// 落库时按具体类型取状态名（type switch），不需要额外的枚举。
//
// 退款不进这个聚合：「已退款」是 refund 聚合的事，订单状态到 Paid 就停住了。
type Order interface {
	ID() int64         // 标识符放在第一位
	CustomerID() int64 // 常用聚合根归属标识
	// IsPending 订单是否还处于**待支付**状态（能不能收款 / 改优惠看它）。
	// 只想知道「还要不要付款」的调用方（app handler、对账）用它，
	// 不用写 type switch，也不用碰任何状态枚举。
	IsPending() bool
	isOrderState() // 密封约束，防止包外随意实现
}

// 编译期断言：每个状态都必须满足 Order —— 漏了内嵌 OrderMeta 在这里就会编译失败。
var (
	_ Order = PendingOrder{}
	_ Order = PaidOrder{}
	_ Order = CanceledOrder{}
	_ Order = ClosedOrder{}
)
