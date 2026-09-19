package order

// Package order 是下单用例（写侧）。
//
// 一个用例一个 Handler：拿到写侧仓库接口 + 领域判定服务 + 事务端口，
// 按「先判、过了才写」的顺序来。

import (
	appRepo "payment/internal/app/command/appRepo"
	orderagg "payment/internal/domain/aggregate/order"
	"payment/internal/domain/aggregate/orderItem"
	orderitemscheck "payment/internal/domain/service/orderItemsCheck"
)

// Handler 下单用例。
//
// 依赖：两个**聚合包里的仓储接口**（order / orderItem）+ 领域判定服务 + 事务端口，
// 按「先判、过了才写」的顺序来。
//
// ⚠️ 暂时不需要 product 仓储 —— 下单不扣库存（等支付成功后再扣），
// 库存预检由 itemsCheck 自己拿商品查。
type Handler struct {
	orderCommand     orderagg.OrderCommand
	orderItemCommand orderItem.OrderItemCommand
	itemsCheck       *orderitemscheck.Service
	tx               appRepo.Transaction
}

// NewHandler 构造：依赖由上层（wire）注入。
func NewHandler(
	orderCommand orderagg.OrderCommand,
	orderItemCommand orderItem.OrderItemCommand,
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
