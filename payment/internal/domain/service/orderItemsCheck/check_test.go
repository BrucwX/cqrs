package orderitemscheck

import (
	"context"
	"errors"
	"testing"

	"payment/internal/domain/aggregate/customer"
	"payment/internal/domain/aggregate/orderItem"
	"payment/internal/domain/aggregate/product"
)

// 领域服务只依赖仓库**接口**，所以这里不用数据库：塞几个假实现，
// 顺便记下「问过哪些商品 / 客户」，好断言去重与查询次数。
//
// 商品与客户都用 Reconstitute 造 —— 构造函数会随机生成 ID，而订单行里写的是具体
// 商品 ID，两边必须对得上。

const (
	bookID  int64 = 100 // 有限类商品：还剩 5 本
	flowID  int64 = 101 // 无限类商品：没有库存上限，按 GB 卖
	albumID int64 = 102 // 服务类商品：压根没有库存这回事

	aliceID int64 = 200
	bobID   int64 = 201
)

// --- 假仓库 ---
//
// 服务只声明自己要用的方法（Repo），所以一个假实现就够了 —— 顺便记下
// 「问过哪些商品 / 客户」，好断言去重与查询次数。

type fakeRepo struct {
	items     orderItem.OrderItems
	itemErr   error
	products  map[int64]product.Product
	customers map[int64]customer.Customer

	askedProducts  []int64
	askedCustomers []int64
}

func (f *fakeRepo) MustGetOrderItems(context.Context, []int64) (orderItem.OrderItems, error) {
	if f.itemErr != nil {
		return nil, f.itemErr
	}
	return f.items, nil
}

func (f *fakeRepo) MustGetProduct(_ context.Context, productID int64) (product.Product, error) {
	f.askedProducts = append(f.askedProducts, productID)
	p, ok := f.products[productID]
	if !ok {
		return nil, product.ErrProductNotFound
	}
	return p, nil
}

func (f *fakeRepo) MustGetCustomer(_ context.Context, customerID int64) (customer.Customer, error) {
	f.askedCustomers = append(f.askedCustomers, customerID)
	c, ok := f.customers[customerID]
	if !ok {
		return customer.Customer{}, customer.ErrCustomerNotFound
	}
	return c, nil
}

// --- 夹具 ---

func testType(t *testing.T) product.Type {
	t.Helper()
	productType, err := product.NewType("课程", "course")
	if err != nil {
		t.Fatal(err)
	}
	return productType
}

// limitedProduct 有限类商品：remaining 本库存。
func limitedProduct(t *testing.T, remaining int64) product.Product {
	t.Helper()
	meta := product.ReconstituteMeta(bookID, "教材", "第一版", 3500, testType(t), product.StatusOnSale)
	return product.ReconstituteLimitedProduct(meta, "本", remaining, 0)
}

// unlimitedProduct 无限类商品：没库存上限（问它 HasStockFor 永远是够）。
func unlimitedProduct(t *testing.T) product.Product {
	t.Helper()
	meta := product.ReconstituteMeta(flowID, "流量包", "全国流量", 1, testType(t), product.StatusOnSale)
	return product.ReconstituteUnlimitedProduct(meta, "GB")
}

// serviceProduct 服务类商品：没有库存这回事。
func serviceProduct(t *testing.T) product.Product {
	t.Helper()
	meta := product.ReconstituteMeta(albumID, "会员年卡", "开通一年", 19900, testType(t), product.StatusOnSale)
	return product.ReconstituteServiceProduct(meta)
}

