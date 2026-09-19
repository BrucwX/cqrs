package ordercheck

import "payment/internal/domain/repo/command"

// Service 负责**订单侧**的跨聚合根判定：订单归属、订单状态这类问题。
//
// 取数自己做（拿着写侧仓库接口），判定不留存储。一批**行项**的预检
// （库存、重复拥有、受赠人能不能收）归 orderItemsCheck，别往这里塞。
type Service struct {
	customerCommand command.CustomerCommand
	orderCommand    command.OrderCommand
}

// NewService 构造领域服务：仓库接口由上层（wire）注入。
func NewService(
	customerCommand command.CustomerCommand,
	orderCommand command.OrderCommand,
) *Service {
	return &Service{
		customerCommand: customerCommand,
		orderCommand:    orderCommand,
	}
}
