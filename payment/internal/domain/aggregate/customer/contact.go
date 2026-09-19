package customer

import "strings"

// --- 值对象 (Value Object) ---

// Contact 是客户的联系方式（值对象）：姓名 + 手机号 + 可选邮箱。
//
// 联系方式是整体替换的 —— 改手机号就是换一整个 Contact，不做单字段修改，
// 这样「姓名与手机号必须同时有效」这条约束只能被构造时校验一次。
type Contact struct {
	name  string
	phone string
	email string
}

// NewContact 构造联系方式：姓名与手机号必填，邮箱可选。
func NewContact(name string, phone string, email string) (Contact, error) {
	if strings.TrimSpace(name) == "" {
		return Contact{}, ErrCustomerNameRequired
	}
	if strings.TrimSpace(phone) == "" {
		return Contact{}, ErrCustomerPhoneRequired
	}
	return Contact{name: name, phone: phone, email: email}, nil
}

// ReconstituteContact 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteContact(name string, phone string, email string) Contact {
	return Contact{name: name, phone: phone, email: email}
}

// valid 是否是构造出来的合法联系方式（Customer 用它挡掉零值 Contact 入参）。
func (c Contact) valid() bool {
	return strings.TrimSpace(c.name) != "" && strings.TrimSpace(c.phone) != ""
}

// --- 属性只读访问器 (Getters) ---

func (c Contact) Name() string  { return c.name }
func (c Contact) Phone() string { return c.phone }
func (c Contact) Email() string { return c.email }
