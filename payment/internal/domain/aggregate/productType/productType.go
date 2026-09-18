package productType

import (
	"strings"

	"payment/pkg/uuid"
)

// --- 聚合根 (Aggregate Root) ---

// ProductType 是「一类商品」的聚合根：商品的分类（如 课程 / 数码 / 服务包）。
//
// 它只管分类自身（叫什么、编码、是否启用）。某件商品属于哪一类记在 product 的
// productTypeID 上 —— 商品只持有本聚合的 ID，本聚合不认识它。
//
// code 是业务侧引用的稳定键（如 course / digital），建好后不要改：改了等于断引用。
type ProductType struct {
	id     int64
	name   string
	code   string // 业务编码，如 course / digital
	status Status
}

// NewProductType 建一个商品分类：默认启用。
func NewProductType(name string, code string) (*ProductType, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrProductTypeNameRequired
	}
	if strings.TrimSpace(code) == "" {
		return nil, ErrProductTypeCodeRequired
	}

	return &ProductType{
		id:     uuid.NewInt64(),
		name:   name,
		code:   code,
		status: StatusEnabled,
	}, nil
}

// Reconstitute 仓储恢复：字段全部由外部给定，不做校验。
func Reconstitute(id int64, name string, code string, status Status) *ProductType {
	return &ProductType{
		id:     id,
		name:   name,
		code:   code,
		status: status,
	}
}

// --- 启用 / 停用 ---

// Enable 启用：启用后新商品才能挂到这个分类上。
func (t *ProductType) Enable() error {
	if t.status == StatusEnabled {
		return ErrProductTypeAlreadyEnabled
	}
	t.status = StatusEnabled
	return nil
}

// Disable 停用：不再接新商品，已经挂在这个分类下的商品不受影响。
func (t *ProductType) Disable() error {
	if t.status == StatusDisabled {
		return ErrProductTypeAlreadyDisabled
	}
	t.status = StatusDisabled
	return nil
}

// --- 分类信息 ---

// Rename 改名（只影响展示；业务侧按 code 引用，不受影响）。
func (t *ProductType) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrProductTypeNameRequired
	}
	t.name = name
	return nil
}

// IsEnabled 当前是否可用（能挂新商品）。
func (t *ProductType) IsEnabled() bool { return t.status == StatusEnabled }

// --- 属性只读访问器 (Getters) ---

func (t *ProductType) ID() int64      { return t.id }
func (t *ProductType) Name() string   { return t.name }
func (t *ProductType) Code() string   { return t.code }
func (t *ProductType) Status() Status { return t.status }
