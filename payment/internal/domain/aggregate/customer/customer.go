package customer

import "payment/pkg/uuid"

// --- 聚合根 (Aggregate Root) ---

// Customer 是「下单的人」的聚合根：身份、状态与手上的积分。
//
// 订单在 order 聚合、退款在 refund 聚合，它们只持有 customerID 单向引用本聚合。
// 积分直接放在这里：只记一个余额，不记流水、不分类型、不管适用范围。
type Customer struct {
	id        int64
	accountID int64 // 这个客户对应的**登录账号** ID（账号体系在别处）；聚合自己的主键是 id
	contact   Contact
	status    Status
	points    int64
}

// NewCustomer 注册客户：默认「正常」，冻结要走 Freeze（风控动作）。
func NewCustomer(contact Contact) (*Customer, error) {
	if !contact.valid() {
		return nil, ErrCustomerContactRequired
	}

	return &Customer{
		id:      uuid.NewInt64(),
		contact: contact,
		status:  StatusActive,
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(id int64, contact Contact, status Status, points int64) *Customer {
	return &Customer{
		id:      id,
		contact: contact,
		status:  status,
		points:  points,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// ChangeContact 换联系方式（姓名/手机/邮箱整体替换）。
func (c *Customer) ChangeContact(contact Contact) error {
	if !contact.valid() {
		return ErrCustomerContactRequired
	}
	c.contact = contact
	return nil
}

// Freeze 冻结（风控 / 违规）：冻结后不能再下单、不能再申请退款。
func (c *Customer) Freeze() error {
	if c.status == StatusFrozen {
		return ErrCustomerAlreadyFrozen
	}
	c.status = StatusFrozen
	return nil
}

// Unfreeze 解冻。
func (c *Customer) Unfreeze() error {
	if c.status != StatusFrozen {
		return ErrCustomerNotFrozen
	}
	c.status = StatusActive
	return nil
}

// CanTrade 当前是否允许交易（下单 / 退款）。
func (c *Customer) CanTrade() bool { return c.status == StatusActive }

// --- 积分 (Points) ---

// EarnPoints 增加积分（订单支付返积分 / 活动赠送）。
func (c *Customer) EarnPoints(points int64) error {
	if points <= 0 {
		return ErrInvalidPoints
	}
	c.points += points
	return nil
}

// SpendPoints 扣减积分（下单抵扣 / 兑换），余额不足直接拒绝。
func (c *Customer) SpendPoints(points int64) error {
	if points <= 0 {
		return ErrInvalidPoints
	}
	if points > c.points {
		return ErrInsufficientPoints
	}
	c.points -= points
	return nil
}

// AdjustPoints 人工调整（补偿 / 纠错）：delta 可正可负，但调整后余额不能为负。
func (c *Customer) AdjustPoints(delta int64) error {
	if delta == 0 {
		return ErrInvalidPoints
	}
	if c.points+delta < 0 {
		return ErrInsufficientPoints
	}
	c.points += delta
	return nil
}

// CanAffordPoints 积分够不够扣（只读，抵扣前预判用）。
func (c *Customer) CanAffordPoints(points int64) bool {
	return points > 0 && points <= c.points
}

// --- 属性只读访问器 (Getters) ---

func (c *Customer) ID() int64        { return c.id }
func (c *Customer) AccountID() int64 { return c.accountID }
func (c *Customer) Contact() Contact { return c.contact }
func (c *Customer) Status() Status   { return c.status }
func (c *Customer) Points() int64    { return c.points }

// 下面是联系方式的透传访问器，省得调用方每次都写 c.Contact().Name()。

func (c *Customer) Name() string  { return c.contact.Name() }
func (c *Customer) Phone() string { return c.contact.Phone() }
func (c *Customer) Email() string { return c.contact.Email() }
