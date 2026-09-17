package classroom

import (
	"errors"

	"github.com/google/uuid"
)

// --- 聚合根 (Aggregate Root) ---

type Classroom struct {
	id        string
	location  Location
	capacity  int
	allocated int
	status    Status
}

// NewClassroom 新建教室（ID 由聚合自己生成）
//
// location 与 capacity 必填；status 传 nil 就默认「可用」，
// 报废的教室不能直接建出来。
func NewClassroom(location *Location, capacity *int, status *Status) (*Classroom, error) {
	if location == nil || location.building == "" || location.room == "" {
		return nil, ErrInvalidLocation
	}
	if capacity == nil || *capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	created := &Classroom{
		id:       uuid.New().String(),
		location: *location,
		capacity: *capacity,
		status:   StatusAvailable,
	}

	if status != nil {
		if err := created.ChangeStatus(*status); err != nil {
			return nil, err
		}
	}
	return created, nil
}

// Reconstitute 仓储/种子数据恢复：ID、已分配座位与状态都由外部给定。
func Reconstitute(
	id string,
	location Location,
	capacity int,
	allocated int,
	status Status,
) *Classroom {
	return &Classroom{
		id:        id,
		location:  location,
		capacity:  capacity,
		allocated: allocated,
		status:    status,
	}
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

// UpdateLocation 调整教室位置
func (c *Classroom) UpdateLocation(location Location) error {
	if location.building == "" || location.room == "" {
		return ErrInvalidLocation
	}
	c.location = location
	return nil
}

// ChangeStatus 切换教室状态
//
// 只支持在「可用 / 维护中」之间切；报废（StatusDecommissioned）
// 要经过专门的废弃流程，这里直接拒绝。
func (c *Classroom) ChangeStatus(status Status) error {
	switch status {
	case StatusAvailable:
		c.FinishMaintenance()
	case StatusUnderMaintenance:
		c.StartMaintenance()
	default:
		return ErrUnsupportedStatusChange
	}
	return nil
}

// Update 按非 nil 的字段更新教室，nil 的字段保持原值。
func (c *Classroom) Update(location *Location, capacity *int, status *Status) error {
	if location != nil {
		if err := c.UpdateLocation(*location); err != nil {
			return err
		}
	}
	if capacity != nil {
		if err := c.UpdateCapacity(*capacity); err != nil {
			return err
		}
	}
	if status != nil {
		if err := c.ChangeStatus(*status); err != nil {
			return err
		}
	}
	return nil
}

// --- 属性只读访问器 (Getters) ---

func (c *Classroom) ID() string          { return c.id }
func (c *Classroom) Location() Location  { return c.location }
func (c *Classroom) Capacity() int       { return c.capacity }
func (c *Classroom) AllocatedSeats() int { return c.allocated }
func (c *Classroom) AvailableSeats() int { return c.capacity - c.allocated }
func (c *Classroom) Status() Status      { return c.status }
func (c *Classroom) IsAvailable() bool   { return c.status == StatusAvailable }
