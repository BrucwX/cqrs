package orderitemscheck

import (
	"context"

	"payment/internal/domain/aggregate/orderItem"
)

// CanOrderCreate 下单前的预检：这批行项能不能拿去下单。
//
// 现在只有一条：**都还是草稿** —— 有一行已经归到别的订单上，就不能再拿去建新单
// （否则同一批行项会被两张订单共用）。
//
// ⚠️ 预检只是「提前说不行」：判完之后到真正下单之间，这批行项仍可能被别人用掉，
// 兜底靠的是下单时把行项 ConfirmAll 到订单上那一步（这段 TOCTOU 谁都躲不掉）。
func (s *Service) CanOrderCreate(ctx context.Context, orderItemIDs []int64) (bool, error) {
	orderItems, err := s.repo.MustGetOrderItems(ctx, orderItemIDs)
	if err != nil {
		return false, err
	}
	// 要下单的这批行项必须**都还是草稿**：有一行已经归到别的订单上，就不能再拿去
	// 建新单（否则同一批行项会被两张订单共用）。
	if !orderItems.AllDraft() {
		return false, orderItem.ErrOrderItemAlreadyConfirmed
	}

	return true, nil
}
