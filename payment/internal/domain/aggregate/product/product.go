package product

import (
	"strings"

	"payment/pkg/uuid"
)

// --- 聚合根 (Aggregate Root) ---

// Product 是「一件可以卖的东西」的聚合根：课程 / 服务包这类可售项目。
//
// 它只管商品自身（名字、描述、价格、上下架）；一笔交易买了几件、花了多少钱，
// 是 order / orderItem 的事。
//
// 价格用「分」的 int64。下架不是删除 —— 已经卖出去的订单行项持有下单时的价格
// 快照，不受商品改名/调价/下架影响。
type Product struct {
	id          int64
	name        string
	description string
	price       int64 // 售价（分）
	stock       Stock // 库存：是否不限量 / 剩余可售 / 销售量
	status      Status
}

// NewProduct 建商品：默认「未上架」，要 PutOnSale 之后才能被下单。
func NewProduct(name string, description string, price int64, stock Stock) (*Product, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrProductNameRequired
	}
	if price <= 0 {
		return nil, ErrInvalidPrice
	}
	if !stock.valid() {
		return nil, ErrInvalidStock
	}

	return &Product{
		id:          uuid.NewInt64(),
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		status:      StatusDraft,
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(
	id int64,
	name string,
	description string,
	price int64,
	stock Stock,
	status Status,
) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		status:      status,
	}
}

// --- 上架 / 下架 ---

// PutOnSale 上架：只有上架的商品能被下单。
func (p *Product) PutOnSale() error {
	if p.status == StatusOnSale {
		return ErrProductAlreadyOnSale
	}
	p.status = StatusOnSale
	return nil
}

// TakeOffShelf 下架：之后不再接新单，已经卖出去的不受影响。
func (p *Product) TakeOffShelf() error {
	if p.status == StatusOffShelf {
		return ErrProductAlreadyOffShelf
	}
	p.status = StatusOffShelf
	return nil
}

// --- 商品信息 ---

// ChangePrice 调价（分）：只影响之后下的单，已下的单用快照价。
func (p *Product) ChangePrice(price int64) error {
	if price <= 0 {
		return ErrInvalidPrice
	}
	p.price = price
	return nil
}

// Rename 改名：同样只影响之后的展示，订单行项持有的是下单时的名字快照。
func (p *Product) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrProductNameRequired
	}
	p.name = name
	return nil
}

// UpdateDescription 改描述（可以为空，所以不做校验）。
func (p *Product) UpdateDescription(description string) {
	p.description = description
}

// --- 库存 ---

// ReserveStock 出库 n 件（下单时调用）：剩余减、销售量加。
//
// 剩余不够直接拒绝；不限量库存永远够。
func (p *Product) ReserveStock(n int64) error {
	next, err := p.stock.Reserve(n)
	if err != nil {
		return err
	}
	p.stock = next
	return nil
}

// ReleaseStock 退单 n 件（订单取消 / 退款时调用）：剩余加回来、销售量减回去。
func (p *Product) ReleaseStock(n int64) error {
	next, err := p.stock.Release(n)
	if err != nil {
		return err
	}
	p.stock = next
	return nil
}

// RestockStock 进货 n 件（补货时调用）：只加剩余，不动销售量。
func (p *Product) RestockStock(n int64) error {
	next, err := p.stock.Restock(n)
	if err != nil {
		return err
	}
	p.stock = next
	return nil
}

// IsOnSale 是否在售（可以被下单）。
func (p *Product) IsOnSale() bool { return p.status == StatusOnSale }

// --- 属性只读访问器 (Getters) ---

func (p *Product) ID() int64           { return p.id }
func (p *Product) Name() string        { return p.name }
func (p *Product) Description() string { return p.description }
func (p *Product) Price() int64        { return p.price }
func (p *Product) Stock() Stock        { return p.stock }
func (p *Product) Status() Status      { return p.status }
