package orderitemscheck

import "payment/internal/domain/repo/command"

// Service 负责**一批订单行**的预检：这些行项能不能拿去下单。
//
// 它只读三种聚合：行项（要查的这批）、商品（库存 / 重复拥有）、客户（受赠人能不能收）。
// 订单本身的归属与状态判定归 orderCheck —— 一个服务对一类问题，别互相塞。
type Service struct {
	orderItemCommand command.OrderItemCommand
	productCommand   command.ProductCommand
	customerCommand  command.CustomerCommand
}

// NewService 构造领域服务：仓库接口由上层（wire）注入。
func NewService(
	orderItemCommand command.OrderItemCommand,
	productCommand command.ProductCommand,
	customerCommand command.CustomerCommand,
) *Service {
	return &Service{
		orderItemCommand: orderItemCommand,
		productCommand:   productCommand,
		customerCommand:  customerCommand,
	}
}
