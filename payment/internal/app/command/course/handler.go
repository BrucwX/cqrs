package course

// Package course 是「买课」用例（写侧）。
//
// 一个用例一个 Handler：拿到写侧仓储端口 + 领域判定服务 + 事务端口，
// 按「先判、过了才写」的顺序来。

import (
	appRepo "payment/internal/app/command/appRepo"
	appService "payment/internal/app/command/appService"
	"payment/internal/domain/repo"
	orderitemscheck "payment/internal/domain/service/orderItemsCheck"
)

type Handler struct {
	orderCommand     repo.OrderCommand
	orderItemCommand repo.OrderItemCommand
	itemsCheck       *orderitemscheck.Service
	courseService    appService.CourseService
	tx               appRepo.Transaction
}

// NewHandler 构造：依赖由上层（wire）注入。
func NewHandler(
	orderCommand repo.OrderCommand,
	orderItemCommand repo.OrderItemCommand,
	itemsCheck *orderitemscheck.Service,
	tx appRepo.Transaction,
) *Handler {
	return &Handler{
		orderCommand:     orderCommand,
		orderItemCommand: orderItemCommand,
		itemsCheck:       itemsCheck,
		tx:               tx,
	}
}
