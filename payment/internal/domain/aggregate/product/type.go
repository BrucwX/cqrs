package product

import "strings"

// --- 值对象 (Value Object) ---

// Type 是「商品分类」值对象：这件商品属于哪一类（课程 / 数码 / 服务包）。
//
// 分类不是一个聚合 —— 它随商品一起存、随商品一起读。值对象没有自己的 id：
//   - **code 就是它的标识**，也是别处引用它的键（券的适用范围 Scope 按分类 code 限定）
//   - name 只是展示名，可以改
//
// 没有「启用 / 停用」—— 那是分类作为聚合才管得住的状态；现在分类跟着商品走，
// 要停售某个分类就直接改商品状态。
type Type struct {
	name string // 展示名
	code string // 业务稳定键（如 course / digital）：标识 + 被引用，建好后不该改
}

// NewType 构造分类值对象：name 与 code 都必填。
func NewType(name string, code string) (Type, error) {
	t := Type{name: name, code: code}
	if !t.valid() {
		return Type{}, ErrInvalidProductType
	}
	return t, nil
}

// ReconstituteType 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteType(name string, code string) Type {
	return Type{name: name, code: code}
}

// valid 是否是构造出来的合法分类（挡住零值 Type）。
func (t Type) valid() bool {
	return strings.TrimSpace(t.name) != "" && strings.TrimSpace(t.code) != ""
}

// Is 两个分类是不是同一个（按 code 判）。
func (t Type) Is(other Type) bool { return t.code == other.code }

// --- 属性只读访问器 (Getters) ---

func (t Type) Name() string { return t.name }
func (t Type) Code() string { return t.code }
