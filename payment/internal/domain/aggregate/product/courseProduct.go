package product

import "strings"

// --- 具体商品类型 (Concrete Product) ---

// ProductType 课程商品的分类 code：**这一类商品的 Type() 永远是它**，调用方不用传、也没得改。
//
// 券的适用范围就是按这个 code 匹配的 —— 「只限课程类」的券在 Scope 里直接写 product.ProductType。
const ProductType = "offlineCourse"

type OfflineCourseP struct {
	ProductMeta
	targetID string
}

// 编译期断言：漏了内嵌 ProductMeta（或方法签名写歪）在这里就会编译失败。
var _ Product = (*OfflineCourseP)(nil)

// NewOfflineCourseP 建课程商品：分类固定是 CourseType，默认「未上架」，要 PutOnSale 之后才能被下单。
func NewOfflineCourseP(
	name string,
	description string,
	price int64,
	targetID string,
) (*OfflineCourseP, error) {
	if err := validateCommon(name, price); err != nil {
		return nil, err
	}
	if strings.TrimSpace(targetID) == "" {
		return nil, ErrProductTargetRequired
	}

	return &OfflineCourseP{
		ProductMeta: newMeta(name, description, price),
		targetID:    targetID,
	}, nil
}

// ReconstituteOfflineCourseP 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteOfflineCourseP(meta ProductMeta, targetID string) *OfflineCourseP {
	return &OfflineCourseP{
		ProductMeta: meta,
		targetID:    targetID,
	}
}

// --- 只有课程商品才有的行为 ---

// ChangeTarget 改挂到另一门课程上（原来的课下架了，这个商品位改成卖新开的课）。
func (c *OfflineCourseP) ChangeTarget(targetID string) error {
	if strings.TrimSpace(targetID) == "" {
		return ErrProductTargetRequired
	}
	c.targetID = targetID
	return nil
}

// --- 属性只读访问器 (Getters) ---
func (c *OfflineCourseP) Type() string { return ProductType }


func (c *OfflineCourseP) TargetId() string { return c.targetID }