// activeCustomer 正常客户；frozen 时造一个被冻结的。
func testCustomer(t *testing.T, customerID int64, frozen bool) customer.Customer {
	t.Helper()
	contact, err := customer.NewContact("张三", "13800000000", "z@example.com")
	if err != nil {
		t.Fatal(err)
	}
	status := customer.StatusActive
	if frozen {
		status = customer.StatusFrozen
	}
	return *customer.Reconstitute(customerID, contact, status, 0)
}

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
func newService(
	items orderItem.OrderItems,
	itemErr error,
	products []product.Product,
	customers []customer.Customer,
) (*Service, *fakeRepo) {
	repo := &fakeRepo{
		items:     items,
		itemErr:   itemErr,
		products:  make(map[int64]product.Product, len(products)),
		customers: make(map[int64]customer.Customer, len(customers)),
	}
	for _, p := range products {
		repo.products[p.ID()] = p
	}
	for _, c := range customers {
		repo.customers[c.ID()] = c
	}
	return NewService(repo), repo
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

	t.Run("全是草稿 + 库存够 + 受赠人正常 → 通过", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(
			orderItem.OrderItems{draft(t, aliceID, bookID, 2), draft(t, bobID, flowID, 10)},
			nil,
			[]product.Product{limitedProduct(t, 5), unlimitedProduct(t)},
			[]customer.Customer{testCustomer(t, aliceID, false), testCustomer(t, bobID, false)},
		)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, true, nil, ok, err)
	})

	t.Run("有行项已经归到订单上 → ErrOrderItemAlreadyConfirmed", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(
			orderItem.OrderItems{draft(t, aliceID, bookID, 1), confirmed(t, 2, aliceID, bookID)},
			nil,
			[]product.Product{limitedProduct(t, 5)},
			[]customer.Customer{testCustomer(t, aliceID, false)},
		)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, false, orderItem.ErrOrderItemAlreadyConfirmed, ok, err)
	})

	t.Run("混合状态时先拦「不是草稿」，不去查库存", func(t *testing.T) {
		t.Parallel()

		service, repo := newService(
			orderItem.OrderItems{draft(t, aliceID, bookID, 3), confirmed(t, 2, aliceID, bookID)},
			nil,
			[]product.Product{limitedProduct(t, 1)}, // 库存也不够，但轮不到它报
			[]customer.Customer{testCustomer(t, aliceID, false)},
		)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, false, orderItem.ErrOrderItemAlreadyConfirmed, ok, err)
		if len(repo.askedProducts) != 0 {
			t.Fatalf("草稿检查没过就不该去查商品，实际查了 %v", repo.askedProducts)
		}
	})

	t.Run("同一商品两行加起来超库存 → 拦下", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(
			orderItem.OrderItems{draft(t, aliceID, bookID, 3), draft(t, aliceID, bookID, 3)},
			nil,
			[]product.Product{limitedProduct(t, 5)},
			[]customer.Customer{testCustomer(t, aliceID, false)},
		)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, false, nil, ok, err)
	})

	t.Run("无限类 / 服务类商品不看库存 → 通过", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(
			orderItem.OrderItems{draft(t, aliceID, flowID, 999), draft(t, aliceID, albumID, 999)},
			nil,
			[]product.Product{unlimitedProduct(t), serviceProduct(t)},
			[]customer.Customer{testCustomer(t, aliceID, false)},
		)

		ok, err := service.CanOrderCreate(ctx, []int64{1, 2})
		assertResult(t, true, nil, ok, err)
	})

	t.Run("受赠人被冻结 → 拦下", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(
			orderItem.OrderItems{draft(t, bobID, bookID, 1)},
			nil,
			[]product.Product{limitedProduct(t, 5)},
			[]customer.Customer{testCustomer(t, bobID, true)},
		)

		ok, err := service.CanOrderCreate(ctx, []int64{1})
		assertResult(t, false, nil, ok, err)
	})

	t.Run("取行项失败 → 错误透传", func(t *testing.T) {
		t.Parallel()

		boom := errors.New("db 挂了")
		service, _ := newService(nil, boom, nil, nil)

		ok, err := service.CanOrderCreate(ctx, []int64{1})
		assertResult(t, false, boom, ok, err)
	})

	t.Run("商品不存在 → 错误透传", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(
			orderItem.OrderItems{draft(t, aliceID, 999, 1)},
			nil,
			[]product.Product{limitedProduct(t, 5)},
			[]customer.Customer{testCustomer(t, aliceID, false)},
		)

		ok, err := service.CanOrderCreate(ctx, []int64{1})
		assertResult(t, false, product.ErrProductNotFound, ok, err)
	})

	t.Run("受赠人不存在 → 错误透传", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(
			orderItem.OrderItems{draft(t, 999, bookID, 1)},
			nil,
			[]product.Product{limitedProduct(t, 5)},
			nil,
		)

		ok, err := service.CanOrderCreate(ctx, []int64{1})
		assertResult(t, false, customer.ErrCustomerNotFound, ok, err)
	})
}

