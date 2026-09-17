package course

import "errors"

// Capacity 配额值对象，负责约束容量不变量
type Capacity struct {
	max      int
	enrolled int
}

func NewCapacity(max int, enrolled int) (Capacity, error) {
	if max <= 0 || enrolled < 0 || enrolled > max {
		return Capacity{}, errors.New("invalid capacity boundaries")
	}
	return Capacity{max: max, enrolled: enrolled}, nil
}

func (c Capacity) IsFull() bool  { return c.enrolled >= c.max }
func (c Capacity) Max() int      { return c.max }
func (c Capacity) Enrolled() int { return c.enrolled }
