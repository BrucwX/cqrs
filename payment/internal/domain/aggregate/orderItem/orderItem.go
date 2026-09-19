package orderItem

// --- 聚合根接口 ---

// OrderItem 是订单行聚合根的接口：公共元数据 + 密封。
//
// 「这一行处于什么状态」由**具体类型**表达，不再有状态枚举：
//
//	DraftOrderItem      草稿行：正在拼装（像购物车里的那一行），**不属于任何订单**
//	ConfirmedOrderItem  正式行：已经 Confirm 到某张订单上，从此有 orderID
//
// 于是「草稿行属于哪张订单」这个问题在**编译期**就不存在 —— `DraftOrderItem` 没有
// orderID 字段、也没有 OrderID() 方法，不是「orderID 为 0」那种靠约定，而是根本拿不到；
// 想挂到订单上，只能走 Confirm 换成正式行。
//
// 状态迁移一律**返回新状态**（值语义，不改旧值）：
//
//	draft, err := orderItem.NewDraftOrderItem(customerID, productID, productName, unitPrice, quantity)
//	confirmed, err := draft.Confirm(orderID)
//
// 落库时按具体类型取状态名（type switch），不需要额外的枚举。
type OrderItem interface {
	ID() int64         // 标识符放在第一位
	CustomerID() int64 // 这一行**给谁**（受赠人）：自己买给自己＝下单的客户，送礼时不是
	ProductID() int64  // 常用聚合根归属标识：买的是哪件商品
	Quantity() int32   // 下单数量（库存预检要按商品汇总它）
	// IsDraft 这一行是不是还没归到订单上的**草稿行**（正在拼装）。
	// 还没确认的行项不该参与扣库存 / 算订单金额，用它一票闸在前面。
	IsDraft() bool
	isOrderItemState() // 密封约束，防止包外随意实现
}

// 编译期断言：两种状态都必须满足 OrderItem —— 漏了内嵌 OrderItemMeta 在这里就会编译失败。
var (
	_ OrderItem = DraftOrderItem{}
	_ OrderItem = ConfirmedOrderItem{}
)
