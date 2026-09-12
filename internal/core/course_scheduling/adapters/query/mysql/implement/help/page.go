package help

// defaultPageSize 分页默认值与内存适配器保持一致，两个实现互换时行为不变。
const defaultPageSize = 20

// LimitOffset 把 1-based 的页码换算成 SQL 的 LIMIT / OFFSET。
//
// page <= 0 视为第 1 页；pageSize <= 0 用默认页大小。
// 越界页由 SQL 自然返回 0 行，不用特判。
func LimitOffset(page, pageSize int) (limit, offset int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	return pageSize, (page - 1) * pageSize
}
