package refund

// --- 聚合根接口 ---

// Refund 是退款单聚合根的接口：公共元数据 + 密封。
//
// 「退款单处于什么状态」由**具体类型**表达（AppliedRefund / ApprovedRefund / …），
// 不再有状态枚举。于是「没审批的退款单不能打款」这类规则在**编译期**就被拦住了：
// 只有 ApprovedRefund 有 MarkPaid，而 RejectedRefund / PaidRefund / CanceledRefund
// 是终态，什么动作都没有，各自只暴露本状态允许的动作。
//
// 状态迁移一律**返回新状态**（值语义，不改旧值）：
//
//	applied, err := refund.NewRefund(orderID, customerID, amount, reason, now)
//	approved, err := applied.Approve(now)
//	paid := approved.MarkPaid(now)
//	rejected, err := applied.Reject("不符合退款条件", now)
//
// 落库时按具体类型取状态名（type switch），不需要额外的枚举。
//
// 与订单聚合的分工：订单那边只管到「钱收到了」（PaidOrder 即终态），
// 退多少钱、退几件、审批与打款全部落在本聚合（配合 refundItem 记录退的是哪几行）。
type Refund interface {
	ID() int64         // 标识符放在第一位
	OrderID() int64    // 常用聚合根归属标识：退的是哪张订单
	CustomerID() int64 // 申请人
	Meta() RefundMeta
	isRefundState() // 密封约束，防止包外随意实现
}

// 编译期断言：每个状态都必须满足 Refund —— 漏了内嵌 RefundMeta 在这里就会编译失败。
var (
	_ Refund = AppliedRefund{}
	_ Refund = ApprovedRefund{}
	_ Refund = RejectedRefund{}
	_ Refund = PaidRefund{}
	_ Refund = CanceledRefund{}
)
