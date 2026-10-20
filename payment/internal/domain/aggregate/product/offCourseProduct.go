package product

import (
	"strings"
)

// --- 具体商品类型 (Concrete Product) ---

// ProductType 课程商品的分类：**这一类商品的 Type() 永远是它**，调用方不用传、也没得改。
//
// 券的适用范围就是按这个 code 匹配的 —— 「只限课程类」的券在 Scope 里直接写
// string(product.ProductType)。
const productType Type = "offlineCourse"

type OffCourseProduct struct {
	ProductMeta
	targetID string
}

// 编译期断言：漏了内嵌 ProductMeta（或方法签名写歪）在这里就会编译失败。
var _ Product = (*OffCourseProduct)(nil)

// NewOfflineCourseP 建课程商品：分类固定是 ProductType，默认「未上架」，要 PutOnSale 之后才能被下单。
func NewOfflineCourseP(
	name string,
	description string,
	price int64,
	targetID string,
) (*OffCourseProduct, error) {
	if err := validateCommon(name, price); err != nil {
		return nil, err
	}
	if strings.TrimSpace(targetID) == "" {
		return nil, ErrProductTargetRequired
	}

	return &OffCourseProduct{
		ProductMeta: newMeta(name, description, price),
		targetID:    targetID,
	}, nil
}

// ReconstituteOfflineCourseP 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteOfflineCourseP(meta ProductMeta, targetID string) *OffCourseProduct {
	return &OffCourseProduct{
		ProductMeta: meta,
		targetID:    targetID,
	}
}

// --- 属性只读访问器 (Getters) ---
func (c *OffCourseProduct) Type() Type { return productType }

func (c *OffCourseProduct) TargetId() string { return c.targetID }
