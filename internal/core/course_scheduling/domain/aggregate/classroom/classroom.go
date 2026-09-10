package classroom

import "errors"

// --- 聚合根 (Aggregate Root) ---

type Classroom struct {
	id        string
	location  Location
	capacity  int
	allocated int
	status    Status
}

// 创建教室聚合根
func NewClassroom(id string, location Location, capacity int) (*Classroom, error) {
	if id == "" {
		return nil, errors.New("classroom ID is required")
	}
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	return &Classroom{
		id:        id,
		location:  location,
		capacity:  capacity,
		allocated: 0,
		status:    StatusAvailable,
	}, nil
}

// --- 核心业务行为 (Domain Behaviors) ---

// CanAccommodate 检查教室是否能够容纳指定人数上课（校验状态与可用容量）
func (c *Classroom) CanAccommodate(requiredSeats int) error {
	if c.status != StatusAvailable {
		return ErrClassroomUnavailable
	}
	if requiredSeats > c.capacity-c.allocated {
		return ErrCapacityExceeded
	}
	return nil
}

// AllocateSeats 分配座位，容量递减
func (c *Classroom) AllocateSeats(n int) error {
	if n <= 0 {
		return errors.New("allocated seats must be positive")
	}
	if c.allocated+n > c.capacity {
		return ErrCapacityExceeded
	}
	c.allocated += n
	return nil
}

// ReleaseSeats 释放座位，容量递增
func (c *Classroom) ReleaseSeats(n int) error {
	if n <= 0 {
		return errors.New("released seats must be positive")
	}
	if n > c.allocated {
		return errors.New("cannot release more seats than allocated")
	}
	c.allocated -= n
	return nil
}

// StartMaintenance 开启维护状态（如报修后停用排课）
func (c *Classroom) StartMaintenance() {
	c.status = StatusUnderMaintenance
}

// FinishMaintenance 完成维护，恢复可用
func (c *Classroom) FinishMaintenance() {
	c.status = StatusAvailable
}

// UpdateCapacity 调整教室容量
func (c *Classroom) UpdateCapacity(newCapacity int) error {
	if newCapacity <= 0 {
		return ErrInvalidCapacity
	}
	c.capacity = newCapacity
	return nil
}

// --- 属性只读访问器 (Getters) ---

func (c *Classroom) ID() string            { return c.id }
func (c *Classroom) Location() Location    { return c.location }
func (c *Classroom) Capacity() int         { return c.capacity }
func (c *Classroom) AllocatedSeats() int   { return c.allocated }
func (c *Classroom) AvailableSeats() int   { return c.capacity - c.allocated }
func (c *Classroom) Status() Status        { return c.status }
func (c *Classroom) IsAvailable() bool     { return c.status == StatusAvailable }
