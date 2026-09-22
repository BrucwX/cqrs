package ports

import "context"

// ClosePayQueue 待关单信号队列：多个人往里塞（下单、主动取消…），一个消费者到点取出来办掉。
//
// 入队、取件、确认是**同一份存储**的三件事，所以放一个接口里 —— 适配器一个类型全包
// （比如 Redis：Notify 写 ZSET、score = 到点时刻；Due 取 score <= now 的；Ack 删掉）。
//
// 「生产者 / 消费者」只是两个**角色**，不是两个接口：下单那边只用得到 Notify，
// 轮询那边只用得到 Due / Ack。真嫌它宽，就在用的人那边再声明一个窄接口（Go 的老规矩：
// 接口归调用方），别在这里拆成两个共享接口。
type ClosePayQueue interface {
	// Notify 入队一条信号：这张单到点该看看要不要关（生产者角色）。
	Notify(ctx context.Context, orderID int64) error
	// Due 取一批到点的信号（消费端轮询用；可能为空）。
	Due(ctx context.Context, limit int) ([]int64, error)
	// Ack 确认这条信号处理完了。没 Ack 的之后会被重新取到（至少一次投递），
	// 所以消费动作（CloseIfUnpaid）必须幂等。
	Ack(ctx context.Context, orderID int64) error
}

// ClosePayConsumer 消费端的行为约定：轮询 → 办事。
//
// ⚠️ 它**不是队列能力**，是 app 层（autoClosePay）实现的东西，队列适配器不认识它。
// 分开是为了别让适配器背业务：适配器只管信号怎么存、怎么取；「这张单该不该关、关了之后
// 还欠谁什么」是用例的事。
//
// 只有两个方法：驱动它的轮询（Start）+ 真正办事的那个动作（CloseIfUnpaid）。
// 不再单列一个「消费一条信号」——那本来就是「调 CloseIfUnpaid 办事」，多一层转发
// 只会让人猜它俩到底谁调谁。
type ClosePayConsumer interface {
	// Start 定时轮询：取到点的信号，逐条交给 CloseIfUnpaid，办成了才算消费掉。
	Start(ctx context.Context) error
	// CloseIfUnpaid 具体的消费行为函数：把一张过了支付窗口还没付的订单关掉（幂等）。
	CloseIfUnpaid(ctx context.Context, orderID int64) error
}
