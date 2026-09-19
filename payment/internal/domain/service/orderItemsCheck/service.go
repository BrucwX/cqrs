package orderitemscheck

import (
	"context"

	"payment/internal/domain/aggregate/customer"
	"payment/internal/domain/aggregate/orderItem"
	"payment/internal/domain/aggregate/product"
)

// Repo 是**本服务自己要用的那几个方法**：需要什么就声明什么，不去依赖聚合包里的
// OrderItemCommand / ProductCommand / CustomerCommand。实现侧只要签名对得上就能注入。
type Repo interface {
	MustGetOrderItems(ctx context.Context, orderItemIDs []int64) (orderItem.OrderItems, error)
	MustGetProduct(ctx context.Context, productID int64) (product.Product, error)
	MustGetCustomer(ctx context.Context, customerID int64) (customer.Customer, error)
}

// Service 负责**一批订单行**的预检：这些行项能不能拿去下单。
//
// 它只读三种聚合：行项（要查的这批）、商品（库存 / 重复拥有）、客户（受赠人能不能收）。
// 订单本身的归属与状态判定归 orderCheck —— 一个服务对一类问题，别互相塞。
type Service struct {
	repo Repo
}

// NewService 构造领域服务：Repo 由上层（wire）注入。
func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}
