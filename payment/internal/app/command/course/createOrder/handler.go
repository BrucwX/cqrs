// Package course 是「买课」用例（写侧）。
//
// 一个用例一个 Handler：拿到写侧仓储端口 + 领域判定服务 + 事务端口，
// 按「先判、过了才写」的顺序来。
//
// 本文件只放「Handler 是什么、依赖谁」；具体用例在各自的文件里（createOrder.go …）。
package course

import (
	"payment/internal/domain/repo"
	"payment/internal/ports"
)

type Handler struct {
	orderCommand     repo.OrderCommand
	orderItemCommand repo.OrderItemCommand
	customerCommand  repo.CustomerCommand
	productCommand   repo.ProductCommand
	courseService    ports.CourseService
	pay              ports.PayService
	tx               ports.Transaction
	idservice        ports.IdempotencyService
	p                ports.ClosePayProducer
}

// NewHandler 构造：依赖由上层（wire）注入。
func NewHandler(
	orderCommand repo.OrderCommand,
	orderItemCommand repo.OrderItemCommand,
	customerCommand repo.CustomerCommand,
	productCommand repo.ProductCommand,
	courseService ports.CourseService,
	pay ports.PayService,
	tx ports.Transaction,
	idservice ports.IdempotencyService,
	p ports.ClosePayProducer,
) *Handler {
	return &Handler{
		orderCommand:     orderCommand,
		orderItemCommand: orderItemCommand,
		customerCommand:  customerCommand,
		productCommand:   productCommand,
		courseService:    courseService,
		pay:              pay,
		tx:               tx,
		idservice:        idservice,
		p:                p,
	}
}
