package ordercheck

import "payment/internal/domain/repo"

// Service 负责**订单侧**的跨聚合根判定：订单归属、订单状态这类问题。
//
// 判定不留存储：取数交给 repo（OrderCheckRepo 只声明本服务用得到的那几个方法），
// 规则写在方法里。一批**行项**的预检（库存、重复拥有、受赠人能不能收）归
// orderItemsCheck —— 一个服务对一类问题。
type Service struct {
	repo repo.OrderCheckRepo
}

// NewService 构造领域服务：仓储端口由上层（wire）注入。
func NewService(repo repo.OrderCheckRepo) *Service {
	return &Service{repo: repo}
}
