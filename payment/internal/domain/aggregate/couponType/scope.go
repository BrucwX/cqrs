package couponType

import (
	"maps"
	"slices"
)

// --- 值对象 (Value Object) ---

// Scope 是「这张券能用在哪类商品上」的值对象。
//
// 两种形态：
//   - 通用：任何商品都能用（NewUniversalScope）；
//   - 限定：只能用在指定**商品类别**（productType）下的商品（NewLimitedScope）。
//
// 按类别而不是按单个商品限定：新商品只要归到该类别就自动可用，不用维护白名单。
type Scope struct {
	universal      bool
	productTypeIDs map[int64]struct{}
}

// NewUniversalScope 通用券：全场商品都能用。
func NewUniversalScope() Scope {
	return Scope{universal: true}
}

// NewLimitedScope 限定券：只能用在指定商品类别下，至少要给一个类别。
//
// 入参会丢掉非正数并去重。
func NewLimitedScope(productTypeIDs []int64) (Scope, error) {
	ids := toProductTypeIDSet(productTypeIDs)
	if len(ids) == 0 {
		return Scope{}, ErrInvalidScope
	}
	return Scope{productTypeIDs: ids}, nil
}

// ReconstituteScope 仓储恢复：字段全部由外部给定，不做校验（内部会复制一份）。
func ReconstituteScope(universal bool, productTypeIDs map[int64]struct{}) Scope {
	return Scope{universal: universal, productTypeIDs: copyProductTypeIDSet(productTypeIDs)}
}

// Covers 某个商品类别能不能用这张券（限定范围里含它才算）。
func (s Scope) Covers(productTypeID int64) bool {
	if s.universal {
		return true
	}
	_, ok := s.productTypeIDs[productTypeID]
	return ok
}

// valid 是否是构造出来的合法范围（挡住零值 Scope：既非通用、又没有类别）。
func (s Scope) valid() bool {
	return s.universal || len(s.productTypeIDs) > 0
}

// toProductTypeIDSet 丢掉非正数并去重，把入参收成集合。
func toProductTypeIDSet(productTypeIDs []int64) map[int64]struct{} {
	set := make(map[int64]struct{}, len(productTypeIDs))
	for _, id := range productTypeIDs {
		if id <= 0 {
			continue
		}
		set[id] = struct{}{}
	}
	return set
}

// copyProductTypeIDSet 复制一份，免得外部拿着同一个 map 改到值对象里面。
func copyProductTypeIDSet(productTypeIDs map[int64]struct{}) map[int64]struct{} {
	if len(productTypeIDs) == 0 {
		return nil
	}
	set := make(map[int64]struct{}, len(productTypeIDs))
	for id := range productTypeIDs {
		set[id] = struct{}{}
	}
	return set
}

// --- 属性只读访问器 (Getters) ---

// ProductTypeIDs 限定的商品类别 ID（升序；通用时为空）。
//
// 集合本身无序，这里排好序返回，方便展示与测试断言。
func (s Scope) ProductTypeIDs() []int64 {
	return slices.Sorted(maps.Keys(s.productTypeIDs))
}

func (s Scope) IsUniversal() bool { return s.universal }
