// Package autoClosePay 是「超时关单」用例（写侧）：轮询到点后调它，把过了支付窗口
// 还没付的订单关掉 —— 先本地关单，再让渠道侧 / 课程侧跟上。
//
// 它实现 ports.ClosePayConsumer：别人下单时往信号队列里塞一个订单号
// （ports.ClosePayQueue.Notify），这里到点把信号取出来消费掉。
//
// 本文件只放「它是什么、依赖谁」：轮询在 start.go，关单动作在 closeIfUnpaid.go。
package autoClosePay

import (
	"time"

	"payment/internal/domain/repo"
	"payment/internal/ports"
)

// closePayScheduler 超时关单的消费者：轮询信号队列 → 逐条关单。
type closePayScheduler struct {
	orderCommand repo.OrderCommand
	pay          ports.PayService
	// queue 消费端用到的那两下：取到点的信号、处理完确认掉。
	// 入队（Notify）是下单那边的事 —— 同一个队列，两个角色。
	queue         ports.ClosePayQueue
	courseService ports.CourseService
	tx            ports.Transaction
	// 退课要的数据都在行项上（买了哪些课、报名的是哪些账号），
	// 订单聚合只认钱不认货，所以要另外查一次行项 + 商品 + 客户。
	itemCommand     repo.OrderItemCommand
	customerCommand repo.CustomerCommand
	productCommand  repo.ProductCommand

	interval time.Duration
	batch    int
}

// NewClosePayScheduler 构造：依赖由上层（wire）注入。
//
// 返回接口而不是具体类型：外面需要的只是「一个能 Start 的消费者」，具体类型不用露出去
// （不然「导出的构造函数返回未导出类型」这事会被 lint 念）。
func NewClosePayScheduler(
	orderCommand repo.OrderCommand,
	itemCommand repo.OrderItemCommand,
	customerCommand repo.CustomerCommand,
	productCommand repo.ProductCommand,
	pay ports.PayService,
	courseService ports.CourseService,
	queue ports.ClosePayQueue,
	tx ports.Transaction,
) ports.ClosePayConsumer {
	return &closePayScheduler{
		orderCommand:    orderCommand,
		itemCommand:     itemCommand,
		customerCommand: customerCommand,
		productCommand:  productCommand,
		pay:             pay,
		courseService:   courseService,
		queue:           queue,
		tx:              tx,
		interval:        defaultPollInterval,
		batch:           defaultBatch,
	}
}

var _ ports.ClosePayConsumer = (*closePayScheduler)(nil)
