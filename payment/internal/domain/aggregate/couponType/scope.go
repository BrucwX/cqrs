package couponType

import (
	"maps"
	"slices"
	"strings"
)

// --- 值对象 (Value Object) ---

// Scope 是「这张券能用在哪类商品上」的值对象。
//
// 两种形态：
//   - 通用：任何商品都能用（NewUniversalScope）；
//   - 限定：只能用在指定**商品分类**下的商品（NewLimitedScope）。
//
// 按分类（product.Type 的 code）而不是按单个商品限定：新商品只要归到该分类就自动可用，
// 不用维护白名单。分类是值对象、没有 id，所以这里也用 code 当键。
type Scope struct {
	universal        bool
	productTypeCodes map[string]struct{}
}

// NewUniversalScope 通用券：全场商品都能用。
func NewUniversalScope() Scope {
	return Scope{universal: true}
}

// NewLimitedScope 限定券：只能用在指定商品分类下，至少要给一个分类 code。
//
// 入参会去掉空白项并去重。
func NewLimitedScope(productTypeCodes []string) (Scope, error) {
	codes := toProductTypeCodeSet(productTypeCodes)
	if len(codes) == 0 {
		return Scope{}, ErrInvalidScope
	}
	return Scope{productTypeCodes: codes}, nil
}

// ReconstituteScope 仓储恢复：字段全部由外部给定，不做校验（内部会复制一份）。
func ReconstituteScope(universal bool, productTypeCodes map[string]struct{}) Scope {
	return Scope{universal: universal, productTypeCodes: copyProductTypeCodeSet(productTypeCodes)}
}

// Covers 某个商品分类能不能用这张券（限定范围里含它才算）。
func (s Scope) Covers(productTypeCode string) bool {
	if s.universal {
		return true
	}
	_, ok := s.productTypeCodes[productTypeCode]
	return ok
}

// valid 是否是构造出来的合法范围（挡住零值 Scope：既非通用、又没有分类）。
func (s Scope) valid() bool {
	return s.universal || len(s.productTypeCodes) > 0
}

// toProductTypeCodeSet 去掉空白项并去重，把入参收成集合。
func toProductTypeCodeSet(productTypeCodes []string) map[string]struct{} {
	set := make(map[string]struct{}, len(productTypeCodes))
	for _, code := range productTypeCodes {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		set[code] = struct{}{}
	}
	return set
}

// copyProductTypeCodeSet 复制一份，免得外部拿着同一个 map 改到值对象里面。
func copyProductTypeCodeSet(productTypeCodes map[string]struct{}) map[string]struct{} {
	if len(productTypeCodes) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(productTypeCodes))
	for code := range productTypeCodes {
		set[code] = struct{}{}
	}
	return set
}

// --- 属性只读访问器 (Getters) ---

// ProductTypeCodes 限定的商品分类 code（升序；通用时为空）。
//
// 集合本身无序，这里排好序返回，方便展示与测试断言。
func (s Scope) ProductTypeCodes() []string {
	return slices.Sorted(maps.Keys(s.productTypeCodes))
}

func (s Scope) IsUniversal() bool { return s.universal }
