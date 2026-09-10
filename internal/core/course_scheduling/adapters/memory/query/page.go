package query

// defaultPageSize 是 Page 系列查询在调用方未指定 pageSize 时使用的默认页大小。
const defaultPageSize = 20

// paginate 对已按 ID 排序的切片做 1-based 分页。
//
// page <= 0 视为第 1 页；pageSize <= 0 使用 defaultPageSize；
// 越界时返回空切片（不是 nil），方便调用方直接序列化。
func paginate[T any](items []T, page, pageSize int) []T {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	start := (page - 1) * pageSize
	if start >= len(items) {
		return []T{}
	}

	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

// indexBy 按 key 建立索引；出现重复 key 时保留先出现的那一个。
func indexBy[K comparable, V any](items []V, key func(V) K) map[K]V {
	out := make(map[K]V, len(items))
	for _, item := range items {
		k := key(item)
		if _, ok := out[k]; !ok {
			out[k] = item
		}
	}
	return out
}
