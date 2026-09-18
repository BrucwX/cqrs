package couponType

// --- 值对象 (Value Object) ---

// Quota 是「发放配额」值对象：是否限量、发行总量、已发数量。
//
// 值对象不可变 —— 占名额返回一份**新的**配额，由 CouponType 换掉自己那份，
// 免得有人拿着指针就地改。
//
// 不限量时 total 记 0、也不参与判断（issued 仍如实记录发过多少）。
type Quota struct {
	limited bool  // 是否限量
	total   int64 // 发行总量（不限量时为 0）
	issued  int64 // 已发数量
}

// NewUnlimitedQuota 不限量（配额零值就是它）。
func NewUnlimitedQuota() Quota {
	return Quota{}
}

// NewLimitedQuota 限量 total 张，total 必须大于 0。
func NewLimitedQuota(total int64) (Quota, error) {
	if total <= 0 {
		return Quota{}, ErrInvalidQuota
	}
	return Quota{limited: true, total: total}, nil
}

// ReconstituteQuota 仓储恢复：字段全部由外部给定，不做校验。
func ReconstituteQuota(limited bool, total int64, issued int64) Quota {
	return Quota{limited: limited, total: total, issued: issued}
}

// Reserve 占 n 个名额，返回占完之后的配额；不够就报 ErrCouponTypeQuotaExhausted。
//
// 返回新值而不是改自己：调用方要 `q, err = q.Reserve(1)` 这样换掉旧的那份。
func (q Quota) Reserve(n int64) (Quota, error) {
	if n <= 0 {
		return q, ErrInvalidQuota
	}
	if q.limited && q.issued+n > q.total {
		return q, ErrCouponTypeQuotaExhausted
	}
	q.issued += n
	return q, nil
}

// Remaining 剩余可发数量；不限量返回 -1。
func (q Quota) Remaining() int64 {
	if !q.limited {
		return -1
	}
	return q.total - q.issued
}

// IsExhausted 是否已发完（不限量永远不会发完）。
func (q Quota) IsExhausted() bool {
	return q.limited && q.issued >= q.total
}

// valid 是否是自洽的配额（限定总量必须大于 0；不限量随便）。
func (q Quota) valid() bool {
	return !q.limited || q.total > 0
}

// --- 属性只读访问器 (Getters) ---

func (q Quota) IsLimited() bool { return q.limited }
func (q Quota) Total() int64    { return q.total }
func (q Quota) Issued() int64   { return q.issued }