// --- CheckProductQuantity ---

func TestCheckProductQuantity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("同一商品多行会累加后比库存", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(nil, nil, []product.Product{limitedProduct(t, 5)}, nil)

		ok, err := service.CheckProductQuantity(ctx, orderItem.OrderItems{
			draft(t, aliceID, bookID, 2),
			draft(t, bobID, bookID, 3),
		})
		assertResult(t, true, nil, ok, err)

		ok, err = service.CheckProductQuantity(ctx, orderItem.OrderItems{
			draft(t, aliceID, bookID, 2),
			draft(t, bobID, bookID, 4),
		})
		assertResult(t, false, nil, ok, err)
	})

	t.Run("每件商品只查一次库存", func(t *testing.T) {
		t.Parallel()

		service, repo := newService(nil, nil, []product.Product{limitedProduct(t, 5)}, nil)

		ok, err := service.CheckProductQuantity(ctx, orderItem.OrderItems{
			draft(t, aliceID, bookID, 1),
			draft(t, bobID, bookID, 1),
			draft(t, aliceID, bookID, 1),
		})
		assertResult(t, true, nil, ok, err)
		if len(repo.askedProducts) != 1 {
			t.Fatalf("三行同一个商品只该查一次，实际 %v", repo.askedProducts)
		}
	})

	t.Run("卖光了 → 拦下", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(nil, nil, []product.Product{limitedProduct(t, 0)}, nil)

		ok, err := service.CheckProductQuantity(ctx, orderItem.OrderItems{draft(t, aliceID, bookID, 1)})
		assertResult(t, false, nil, ok, err)
	})

	t.Run("空集合 → 通过（没东西要查）", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(nil, nil, nil, nil)

		ok, err := service.CheckProductQuantity(ctx, orderItem.OrderItems{})
		assertResult(t, true, nil, ok, err)
	})
}

// --- CheckCustomerReceive ---

func TestCheckCustomerReceive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("受赠人都正常 → 通过", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(nil, nil, nil, []customer.Customer{
			testCustomer(t, aliceID, false),
			testCustomer(t, bobID, false),
		})

		ok, err := service.CheckCustomerReceive(ctx, orderItem.OrderItems{
			draft(t, aliceID, bookID, 1),
			draft(t, bobID, flowID, 1),
		})
		assertResult(t, true, nil, ok, err)
	})

	t.Run("同一个受赠人出现多行 → 只查一次", func(t *testing.T) {
		t.Parallel()

		service, repo := newService(nil, nil, nil,
			[]customer.Customer{testCustomer(t, aliceID, false)})

		ok, err := service.CheckCustomerReceive(ctx, orderItem.OrderItems{
			draft(t, aliceID, bookID, 1),
			draft(t, aliceID, flowID, 1),
			draft(t, aliceID, albumID, 1),
		})
		assertResult(t, true, nil, ok, err)
		if len(repo.askedCustomers) != 1 {
			t.Fatalf("三行同一个受赠人只该查一次，实际 %v", repo.askedCustomers)
		}
	})

	t.Run("其中一个人被冻结 → 拦下", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(nil, nil, nil, []customer.Customer{
			testCustomer(t, aliceID, false),
			testCustomer(t, bobID, true),
		})

		ok, err := service.CheckCustomerReceive(ctx, orderItem.OrderItems{
			draft(t, aliceID, bookID, 1),
			draft(t, bobID, bookID, 1),
		})
		assertResult(t, false, nil, ok, err)
	})

	t.Run("受赠人不存在 → 错误透传", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(nil, nil, nil, nil)

		ok, err := service.CheckCustomerReceive(ctx, orderItem.OrderItems{draft(t, bobID, bookID, 1)})
		assertResult(t, false, customer.ErrCustomerNotFound, ok, err)
	})

	t.Run("空集合 → 通过", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(nil, nil, nil, nil)

		ok, err := service.CheckCustomerReceive(ctx, orderItem.OrderItems{})
		assertResult(t, true, nil, ok, err)
	})
}
