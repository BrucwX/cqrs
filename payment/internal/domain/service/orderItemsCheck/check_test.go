package orderitemscheck

import (
	"context"
	"errors"
	"testing"

	"payment/internal/domain/aggregate/orderItem"
)

// 领域服务只依赖仓库**接口**，所以这里不用数据库：塞一个假实现，
// 只实现服务真要用的那一个方法（repo.OrderItemsCheckRepo 现在只有 MustGetOrderItems），
// 顺便记下被问过的行项 ID，好断言「取数只取一次」。

const (
	bookID  int64 = 100
	aliceID int64 = 200
	bobID   int64 = 201
)

// --- 假仓库 ---

type fakeRepo struct {
	items   orderItem.OrderItems
	itemErr error

	askedItemIDs [][]int64 // 每次被问的行项 ID
}

func (f *fakeRepo) MustGetOrderItems(_ context.Context, orderItemIDs []int64) (orderItem.OrderItems, error) {
	f.askedItemIDs = append(f.askedItemIDs, orderItemIDs)
	if f.itemErr != nil {
		return nil, f.itemErr
	}
	return f.items, nil
}

// --- 夹具 ---

// draft 草稿行：还没归到订单上（下单前的那批就该长这样）。
func draft(t *testing.T, customerID int64, productID int64, quantity int32) orderItem.DraftOrderItem {
	t.Helper()
	item, err := orderItem.NewDraftOrderItem(customerID, productID, "教材", 3500, quantity)
	if err != nil {
		t.Fatal(err)
	}
	return item
}

// confirmed 正式行：已经归到订单 1 上了。
func confirmed(t *testing.T, id int64, customerID int64, productID int64) orderItem.ConfirmedOrderItem {
	t.Helper()
	meta := orderItem.ReconstituteMeta(id, customerID, productID, "教材", 3500, 1)
	return orderItem.ReconstituteConfirmedOrderItem(meta, 1)
}

// newService 组装一个被假仓库喂饱的业务服务。
func newService(items orderItem.OrderItems, itemErr error) (*Service, *fakeRepo) {
	fake := &fakeRepo{items: items, itemErr: itemErr}
	return NewService(fake), fake
}

// assertResult 断言一次检查的 (bool, error)：给了 wantErr 就只认那个错误。
func assertResult(t *testing.T, wantOK bool, wantErr error, ok bool, err error) {
	t.Helper()
	if wantErr != nil {
		if !errors.Is(err, wantErr) {
			t.Fatalf("期望错误 %v，实际 %v", wantErr, err)
		}
		return
	}
	if err != nil {
		t.Fatalf("不该有错误，实际 %v", err)
	}
	if ok != wantOK {
		t.Fatalf("期望 ok=%v，实际 ok=%v", wantOK, ok)
	}
}

// --- CanOrderCreate ---

func TestCanOrderCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("全是草稿 → 通过", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(orderItem.OrderItems{
			draft(t, aliceID, bookID, 2),
			draft(t, bobID, bookID, 1),
		}, nil)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, true, nil, ok, err)
	})

	t.Run("单行草稿 → 通过", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(orderItem.OrderItems{draft(t, aliceID, bookID, 1)}, nil)

		ok, err := service.CanOrderCreate(ctx, []int64{1})
		assertResult(t, true, nil, ok, err)
	})

	t.Run("有一行已经归到订单上 → ErrOrderItemAlreadyConfirmed", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(orderItem.OrderItems{
			draft(t, aliceID, bookID, 1),
			confirmed(t, 2, aliceID, bookID),
		}, nil)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, false, orderItem.ErrOrderItemAlreadyConfirmed, ok, err)
	})

	t.Run("正式行排在最前面也一样拦下（跟顺序无关）", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(orderItem.OrderItems{
			confirmed(t, 1, aliceID, bookID),
			draft(t, aliceID, bookID, 1),
		}, nil)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, false, orderItem.ErrOrderItemAlreadyConfirmed, ok, err)
	})

	t.Run("全部都不是草稿 → 拦下", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(orderItem.OrderItems{
			confirmed(t, 1, aliceID, bookID),
			confirmed(t, 2, bobID, bookID),
		}, nil)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, false, orderItem.ErrOrderItemAlreadyConfirmed, ok, err)
	})

	t.Run("空集合 → 通过（没有行项要下单）", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(orderItem.OrderItems{}, nil)

		ok, err := service.CanOrderCreate(ctx, nil)
		assertResult(t, true, nil, ok, err)
	})

	t.Run("取行项失败 → 错误透传", func(t *testing.T) {
		t.Parallel()

		boom := errors.New("db 挂了")
		service, _ := newService(nil, boom)

		ok, err := service.CanOrderCreate(ctx, []int64{1})
		assertResult(t, false, boom, ok, err)
	})

	t.Run("行项只取一次（草稿判定用不着别的聚合）", func(t *testing.T) {
		t.Parallel()

		service, repo := newService(orderItem.OrderItems{
			draft(t, aliceID, bookID, 1),
			draft(t, bobID, bookID, 1),
		}, nil)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, true, nil, ok, err)
		if len(repo.askedItemIDs) != 1 {
			t.Fatalf("行项只该取一次，实际取了 %d 次：%v", len(repo.askedItemIDs), repo.askedItemIDs)
		}
	})
}
