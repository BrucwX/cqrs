package orderitemscheck

import "payment/internal/domain/repo"

// Service 负责**一批订单行**的预检：这些行项能不能拿去下单。
//
// 它只读三种聚合：行项（要查的这批）、商品（库存 / 重复拥有）、客户（受赠人能不能收），
// 需要的三个方法声明在 repo.OrderItemsCheckRepo 里 —— 服务不背「完整仓储」的包袱。
// 订单本身的归属与状态判定归 orderCheck。
type Service struct {
	repo repo.OrderItemsCheckRepo
}

// NewService 构造领域服务：仓储端口由上层（wire）注入。
func NewService(repo repo.OrderItemsCheckRepo) *Service {
	return &Service{repo: repo}
}
