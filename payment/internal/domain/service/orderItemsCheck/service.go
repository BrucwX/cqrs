package orderitemscheck

import "payment/internal/domain/repo"

// Service 负责**一批订单行**的预检：这些行项能不能拿去下单。
//
// 它只读一种聚合：要查的那批行项；需要的方法声明在 repo.OrderItemsCheckRepo 里
// —— 服务不背「完整仓储」的包袱。订单本身的归属与状态判定归 orderCheck。
type Service struct {
	repo repo.OrderItemsCheckRepo
}

// NewService 构造领域服务：仓储端口由上层（wire）注入。
func NewService(repo repo.OrderItemsCheckRepo) *Service {
	return &Service{repo: repo}
}
