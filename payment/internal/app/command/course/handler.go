// Package course 是「买课」用例（写侧）。
//
// 一个 Handler 挂这一条业务线上的多个用例（按文件分：createOrder.go 下单、
// closeIfUnpaid.go 超时关单），依赖一次性注入 —— 下单入队、关单消费的是同一个
// 延时队列，所以两边合用一个 Handler，不做「一个用例一个 Handler」。
//
// 本文件只放「Handler 是什么、依赖谁」；具体用例在各自的文件里。
package course

import (
	"payment/internal/domain/repo"
	"payment/internal/ports"
)

// Handler 买课写侧处理器：下单（CreateCOrder）与超时关单（CloseCourseIfUnpaid）。
//
// 依赖：订单 / 行项 / 客户 / 商品的写侧仓储，课程侧判定服务（校验报名、退课），
// 支付渠道端口（发起支付、渠道侧关单），事务端口，等幂服务，以及延时队列。
type Handler struct {
	orderCommand     repo.OrderCommand
	orderItemCommand repo.OrderItemCommand
	customerCommand  repo.CustomerCommand
	productCommand   repo.ProductCommand
	courseService    ports.CourseService
	pay              ports.PayService
	tx               ports.Transaction
	idservice        ports.IdempotencyService
	// queue 延时队列，一个队列两个角色：下单时入队（Notify），
	// 超时关单时等信号、处理完确认掉。
	queue ports.DelayScheduler
}

// NewHandler 构造：依赖由上层（wire）注入。
//
// 构造**只是装配**：不起 goroutine、不碰 ctx。要开始消费超时关单信号，
// 由上层在进程生命周期里调 Start(ctx)（见 start.go）。
func NewHandler(
	orderCommand repo.OrderCommand,
	orderItemCommand repo.OrderItemCommand,
	customerCommand repo.CustomerCommand,
	productCommand repo.ProductCommand,
	courseService ports.CourseService,
	pay ports.PayService,
	tx ports.Transaction,
	idservice ports.IdempotencyService,
	queue ports.DelayScheduler,
) *Handler {
	h := &Handler{
		orderCommand:     orderCommand,
		orderItemCommand: orderItemCommand,
		customerCommand:  customerCommand,
		productCommand:   productCommand,
		courseService:    courseService,
		pay:              pay,
		tx:               tx,
		idservice:        idservice,
		queue:            queue,
	}

	return h
}
