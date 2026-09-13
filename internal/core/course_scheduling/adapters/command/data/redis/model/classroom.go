package model

import (
	"fmt"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
)

// Classroom 是 Redis 快照（对应 classroom），领域模型是 classroom.Classroom。
//
// 值对象 Location 就地打平为 Building / Floor / Room。
//
// 注意：快照里**没有** created_at / updated_at —— 领域模型里教室也不带时间戳，
// 字段与领域模型保持一致，不额外补字段。
type Classroom struct {
	ID        string `json:"id"`        // 教室 ID（uuid）
	Building  string `json:"building"`  // 楼栋        <- Location.Building
	Floor     int    `json:"floor"`     // 楼层        <- Location.Floor
	Room      string `json:"room"`      // 房间号      <- Location.Room
	Capacity  int    `json:"capacity"`  // 总容量（座位数）
	Allocated int    `json:"allocated"` // 已分配座位数

	// 1 可用 / 2 维护中 / 3 已报废
	Status uint8 `json:"status"`
}

// ClassroomToRedis 写路径：Location 值对象打平成三个字段。
func ClassroomToRedis(do *classroom.Classroom) (*Classroom, error) {
	if do == nil {
		return nil, ErrClassroomToRedis
	}
	location := do.Location()
	return &Classroom{
		ID:        do.ID(),
		Building:  location.Building(),
		Floor:     location.Floor(),
		Room:      location.Room(),
		Capacity:  do.Capacity(),
		Allocated: do.AllocatedSeats(),
		Status:    uint8(do.Status()),
	}, nil
}

// ClassroomFromRedis 读路径：三个字段重新拼回 Location，楼栋/房间号为空会被拒绝。
func ClassroomFromRedis(po *Classroom) (*classroom.Classroom, error) {
	if po == nil {
		return nil, ErrClassroomFromRedis
	}
	location, err := classroom.NewLocation(po.Building, po.Floor, po.Room)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrClassroomFromRedis, err)
	}
	return classroom.Reconstitute(
		po.ID, location, po.Capacity, po.Allocated, classroom.Status(po.Status),
	), nil
}
