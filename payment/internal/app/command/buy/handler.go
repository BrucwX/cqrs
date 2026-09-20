package buy

// Package buy 是「买课」用例（写侧）。
//
// 一个用例一个 Handler：拿到写侧仓储端口 + 领域判定服务 + 事务端口，
// 按「先判、过了才写」的顺序来。

import (
	appRepo "payment/internal/app/command/appRepo"
	"payment/internal/domain/repo"
	timer "payment/internal/pkg"
)

type Handler struct {
	orderCommand     repo.OrderCommand
	orderItemCommand repo.OrderItemCommand
	customerCommand  repo.CustomerCommand
	productCommand   repo.ProductCommand
	courseService    appRepo.CourseService
	pay              appRepo.PayService
	tx               appRepo.Transaction
	scheduler        timer.Scheduler
}

// NewHandler 构造：依赖由上层（wire）注入。
func NewHandler(
	orderCommand repo.OrderCommand,
	orderItemCommand repo.OrderItemCommand,
	customerCommand repo.CustomerCommand,
	productCommand repo.ProductCommand,
	courseService appRepo.CourseService,
	pay appRepo.PayService,
	tx appRepo.Transaction,
	scheduler timer.Scheduler,
) *Handler {
	return &Handler{
		orderCommand:     orderCommand,
		orderItemCommand: orderItemCommand,
		customerCommand:  customerCommand,
		productCommand:   productCommand,
		courseService:    courseService,
		pay:              pay,
		tx:               tx,
		scheduler:        scheduler,
	}
}
