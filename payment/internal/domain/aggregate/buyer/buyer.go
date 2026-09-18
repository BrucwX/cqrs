package buyer

import "payment/pkg/uuid"

// --- 聚合根 (Aggregate Root) ---

// Buyer 是「下单的人」的聚合根：身份、状态与手上的积分。
//
// 订单在 order 聚合、退款在 refund 聚合，它们只持有 buyerID 单向引用本聚合。
// 积分直接放在这里：只记一个余额，不记流水、不分类型、不管适用范围。
type Buyer struct {
	id      int64
	userID  int64
	contact Contact
	status  Status
	points  int64
}

// NewBuyer 注册买家：默认「正常」，冻结要走 Freeze（风控动作）。
func NewBuyer(contact Contact) (*Buyer, error) {
	if !contact.valid() {
		return nil, ErrBuyerContactRequired
	}

	return &Buyer{
		id:      uuid.NewInt64(),
		contact: contact,
		status:  StatusActive,
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(id int64, contact Contact, status Status, points int64) *Buyer {
	return &Buyer{
		id:      id,
		contact: contact,
		status:  status,
		points:  points,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// ChangeContact 换联系方式（姓名/手机/邮箱整体替换）。
func (b *Buyer) ChangeContact(contact Contact) error {
	if !contact.valid() {
		return ErrBuyerContactRequired
	}
	b.contact = contact
	return nil
}

// Freeze 冻结（风控 / 违规）：冻结后不能再下单、不能再申请退款。
func (b *Buyer) Freeze() error {
	if b.status == StatusFrozen {
		return ErrBuyerAlreadyFrozen
	}
	b.status = StatusFrozen
	return nil
}

// Unfreeze 解冻。
func (b *Buyer) Unfreeze() error {
	if b.status != StatusFrozen {
		return ErrBuyerNotFrozen
	}
	b.status = StatusActive
	return nil
}

// CanTrade 当前是否允许交易（下单 / 退款）。
func (b *Buyer) CanTrade() bool { return b.status == StatusActive }

// --- 积分 (Points) ---

// EarnPoints 增加积分（订单支付返积分 / 活动赠送）。
func (b *Buyer) EarnPoints(points int64) error {
	if points <= 0 {
		return ErrInvalidPoints
	}
	b.points += points
	return nil
}

// SpendPoints 扣减积分（下单抵扣 / 兑换），余额不足直接拒绝。
func (b *Buyer) SpendPoints(points int64) error {
	if points <= 0 {
		return ErrInvalidPoints
	}
	if points > b.points {
		return ErrInsufficientPoints
	}
	b.points -= points
	return nil
}

// AdjustPoints 人工调整（补偿 / 纠错）：delta 可正可负，但调整后余额不能为负。
func (b *Buyer) AdjustPoints(delta int64) error {
	if delta == 0 {
		return ErrInvalidPoints
	}
	if b.points+delta < 0 {
		return ErrInsufficientPoints
	}
	b.points += delta
	return nil
}

// CanAffordPoints 积分够不够扣（只读，抵扣前预判用）。
func (b *Buyer) CanAffordPoints(points int64) bool {
	return points > 0 && points <= b.points
}

// --- 属性只读访问器 (Getters) ---

func (b *Buyer) ID() int64        { return b.id }
func (b *Buyer) UserID() int64    { return b.userID }
func (b *Buyer) Contact() Contact { return b.contact }
func (b *Buyer) Status() Status   { return b.status }
func (b *Buyer) Points() int64    { return b.points }

// 下面是联系方式的透传访问器，省得调用方每次都写 b.Contact().Name()。

func (b *Buyer) Name() string  { return b.contact.Name() }
func (b *Buyer) Phone() string { return b.contact.Phone() }
func (b *Buyer) Email() string { return b.contact.Email() }
