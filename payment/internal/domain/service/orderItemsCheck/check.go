package orderitemscheck

import (
	"context"

	"payment/internal/domain/aggregate/orderItem"
)

// CanOrderCreate 下单前的预检：把几个单项检查串起来，一个不过就整体不过。
//
// 真正的规则都在各自的 CheckXxx 里（商品库存 / 受赠人可交易），这里只负责「取数 + 逐个过」。
//
// ⚠️ 预检只是「提前说不行」：真正扣库存是下单时 LimitedProduct.ReserveStock 的事，
// 那一步才挡得住并发（这里到那里之间的 TOCTOU 谁都躲不掉）。
func (s *Service) CanOrderCreate(ctx context.Context, orderItemIDs []int64) (bool, error) {
	orderItems, err := s.orderItemCommand.MustGetOrderItems(ctx, orderItemIDs)
	if err != nil {
		return false, err
	}

	if orderItems.AllDraft() {
		return false, orderItem.ErrOrderItemNotBelongOrder
	}

	for _, check := range []func(context.Context, orderItem.OrderItems) (bool, error){
		s.CheckProductQuantity,
		s.CheckCustomerReceive,
	} {
		ok, err := check(ctx, orderItems)
		if err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

// CheckProductQuantity 每件商品的库存够不够。
//
// 汇总（同一个商品多行会累加）是集合自己的事，判定是商品自己的事（HasStockFor）——
// 服务类商品压根没有库存这回事，无限类商品没有上限，两个都会答「够」。
func (s *Service) CheckProductQuantity(ctx context.Context, orderItems orderItem.OrderItems) (bool, error) {
	// 按商品汇总（同一个商品多行会累加）是集合自己的事，这里只管一个个商品去问库存。
	// 接出来存一份：汇总会遍历一遍并新建 map，别在循环条件里反复求值（读起来也容易误会）。
	needByProduct := orderItems.QuantityByProduct()
	for productID, need := range needByProduct {
		p, err := s.productCommand.MustGetProduct(ctx, productID)
		if err != nil {
			return false, err
		}
		if !p.HasStockFor(need) {
			return false, nil
		}
	}

	return true, nil
}

// CheckCustomerReceive 每一行的**受赠人**都得存在、而且能交易（冻结的客户不能收东西）。
//
// 自己买给自己、还是买来送人，规则是同一个 —— 收东西的人本身得是「可交易」状态。
func (s *Service) CheckCustomerReceive(ctx context.Context, orderItems orderItem.OrderItems) (bool, error) {
	for _, customerID := range orderItems.CustomerIDs() {
		customer, err := s.customerCommand.MustGetCustomer(ctx, customerID)
		if err != nil {
			return false, err
		}
		if !customer.CanTrade() {
			return false, nil
		}
	}
	return true, nil
}
