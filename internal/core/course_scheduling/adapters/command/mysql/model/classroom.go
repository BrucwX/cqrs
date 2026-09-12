package model

import "cqrs/internal/core/course_scheduling/domain/aggregate/classroom"

// Classroom 对应表 classroom，领域模型是 classroom.Classroom。
//
// 值对象 Location 就地打平为 Building / Floor / Room。
//
// 注意：这张表**没有** created_at / updated_at —— 领域模型里教室也不带时间戳，
// 列清单与建表脚本保持一致，不额外补列。
type Classroom struct {
	ID        string // varchar(36) 教室 ID（uuid）
	Building  string // varchar(64) 楼栋        <- Location.Building
	Floor     int    // int         楼层        <- Location.Floor
	Room      string // varchar(32) 房间号      <- Location.Room
	Capacity  int    // int         总容量（座位数）
	Allocated int    // int         已分配座位数

	// uint8 tinyint unsigned 1 可用 / 2 维护中 / 3 已报废
	Status uint8

	LockVersion uint64 // bigint unsigned 乐观锁版本号
}

// ClassroomToPO 写路径：Location 值对象打平成三列。
func ClassroomToPO(do *classroom.Classroom) *Classroom {
	location := do.Location()
	return &Classroom{
		ID:        do.ID(),
		Building:  location.Building(),
		Floor:     location.Floor(),
		Room:      location.Room(),
		Capacity:  do.Capacity(),
		Allocated: do.AllocatedSeats(),
		Status:    uint8(do.Status()),
	}
}

// ClassroomToDO 读路径：三列重新拼回 Location，楼栋/房间号为空会被拒绝。
func ClassroomToDO(po *Classroom) (*classroom.Classroom, error) {
	location, err := classroom.NewLocation(po.Building, po.Floor, po.Room)
	if err != nil {
		return nil, err
	}
	return classroom.Reconstitute(
		po.ID, location, po.Capacity, po.Allocated, classroom.Status(po.Status),
	), nil
}
