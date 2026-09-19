package ordercheck

import (
	"context"

	"payment/internal/domain/aggregate/customer"
	"payment/internal/domain/aggregate/order"
)

// Repo 是**本服务自己要用的那几个方法**：需要什么就声明什么，不去依赖聚合包里的
// CustomerCommand / OrderCommand。实现侧（adapters 的载体）只要签名对得上就能注入 ——
// 服务不认识整个仓储，只认识自己用到的这一小块。
type Repo interface {
	MustGetCustomer(ctx context.Context, customerID int64) (customer.Customer, error)
	MustGetOrder(ctx context.Context, orderID int64) (order.Order, error)
}

// Service 负责**订单侧**的跨聚合根判定：订单归属、订单状态这类问题。
//
// 判定不留存储：取数交给 Repo，规则写在方法里。一批**行项**的预检
// （库存、重复拥有、受赠人能不能收）归 orderItemsCheck，别往这里塞。
type Service struct {
	repo Repo
}

// NewService 构造领域服务：Repo 由上层（wire）注入。
func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}
