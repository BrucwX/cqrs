package product

// --- 值对象 (Value Object) ---

// Stock 是「库存」值对象：是否不限量 + 剩余可售数量 + 销售量。
//
// 三个数各管一件事：
//   - unlimited：是否不限量（课程、虚拟商品这类要多少有多少）
//   - remaining：剩余可售，出库减、进货加、退单加
//   - sold：累计销售量，出库加、退单减、**进货不动**
//
// 值对象不可变 —— 每个操作都返回一份**新的**库存，由 Product 换掉自己那份，
// 免得有人拿着指针就地改。
type Stock struct {
	unlimited bool  // 是否不限量
	remaining int64 // 剩余可售数量（不限量时不参与判断，记 0）
	sold      int64 // 累计销售量（退单会减回去）
}

// NewUnlimitedStock 不限量库存，比如课程、虚拟商品。
func NewUnlimitedStock() Stock {
	return Stock{unlimited: true}
}

// NewLimitedStock 有限库存，初始剩余 remaining 件（>= 0，可以为 0，比如先进货后上架）。
func NewLimitedStock(remaining int64) (Stock, error) {
	if remaining < 0 {
		return Stock{}, ErrInvalidStock
	}
	return Stock{remaining: remaining}, nil
}

// ReconstituteStock 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteStock(unlimited bool, remaining int64, sold int64) Stock {
	return Stock{unlimited: unlimited, remaining: remaining, sold: sold}
}

// Reserve 出库 n 件（下单）：剩余减 n、销售量加 n；剩余不够报 ErrInsufficientStock。
//
// 不限量库存不动剩余，只记销售量。调用方要 `s, err = s.Reserve(1)` 换掉旧的那份。
func (s Stock) Reserve(n int64) (Stock, error) {
	if n <= 0 {
		return s, ErrInvalidStock
	}
	if !s.unlimited {
		if n > s.remaining {
			return s, ErrInsufficientStock
		}
		s.remaining -= n
	}
	s.sold += n
	return s, nil
}

// Release 退单 n 件（订单取消 / 退款）：剩余加 n、销售量减 n —— 卖出去了又退回来。
//
// 退的件数不能超过累计销售量（不然销量就成负数了）。
func (s Stock) Release(n int64) (Stock, error) {
	if n <= 0 || n > s.sold {
		return s, ErrInvalidStock
	}
	if !s.unlimited {
		s.remaining += n
	}
	s.sold -= n
	return s, nil
}

// Restock 进货 n 件（补货）：只加剩余，不动销售量。
func (s Stock) Restock(n int64) (Stock, error) {
	if n <= 0 {
		return s, ErrInvalidStock
	}
	if !s.unlimited {
		s.remaining += n
	}
	return s, nil
}

// Remaining 剩余可售数量（不限量时记 0，判断请用 IsUnlimited）。
func (s Stock) Remaining() int64 { return s.remaining }

// Sold 累计销售量。
func (s Stock) Sold() int64 { return s.sold }

// IsUnlimited 是否不限量。
func (s Stock) IsUnlimited() bool { return s.unlimited }

// IsExhausted 是否已售罄（不限量永远不会售罄）。
func (s Stock) IsExhausted() bool { return !s.unlimited && s.remaining <= 0 }

// valid 是否是自洽的库存（有限库存的剩余不能为负；销售量不能为负）。
func (s Stock) valid() bool {
	return (s.unlimited || s.remaining >= 0) && s.sold >= 0
}
